package txobserver

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	slog "github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/store"
)

type Observer struct {
	ctx context.Context

	bridgeAccount string

	client *horizonclient.Client
	store  *store.DB
	log    *slog.Entry

	bridgeAccountCreateSequence uint32
}

func NewObserver(
	ctx context.Context,
	bridgeAccount string,
	bridgeAccountCreateSequence uint32,
	client *horizonclient.Client,
	store *store.DB,
) *Observer {
	o := &Observer{
		ctx:                         ctx,
		bridgeAccount:               bridgeAccount,
		bridgeAccountCreateSequence: bridgeAccountCreateSequence,
		client:                      client,
		store:                       store,
		log:                         slog.DefaultLogger.WithField("service", "stellar_txobserver"),
	}

	return o
}

func (o *Observer) ProcessNewLedgers() {
	for o.ctx.Err() == nil {
		if ledgerSeq, err := o.store.GetLastLedgerSequence(context.Background()); err != nil {
			o.log.Errorf("Unable to load last ledger sequence from db: %v", err)
		} else {
			if ledgerSeq == 0 {
				ledgerSeq = o.bridgeAccountCreateSequence
			} else {
				ledgerSeq++
			}
			// Get ledger data first to ensure there are no gaps
			if ledger, err := o.client.LedgerDetail(ledgerSeq); err != nil {
				if herr, ok := err.(*horizonclient.Error); ok && herr.Response.StatusCode == http.StatusNotFound {
					// Ledger not found means we reached the latest ledger
					return
				} else {
					o.log.WithField("error", err).Error("Error getting ledger details")
				}
			} else {
				o.log.WithField("sequence", ledgerSeq).Info("Processing ledger...")
				if err := o.processSingleLedger(ledger); err != nil {
					o.log.WithFields(slog.F{"error": err, "sequence": ledgerSeq}).Error("Error processing a single ledger details")
				} else {
					continue // without time.Sleep
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func (o *Observer) processSingleLedger(ledger horizon.Ledger) error {
	err := o.store.Session.Begin()
	if err != nil {
		return errors.Wrap(err, "error starting a transaction")
	}
	defer func() {
		// explicitly ignore return value to make the linter happy
		_ = o.store.Session.Rollback()
	}()
	// Process operations
	cursor := ""
	previousHash := ""
	for {
		ops, err := o.client.Operations(horizonclient.OperationRequest{
			ForLedger:     uint(ledger.Sequence),
			Cursor:        cursor,
			Limit:         200,
			IncludeFailed: true,
			Join:          "transactions",
		})
		if err != nil {
			return errors.Wrap(err, "error getting operations")
		}

		if len(ops.Embedded.Records) == 0 {
			break
		}

		for _, op := range ops.Embedded.Records {
			baseOp := op.GetBase()

			// Update cursor instantly because we can continue later
			cursor = op.PagingToken()

			// Ignore ops not coming from bridge account
			if baseOp.SourceAccount != o.bridgeAccount {
				continue
			}

			tx := baseOp.Transaction
			if tx.MemoType != "hash" || tx.Memo == "" || !tx.Successful ||
				// Skip inserting transactions with multiple ops. Currently Starbridge
				// does not create such transactions but it can change in the future.
				previousHash == baseOp.TransactionHash {
				continue
			}
			memoBytes, err := base64.StdEncoding.DecodeString(tx.Memo)
			if err != nil {
				return errors.Wrapf(err, "error decoding memo: %s", tx.Memo)
			}

			err = o.store.InsertHistoryStellarTransaction(context.TODO(), store.HistoryStellarTransaction{
				Hash:     tx.Hash,
				Envelope: tx.EnvelopeXdr,
				MemoHash: hex.EncodeToString(memoBytes),
			})
			if err != nil {
				return errors.Wrapf(err, "error inserting history transaction: %s", tx.Hash)
			}
			previousHash = baseOp.TransactionHash
		}
	}

	err = o.store.UpdateLastLedgerSequence(context.TODO(), uint32(ledger.Sequence))
	if err != nil {
		return errors.Wrap(err, "error updating last ledger sequence")
	}

	err = o.store.UpdateLastLedgerCloseTime(context.TODO(), ledger.ClosedAt)
	if err != nil {
		return errors.Wrap(err, "error updating last ledger sequence")
	}

	err = o.store.Session.Commit()
	if err != nil {
		return errors.Wrap(err, "error commiting a transaction")
	}

	o.log.WithField("sequence", ledger.Sequence).Info("Processed ledger")
	return nil
}

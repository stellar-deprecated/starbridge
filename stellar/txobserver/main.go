package txobserver

import (
	"context"
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

	ledgerSequence uint32
}

func NewObserver(ctx context.Context, bridgeAccount string, client *horizonclient.Client, store *store.DB) *Observer {
	o := &Observer{
		ctx:           ctx,
		bridgeAccount: bridgeAccount,
		client:        client,
		store:         store,
		log:           slog.DefaultLogger.WithField("service", "stellar_txobserver"),
	}

	ledgerSeq, err := o.store.GetLastLedgerSequence(context.Background())
	if err != nil {
		o.log.Fatalf("Unable to load last ledger sequence from db: %v", err)
	}

	if ledgerSeq == 0 {
		root, err := o.client.Root()
		if err != nil {
			o.log.Fatalf("Unable to access Horizon (%s) root resource: %v", client.HorizonURL, err)
		}

		o.ledgerSequence = uint32(root.HorizonSequence)
	} else {
		o.ledgerSequence = ledgerSeq
	}

	return o
}

func (o *Observer) ProcessNewLedgers() {
	for o.ctx.Err() == nil {
		// Get ledger data first to ensure there are no gaps
		ledger, err := o.client.LedgerDetail(o.ledgerSequence)
		if err != nil {
			if herr, ok := err.(*horizonclient.Error); ok && herr.Response.StatusCode == http.StatusNotFound {
				// Ledger not found means we reached the latest ledger
				return
			} else {
				o.log.WithField("error", err).Error("Error getting ledger details")
			}
		} else {
			o.log.WithField("sequence", o.ledgerSequence).Info("Processing ledger...")

			err = o.processSingleLedger(ledger)
			if err != nil {
				o.log.WithFields(slog.F{"error": err, "sequence": o.ledgerSequence}).Error("Error processing a single ledger details")
			} else {
				o.ledgerSequence++
				continue // without time.Sleep
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
		_ = o.store.Session.Rollback()
	}()

	// Get latest list of hashes to observe
	outgoingBridgeTransactions, err := o.store.GetPendingOutgoingStellarTransactions(context.TODO())
	if err != nil {
		return errors.Wrap(err, "error getting outgoing bridge transactions")
	}

	// Transform to hash->tx map for faster lookups
	outgoingBridgeTransactionsHashMap := make(map[string]store.OutgoingStellarTransaction)
	for _, tx := range outgoingBridgeTransactions {
		outgoingBridgeTransactionsHashMap[tx.Hash] = tx
	}

	// Process operations
	cursor := ""
	previousHash := ""
	for {
		ops, err := o.client.Operations(horizonclient.OperationRequest{
			ForLedger:     uint(o.ledgerSequence),
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

			memoHash := ""
			if tx.MemoType == "hash" {
				memoHash = tx.Memo
			}

			// Skip inserting transactions with multiple ops. Currently Starbridge
			// does not create such transactions but it can change in the future.
			if previousHash != baseOp.TransactionHash {
				err := o.store.InsertHistoryStellarTransaction(context.TODO(), store.HistoryStellarTransaction{
					Hash:     tx.Hash,
					Envelope: tx.EnvelopeXdr,
					MemoHash: memoHash,
				})
				if err != nil {
					return errors.Wrapf(err, "error inserting history transaction: %s", tx.Hash)
				}
				previousHash = baseOp.TransactionHash
			}

			otx, hashExists := outgoingBridgeTransactionsHashMap[baseOp.TransactionHash]
			if hashExists {
				if tx.Successful {
					otx.State = store.SuccessState
					err := o.store.MarkIncomingEthereumTransactionAsWithdrawn(context.TODO(), otx.IncomingTransactionHash)
					if err != nil {
						return errors.Wrapf(err, "error updating incoming ethereum transaction: %s", otx.IncomingTransactionHash)
					}
				} else {
					otx.State = store.InvalidState
				}
				err := o.store.UpsertOutgoingStellarTransaction(context.TODO(), otx)
				if err != nil {
					return errors.Wrapf(err, "error upserting outgoing transaction: %s", otx.Hash)
				}
			}

			if tx.InnerTransaction != nil {
				otx, hashExists := outgoingBridgeTransactionsHashMap[tx.InnerTransaction.Hash]
				if hashExists {
					if baseOp.Transaction.Successful {
						otx.State = store.SuccessState
						err := o.store.MarkIncomingEthereumTransactionAsWithdrawn(context.TODO(), otx.IncomingTransactionHash)
						if err != nil {
							return errors.Wrapf(err, "error updating incoming ethereum transaction: %s", otx.IncomingTransactionHash)
						}
					} else {
						otx.State = store.InvalidState
					}
					err := o.store.UpsertOutgoingStellarTransaction(context.TODO(), otx)
					if err != nil {
						return errors.Wrapf(err, "error upserting outgoing transaction: %s", otx.Hash)
					}
				}
			}
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

	o.log.WithField("sequence", o.ledgerSequence).Info("Processed ledger")
	return nil
}

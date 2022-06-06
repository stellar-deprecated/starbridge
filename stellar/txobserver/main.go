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

	client *horizonclient.Client
	store  *store.DB
	log    *slog.Entry

	// TODO: last ledgerSequence will be persisted in a DB
	ledgerSequence  uint32
	ledgerCloseTime time.Time
}

func NewObserver(ctx context.Context, client *horizonclient.Client, store *store.DB) *Observer {
	o := &Observer{
		ctx:    ctx,
		client: client,
		store:  store,
		log:    slog.DefaultLogger.WithField("service", "stellar_txobserver"),
	}

	root, err := o.client.Root()
	if err != nil {
		o.log.Fatal("Unable to access Horizon root resource")
	}

	o.ledgerSequence = uint32(root.HorizonSequence)

	return o
}

func (o *Observer) GetLastLedgerCloseTime() (time.Time, error) {
	return o.ledgerCloseTime, nil
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
				o.ledgerCloseTime = ledger.ClosedAt
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
	outgoingBridgeTransactions, err := o.store.GetOutgoingStellarTransactions(context.TODO())
	if err != nil {
		return errors.Wrap(err, "error getting outgoing bridge transactions")
	}

	// If no transactions to observe, skip to next ledger.
	if len(outgoingBridgeTransactions) == 0 {
		o.log.WithField("sequence", o.ledgerSequence).Info("No outgoing bridge transactions, skipping to next ledger")
		return nil
	}

	// Transform to hash->tx map for faster lookups
	outgoingBridgeTransactionsHashMap := make(map[string]store.OutgoingStellarTransaction)
	for _, tx := range outgoingBridgeTransactions {
		outgoingBridgeTransactionsHashMap[tx.Hash] = tx
	}

	// Transform to source->seqnum->tx map for faster lookups
	outgoingBridgeTransactionsSourceMap := make(map[string]map[int64]store.OutgoingStellarTransaction)
	for _, tx := range outgoingBridgeTransactions {
		if outgoingBridgeTransactionsSourceMap[tx.Source] == nil {
			outgoingBridgeTransactionsSourceMap[tx.Source] = make(map[int64]store.OutgoingStellarTransaction)
		}
		outgoingBridgeTransactionsSourceMap[tx.Source][tx.SequenceNumber] = tx
	}

	// Process transactions
	cursor := ""
	for {
		txs, err := o.client.Transactions(horizonclient.TransactionRequest{
			ForLedger:     uint(o.ledgerSequence),
			Cursor:        cursor,
			Limit:         200,
			IncludeFailed: true,
		})
		if err != nil {
			return errors.Wrap(err, "error getting transactions")
		}

		if len(txs.Embedded.Records) == 0 {
			break
		}

		for _, tx := range txs.Embedded.Records {
			otx, hashExists := outgoingBridgeTransactionsHashMap[tx.Hash]
			if hashExists {
				if tx.Successful {
					otx.State = store.SuccessState
					// TODO [important] update IncomingEthereumTransaction state
					// as withdrawn so it's not possible to cancel it.
				} else {
					otx.State = store.FailedState
				}
				err := o.store.UpsertOutgoingStellarTransaction(context.TODO(), otx)
				if err != nil {
					return errors.Wrapf(err, "error upserting outgoing transaction: %s", otx.Hash)
				}
			}

			if tx.InnerTransaction != nil {
				otx, hashExists := outgoingBridgeTransactionsHashMap[tx.InnerTransaction.Hash]
				if hashExists {
					if tx.Successful {
						otx.State = store.SuccessState
						// TODO [important] update IncomingEthereumTransaction state
						// as withdrawn so it's not possible to cancel it.
					} else {
						otx.State = store.FailedState
					}
					err := o.store.UpsertOutgoingStellarTransaction(context.TODO(), otx)
					if err != nil {
						return errors.Wrapf(err, "error upserting outgoing transaction: %s", otx.Hash)
					}
				}
			}

			seqnums, sourceExists := outgoingBridgeTransactionsSourceMap[tx.Account]
			if sourceExists {
				for seqnum, otx := range seqnums {
					if tx.AccountSequence >= seqnum && otx.Hash != tx.Hash &&
						(tx.InnerTransaction == nil ||
							(tx.InnerTransaction != nil && otx.Hash != tx.InnerTransaction.Hash)) {
						// Account sequence number was used for another transaction or bump thus otx transaction
						// is irrevocably invalid and can be marked as expired.
						otx.State = store.ExpiredState
						err := o.store.UpsertOutgoingStellarTransaction(context.TODO(), otx)
						if err != nil {
							return errors.Wrapf(err, "error upserting outgoing transaction: %s", otx.Hash)
						}
					}
				}
			}

			cursor = tx.PagingToken()
		}
	}

	// Mark all expired txs as expired.
	count, err := o.store.MarkOutgoingStellarTransactionExpired(context.TODO(), ledger.ClosedAt)
	if err != nil {
		return errors.Wrap(err, "error marking outgoing transactions as expired")
	}

	if count > 0 {
		o.log.Infof("Marked %d txs are expired", count)
	}

	err = o.store.Session.Commit()
	if err != nil {
		return errors.Wrap(err, "error commiting a transaction")
	}

	o.log.WithField("sequence", o.ledgerSequence).Info("Processed ledger")
	return nil
}

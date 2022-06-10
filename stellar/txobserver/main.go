package txobserver

import (
	"context"
	"net/http"
	"time"

	"github.com/stellar/go/clients/horizonclient"
	slog "github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/store"
)

type Observer struct {
	ctx context.Context

	client *horizonclient.Client
	store  *store.DB
	log    *slog.Entry

	// TODO: this will be persisted in a DB
	ledgerSequence uint32
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
		o.log.Fatalf("Unable to access Horizon (%s) root resource: %v", client.HorizonURL, err)
	}

	o.ledgerSequence = uint32(root.HorizonSequence)

	return o
}

func (o *Observer) ProcessNewLedgers() {
LedgerLoop:
	for {
		if o.ctx.Err() != nil {
			return
		}

		// Get ledger data first to ensure there are no gaps
		ledger, err := o.client.LedgerDetail(o.ledgerSequence)
		if err != nil {
			if herr, ok := err.(*horizonclient.Error); ok && herr.Response.StatusCode == http.StatusNotFound {
				// Ledger not found means we reached the latest ledger
				return
			} else {
				o.log.WithField("error", err).Error("Error getting ledger details")
			}
			time.Sleep(time.Second)
			continue
		}

		o.log.WithField("sequence", o.ledgerSequence).Info("Processing ledger...")

		// Get latest list of hashes to observe
		outgoingBridgeTransactions, err := o.store.GetOutgoingStellarTransactions(context.TODO())
		if err != nil {
			o.log.WithField("error", err).Error("Error getting outgoing bridge transactions")
			time.Sleep(time.Second)
			continue
		}

		// If no transactions to observe, skip to next ledger.
		if len(outgoingBridgeTransactions) == 0 {
			o.log.WithField("sequence", o.ledgerSequence).Info("No outgoing bridge transactions, skiping to next ledger")
			o.ledgerSequence++
			continue
		}

		// Transform to map for faster lookups
		outgoingBridgeTransactionsMap := make(map[string]store.OutgoingStellarTransaction)
		for _, tx := range outgoingBridgeTransactions {
			outgoingBridgeTransactionsMap[tx.Hash] = tx
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
				o.log.WithFields(slog.F{
					"error":  err,
					"ledger": o.ledgerSequence,
					"cursor": cursor,
				}).Error("Error getting transactions")
				time.Sleep(time.Second)
				continue LedgerLoop
			}

			if len(txs.Embedded.Records) == 0 {
				break
			}

			for _, tx := range txs.Embedded.Records {
				otx, hashExists := outgoingBridgeTransactionsMap[tx.Hash]
				if hashExists {
					if tx.Successful {
						otx.State = store.SuccessState
					} else {
						otx.State = store.FailedState
					}
					err := o.store.UpsertOutgoingStellarTransaction(context.TODO(), otx)
					if err != nil {
						o.log.WithFields(slog.F{
							"error": err,
							"hash":  otx.Hash,
						}).Error("Error upserting outgoing transactions")
						time.Sleep(time.Second)
						continue LedgerLoop
					}
				}

				if tx.InnerTransaction != nil {
					otx, hashExists := outgoingBridgeTransactionsMap[tx.InnerTransaction.Hash]
					if hashExists {
						if tx.Successful {
							otx.State = store.SuccessState
						} else {
							otx.State = store.FailedState
						}
						err := o.store.UpsertOutgoingStellarTransaction(context.TODO(), otx)
						if err != nil {
							o.log.WithFields(slog.F{
								"error": err,
								"hash":  otx.Hash,
							}).Error("Error upserting outgoing transactions")
							time.Sleep(time.Second)
							continue LedgerLoop
						}
					}
				}

				cursor = tx.PagingToken()
			}
		}

		// Mark all txs with expired time + buffer as expired.
		expiredBefore := ledger.ClosedAt.Add(-time.Minute)
		count, err := o.store.MarkOutgoingStellarTransactionExpired(context.TODO(), expiredBefore)
		if err != nil {
			o.log.WithField("error", err).Error("Error marking outgoing transactions as expired")
			time.Sleep(time.Second)
			continue LedgerLoop
		}

		o.log.Infof("Marked %d txs are expired", count)

		o.log.WithField("sequence", o.ledgerSequence).Info("Processed ledger")
		o.ledgerSequence++
	}
}

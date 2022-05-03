package txobserver

import (
	"net/http"
	"time"

	"github.com/stellar/go/clients/horizonclient"
	slog "github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/store"
)

type Observer struct {
	Client *horizonclient.Client
	Store  *store.Memory
}

func (o *Observer) Run() error {
	log := slog.DefaultLogger.WithField("service", "stellar_txobserver")

	root, err := o.Client.Root()
	if err != nil {
		log.Fatal("Unable to access Horizon root resource")
	}

	ledgerSequence := uint32(root.HorizonSequence)

	log.Infof("Starting Stellar observer from ledger: %d", ledgerSequence)

LedgerLoop:
	for {
		// TODO: this can slow down catchup
		time.Sleep(time.Second)

		// Get ledger data first to ensure there are no gaps
		ledger, err := o.Client.LedgerDetail(ledgerSequence)
		if err != nil {
			if herr, ok := err.(*horizonclient.Error); ok && herr.Response.StatusCode == http.StatusNotFound {
				log.WithField("sequence", ledgerSequence).Debug("Ledger not found, waiting...")
			} else {
				log.WithField("error", err).Error("Error getting ledger details")
			}
			continue
		}

		log.WithField("sequence", ledgerSequence).Info("Processing ledger...")

		// Get latest list of hashes to observe
		outgoingBridgeTransactions, err := o.Store.GetOutgoingStellarTransactions()
		if err != nil {
			log.WithField("error", err).Error("Error getting outgoing bridge transactions")
			continue
		}

		// If no transactions to observe, skip to next ledger.
		if len(outgoingBridgeTransactions) == 0 {
			log.WithField("sequence", ledgerSequence).Info("No outgoing bridge transactions, skiping to next ledger")
			ledgerSequence++
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
			txs, err := o.Client.Transactions(horizonclient.TransactionRequest{
				ForLedger:     uint(ledgerSequence),
				Cursor:        cursor,
				Limit:         200,
				IncludeFailed: true,
			})
			if err != nil {
				log.WithFields(slog.F{
					"error":  err,
					"ledger": ledgerSequence,
					"cursor": cursor,
				}).Error("Error getting transactions")
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
					err := o.Store.UpsertOutgoingStellarTransaction(otx)
					if err != nil {
						log.WithFields(slog.F{
							"error": err,
							"hash":  otx.Hash,
						}).Error("Error upserting outgoing transactions")
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
						err := o.Store.UpsertOutgoingStellarTransaction(otx)
						if err != nil {
							log.WithFields(slog.F{
								"error": err,
								"hash":  otx.Hash,
							}).Error("Error upserting outgoing transactions")
							continue LedgerLoop
						}
					}
				}

				cursor = tx.PagingToken()
			}
		}

		// Mark all txs with expired time + buffer as expired.
		expiredBefore := ledger.ClosedAt.Add(-time.Minute)
		count, err := o.Store.MarkOutgoingStellarTransactionExpired(expiredBefore)
		if err != nil {
			log.WithField("error", err).Error("Error marking outgoing transactions as expired")
			continue LedgerLoop
		}

		log.Infof("Marked %d txs are expired", count)

		log.WithField("sequence", ledgerSequence).Info("Processed ledger")
		ledgerSequence++
	}
}

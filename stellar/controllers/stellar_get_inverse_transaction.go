package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/stellar/go/strkey"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/store"
)

var validTxHash = regexp.MustCompile("^(0x)?([A-Fa-f0-9]{64})$")

type StellarWithdrawalHandler struct {
	StellarClient          *horizonclient.Client
	Observer               ethereum.Observer
	Store                  *store.DB
	WithdrawalWindow       time.Duration
	EthereumFinalityBuffer uint64
}

func (c *StellarWithdrawalHandler) getDeposit(r *http.Request, w http.ResponseWriter) (store.EthereumDeposit, error) {
	txHash := r.PostFormValue("transaction_hash")
	if !validTxHash.MatchString(txHash) {
		w.WriteHeader(http.StatusBadRequest)
		return store.EthereumDeposit{}, fmt.Errorf("invalid transaction hash")
	}
	parsed, err := strconv.ParseInt(r.PostFormValue("log_index"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return store.EthereumDeposit{}, fmt.Errorf("invalid log index")
	}
	logIndex := uint(parsed)
	depositID := ethereum.DepositID(txHash, logIndex)

	storeDeposit, err := c.Store.GetEthereumDeposit(r.Context(), depositID)
	if err == nil {
		return storeDeposit, nil
	} else if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		return store.EthereumDeposit{}, err
	}

	deposit, err := c.Observer.GetDeposit(r.Context(), txHash, logIndex)
	if ethereum.IsInvalidGetDepositRequest(err) {
		w.WriteHeader(http.StatusBadRequest)
		return store.EthereumDeposit{}, fmt.Errorf("invalid log index")
	} else if err == ethereum.ErrTxHashNotFound {
		w.WriteHeader(http.StatusNotFound)
		return store.EthereumDeposit{}, err
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return store.EthereumDeposit{}, err
	}

	block, err := c.Observer.GetLatestBlock(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return store.EthereumDeposit{}, err
	}
	if deposit.BlockNumber+c.EthereumFinalityBuffer > block.Number {
		w.WriteHeader(http.StatusPreconditionFailed)
		return store.EthereumDeposit{}, fmt.Errorf("need to wait for finality buffer")
	}

	destinationAccountID, err := strkey.Encode(
		strkey.VersionByteAccountID,
		deposit.Destination.Bytes(),
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return store.EthereumDeposit{}, fmt.Errorf("invalid stellar destination account %w", err)
	}
	storeDeposit = store.EthereumDeposit{
		ID:          depositID,
		Token:       deposit.Token.String(),
		Sender:      deposit.Sender.String(),
		Destination: destinationAccountID,
		Amount:      deposit.Amount.String(),
		Hash:        deposit.TxHash.String(),
		LogIndex:    deposit.LogIndex,
		BlockNumber: deposit.BlockNumber,
		Timestamp:   deposit.Time.Unix(),
	}
	if err = c.Store.InsertEthereumDeposit(r.Context(), storeDeposit); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return store.EthereumDeposit{}, err
	}

	return storeDeposit, nil
}

func (c *StellarWithdrawalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	deposit, err := c.getDeposit(r, w)
	if err != nil {
		return
	}

	// Check withdraw expiration
	lastLedgerCloseTime, err := c.Store.GetLastLedgerCloseTime(r.Context())
	if err != nil {
		log.WithField("error", err).Error("Error getting last ledger close time")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if lastLedgerCloseTime.After(time.Unix(deposit.Timestamp, 0).Add(c.WithdrawalWindow)) {
		// TODO send an error msg to the client
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if withdrawal tx was seen without signature request
	exists, err := c.Store.HistoryStellarTransactionExists(r.Context(), deposit.ID)
	if err != nil {
		log.WithField("error", err).Error("Error getting transaction by memo hash")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		log.WithField("depositID", deposit.ID).Info("withdrawal transaction was already executed")
		// TODO send an error msg to the client
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if outgoing transaction exists
	outgoingTransaction, err := c.Store.GetOutgoingStellarTransaction(r.Context(), store.Withdraw, deposit.ID)
	if err != nil && err != sql.ErrNoRows {
		log.WithField("error", err).Error("Error getting an outgoing stellar transaction for ethereum")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err == nil {
		sourceAccount, err := c.StellarClient.AccountDetail(horizonclient.AccountRequest{
			AccountID: deposit.Destination,
		})
		if err != nil {
			log.WithField("error", err).Error("Error getting stellar account")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if sourceAccount.Sequence < outgoingTransaction.Sequence {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(outgoingTransaction.Envelope))
			return
		}
	}

	// Outgoing Stellar transaction does not exist so create signature request.
	// Duplicate requests for the same signatures are not allowed but the error is ignored.
	err = c.Store.InsertSignatureRequest(r.Context(), store.SignatureRequest{
		DepositChain: store.Ethereum,
		Action:       store.Withdraw,
		DepositID:    deposit.ID,
	})
	if err != nil {
		// Ignore duplicate violations
		if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			log.WithField("error", err).Error("Error inserting a signature request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

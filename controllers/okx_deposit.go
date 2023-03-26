package controllers

import (
	"database/sql"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/backend"
	"math/big"
	"net/http"
	"strconv"

	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/store"
)

var (
	InvalidOkxTxHash = problem.P{
		Type:   "invalid_okx_tx_hash",
		Title:  "Invalid Okx Transaction Hash",
		Status: http.StatusBadRequest,
		Detail: "The transaction hash of the Okx transaction is invalid.",
	}
	OkxTxHashNotFound = problem.P{
		Type:   "okx_tx_hash_not_found",
		Title:  "Okx Transaction Hash Not Found",
		Status: http.StatusNotFound,
		Detail: "The Okx transaction cannot be found.",
	}
	OkxTxRequiresMoreConfirmations = problem.P{
		Type:   "okx_tx_requires_more_confirmations",
		Title:  "Okx Transaction Requires More confirmations",
		Status: http.StatusUnprocessableEntity,
		Detail: "Retry later once the transaction has more confirmations.",
	}
)

type OkxDepositHandler struct {
	Observer                   ethereum.Observer
	Store                      *store.DB
	StellarWithdrawalValidator backend.StellarWithdrawalValidator
	OkxFinalityBuffer          uint64
	Token                      string
}

func (c *OkxDepositHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	txHash := r.PostFormValue("hash")
	stellarAddress := r.PostFormValue("stellar_address")
	parsed, err := strconv.ParseInt(r.PostFormValue("log_index"), 10, 32)
	if err != nil {
		problem.Render(r.Context(), w, InvalidLogIndex)
		return
	}
	logIndex := uint(parsed)

	decoded, err := strkey.Decode(strkey.VersionByteAccountID, stellarAddress)
	if err != nil {
		log.WithField("error", err).Error("Error strkey.Decode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var intEncoded big.Int
	intEncoded.SetBytes(decoded)

	depositID := ethereum.DepositID(txHash, logIndex)

	deposit, err := c.Observer.GetDeposit(r.Context(), txHash, logIndex)
	if ethereum.IsInvalidGetDepositRequest(err) {
		problem.Render(r.Context(), w, InvalidDepositLog)
		return
	} else if err == ethereum.ErrTxHashNotFound {
		problem.Render(r.Context(), w, OkxTxHashNotFound)
		return
	} else if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	block, err := c.Observer.GetLatestBlock(r.Context())
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}
	if deposit.BlockNumber+c.OkxFinalityBuffer > block.Number {
		problem.Render(r.Context(), w, OkxTxRequiresMoreConfirmations)
		return
	}

	incomingTx := store.OkxDeposit{
		ID:          depositID,
		Token:       c.Token,
		Sender:      deposit.Sender.String(),
		Hash:        deposit.TxHash.String(),
		LogIndex:    deposit.LogIndex,
		Amount:      deposit.Amount.String(),
		Destination: intEncoded.String(),
		BlockNumber: deposit.BlockNumber,
		BlockTime:   deposit.Time.Unix(),
	}

	err = c.Store.InsertOkxDeposit(r.Context(), incomingTx)
	if err != nil {
		log.WithField("error", err).Error("Error inserting incoming Okx transaction")
		problem.Render(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(incomingTx.Hash))
}

func getOkxDeposit(observer ethereum.Observer, depositStore *store.DB, finalityBuffer uint64, r *http.Request) (store.OkxDeposit, error) {
	txHash := r.PostFormValue("transaction_hash")
	if !validTxHash.MatchString(txHash) {
		return store.OkxDeposit{}, InvalidOkxTxHash
	}
	parsed, err := strconv.ParseInt(r.PostFormValue("log_index"), 10, 32)
	if err != nil {
		return store.OkxDeposit{}, InvalidLogIndex
	}
	logIndex := uint(parsed)
	depositID := ethereum.DepositID(txHash, logIndex)

	log.Info("depositID")
	log.Info(depositID)
	storeDeposit, err := depositStore.GetOkxDeposit(r.Context(), depositID)
	log.Info("storeDeposit")
	log.Info(storeDeposit)
	if err == nil {
		return storeDeposit, nil
	} else if err != sql.ErrNoRows {
		return store.OkxDeposit{}, err
	}

	deposit, err := observer.GetDeposit(r.Context(), txHash, logIndex)
	if ethereum.IsInvalidGetDepositRequest(err) {
		return store.OkxDeposit{}, InvalidDepositLog
	} else if err == ethereum.ErrTxHashNotFound {
		return store.OkxDeposit{}, OkxTxHashNotFound
	} else if err != nil {
		return store.OkxDeposit{}, err
	}

	block, err := observer.GetLatestBlock(r.Context())
	if err != nil {
		return store.OkxDeposit{}, err
	}
	if deposit.BlockNumber+finalityBuffer > block.Number {
		return store.OkxDeposit{}, OkxTxRequiresMoreConfirmations
	}

	storeDeposit = store.OkxDeposit{
		ID:          depositID,
		Token:       deposit.Token.String(),
		Sender:      deposit.Sender.String(),
		Destination: deposit.Destination.String(),
		Amount:      deposit.Amount.String(),
		Hash:        deposit.TxHash.String(),
		LogIndex:    deposit.LogIndex,
		BlockNumber: deposit.BlockNumber,
		BlockTime:   deposit.Time.Unix(),
	}
	if err = depositStore.InsertOkxDeposit(r.Context(), storeDeposit); err != nil {
		return store.OkxDeposit{}, err
	}

	return storeDeposit, nil
}

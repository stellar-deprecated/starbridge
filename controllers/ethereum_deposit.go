package controllers

import (
	"database/sql"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/backend"
	"math/big"
	"net/http"
	"regexp"
	"strconv"

	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/store"
)

var (
	InvalidEthereumTxHash = problem.P{
		Type:   "invalid_ethereum_tx_hash",
		Title:  "Invalid Ethereum Transaction Hash",
		Status: http.StatusBadRequest,
		Detail: "The transaction hash of the ethereum transaction is invalid.",
	}
	InvalidLogIndex = problem.P{
		Type:   "invalid_log_index",
		Title:  "Invalid Log Index",
		Status: http.StatusBadRequest,
		Detail: "The given log index for the Ethereum deposit is invalid.",
	}
	InvalidDepositLog = problem.P{
		Type:   "invalid_deposit_log",
		Title:  "Invalid Deposit Log",
		Status: http.StatusBadRequest,
		Detail: "The given log for the transaction hash is invalid.",
	}
	EthereumTxHashNotFound = problem.P{
		Type:   "ethereum_tx_hash_not_found",
		Title:  "Ethereum Transaction Hash Not Found",
		Status: http.StatusNotFound,
		Detail: "The ethereum transaction cannot be found.",
	}
	EthereumTxRequiresMoreConfirmations = problem.P{
		Type:   "ethereum_tx_requires_more_confirmations",
		Title:  "Ethereum Transaction Requires More confirmations",
		Status: http.StatusUnprocessableEntity,
		Detail: "Retry later once the transaction has more confirmations.",
	}

	validTxHash = regexp.MustCompile("^(0x)?([A-Fa-f0-9]{64})$")
)

type EthereumDeposit struct {
	Observer                   ethereum.Observer
	Store                      *store.DB
	StellarWithdrawalValidator backend.StellarWithdrawalValidator
	EthereumFinalityBuffer     uint64
	Token                      string
}

func (c *EthereumDeposit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		problem.Render(r.Context(), w, EthereumTxHashNotFound)
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
	if deposit.BlockNumber+c.EthereumFinalityBuffer > block.Number {
		problem.Render(r.Context(), w, EthereumTxRequiresMoreConfirmations)
		return
	}

	incomingTx := store.EthereumDeposit{
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

	err = c.Store.InsertEthereumDeposit(r.Context(), incomingTx)
	if err != nil {
		log.WithField("error", err).Error("Error inserting incoming ethereum transaction")
		problem.Render(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(incomingTx.Hash))
}

func getEthereumDeposit(observer ethereum.Observer, depositStore *store.DB, finalityBuffer uint64, r *http.Request) (store.EthereumDeposit, error) {
	txHash := r.PostFormValue("transaction_hash")
	if !validTxHash.MatchString(txHash) {
		return store.EthereumDeposit{}, InvalidEthereumTxHash
	}
	parsed, err := strconv.ParseInt(r.PostFormValue("log_index"), 10, 32)
	if err != nil {
		return store.EthereumDeposit{}, InvalidLogIndex
	}
	logIndex := uint(parsed)
	depositID := ethereum.DepositID(txHash, logIndex)

	storeDeposit, err := depositStore.GetEthereumDeposit(r.Context(), depositID)
	if err == nil {
		return storeDeposit, nil
	} else if err != sql.ErrNoRows {
		return store.EthereumDeposit{}, err
	}

	deposit, err := observer.GetDeposit(r.Context(), txHash, logIndex)
	if ethereum.IsInvalidGetDepositRequest(err) {
		return store.EthereumDeposit{}, InvalidDepositLog
	} else if err == ethereum.ErrTxHashNotFound {
		return store.EthereumDeposit{}, EthereumTxHashNotFound
	} else if err != nil {
		return store.EthereumDeposit{}, err
	}

	block, err := observer.GetLatestBlock(r.Context())
	if err != nil {
		return store.EthereumDeposit{}, err
	}
	if deposit.BlockNumber+finalityBuffer > block.Number {
		return store.EthereumDeposit{}, EthereumTxRequiresMoreConfirmations
	}

	storeDeposit = store.EthereumDeposit{
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
	if err = depositStore.InsertEthereumDeposit(r.Context(), storeDeposit); err != nil {
		return store.EthereumDeposit{}, err
	}

	return storeDeposit, nil
}

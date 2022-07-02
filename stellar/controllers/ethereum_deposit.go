package controllers

import (
	"database/sql"
	"net/http"
	"regexp"
	"strconv"

	"github.com/stellar/go/strkey"
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
	InvalidStellarRecipient = problem.P{
		Type:   "invalid_stellar_recipient",
		Title:  "Invalid Stellar Recipient",
		Status: http.StatusUnprocessableEntity,
		Detail: "The recipient of the deposit is not a valid Stellar address.",
	}

	validTxHash = regexp.MustCompile("^(0x)?([A-Fa-f0-9]{64})$")
)

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

	destinationAccountID, err := strkey.Encode(
		strkey.VersionByteAccountID,
		deposit.Destination.Bytes(),
	)
	if err != nil {
		return store.EthereumDeposit{}, InvalidStellarRecipient
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
	if err = depositStore.InsertEthereumDeposit(r.Context(), storeDeposit); err != nil {
		return store.EthereumDeposit{}, err
	}

	return storeDeposit, nil
}

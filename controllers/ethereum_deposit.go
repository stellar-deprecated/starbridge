package controllers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/ethereum"
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

func getEthereumDeposit(observer ethereum.Observer, finalityBuffer uint64, r *http.Request) (ethereum.Deposit, error) {
	txHash := r.PostFormValue("transaction_hash")
	if !validTxHash.MatchString(txHash) {
		return ethereum.Deposit{}, InvalidEthereumTxHash
	}
	parsed, err := strconv.ParseInt(r.PostFormValue("log_index"), 10, 32)
	if err != nil {
		return ethereum.Deposit{}, InvalidLogIndex
	}
	logIndex := uint(parsed)

	deposit, err := observer.GetDeposit(r.Context(), txHash, logIndex)
	if ethereum.IsInvalidGetDepositRequest(err) {
		return ethereum.Deposit{}, InvalidDepositLog
	} else if err == ethereum.ErrTxHashNotFound {
		return ethereum.Deposit{}, EthereumTxHashNotFound
	} else if err != nil {
		return ethereum.Deposit{}, err
	}

	block, err := observer.GetLatestBlock(r.Context())
	if err != nil {
		return ethereum.Deposit{}, err
	}
	if deposit.BlockNumber+finalityBuffer > block.Number {
		return ethereum.Deposit{}, EthereumTxRequiresMoreConfirmations
	}

	return deposit, nil
}

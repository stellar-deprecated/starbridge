package controllers

import (
	"context"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/concordium"
	"github.com/stellar/starbridge/store"
	"math/big"
	"net/http"
)

var (
//InvalidEthereumTxHash = problem.P{
//	Type:   "invalid_ethereum_tx_hash",
//	Title:  "Invalid Ethereum Transaction Hash",
//	Status: http.StatusBadRequest,
//	Detail: "The transaction hash of the ethereum transaction is invalid.",
//}
//InvalidLogIndex = problem.P{
//	Type:   "invalid_log_index",
//	Title:  "Invalid Log Index",
//	Status: http.StatusBadRequest,
//	Detail: "The given log index for the Ethereum deposit is invalid.",
//}
//InvalidDepositLog = problem.P{
//	Type:   "invalid_deposit_log",
//	Title:  "Invalid Deposit Log",
//	Status: http.StatusBadRequest,
//	Detail: "The given log for the transaction hash is invalid.",
//}
//EthereumTxHashNotFound = problem.P{
//	Type:   "ethereum_tx_hash_not_found",
//	Title:  "Ethereum Transaction Hash Not Found",
//	Status: http.StatusNotFound,
//	Detail: "The ethereum transaction cannot be found.",
//}
//EthereumTxRequiresMoreConfirmations = problem.P{
//	Type:   "ethereum_tx_requires_more_confirmations",
//	Title:  "Ethereum Transaction Requires More confirmations",
//	Status: http.StatusUnprocessableEntity,
//	Detail: "Retry later once the transaction has more confirmations.",
//}
//
//validTxHash = regexp.MustCompile("^(0x)?([A-Fa-f0-9]{64})$")
)

type ConcordiumDepositHandler struct {
	Observer concordium.Observer
	Store    *store.DB
}

func (c *ConcordiumDepositHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	txHash := r.PostFormValue("hash")
	stellarAddress := r.PostFormValue("stellar_address")
	decoded, err := strkey.Decode(strkey.VersionByteAccountID, stellarAddress)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	var intEncoded big.Int
	intEncoded.SetBytes(decoded)
	deposit, err := c.Observer.GetDeposit(r.Context(), txHash)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}
	incomingTx := store.ConcordiumDeposit{
		ID:          deposit.Hash,
		Sender:      deposit.From,
		Amount:      deposit.Amount,
		Destination: intEncoded.String(),
		BlockHash:   deposit.BlockHash,
		BlockTime:   deposit.BlockTime,
	}
	if err = c.Store.InsertConcordiumDeposit(r.Context(), incomingTx); err != nil {
		log.WithField("error", err).Error("Error inserting incoming concordium transaction")
		problem.Render(r.Context(), w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(deposit.Hash))

	//depositID := ethereum.DepositID(txHash, logIndex)
	//
	//deposit, err := c.Observer.GetDeposit(r.Context(), txHash, logIndex)
	//if ethereum.IsInvalidGetDepositRequest(err) {
	//	problem.Render(r.Context(), w, InvalidDepositLog)
	//	return
	//} else if err == ethereum.ErrTxHashNotFound {
	//	problem.Render(r.Context(), w, EthereumTxHashNotFound)
	//	return
	//} else if err != nil {
	//	problem.Render(r.Context(), w, err)
	//	return
	//}
	//
	//block, err := c.Observer.GetLatestBlock(r.Context())
	//if err != nil {
	//	problem.Render(r.Context(), w, err)
	//	return
	//}
	//if deposit.BlockNumber+c.EthereumFinalityBuffer > block.Number {
	//	problem.Render(r.Context(), w, EthereumTxRequiresMoreConfirmations)
	//	return
	//}
	//
	//incomingTx := store.EthereumDeposit{
	//	ID:          depositID,
	//	Token:       c.Token,
	//	Sender:      deposit.Sender.String(),
	//	Hash:        deposit.TxHash.String(),
	//	LogIndex:    deposit.LogIndex,
	//	Amount:      deposit.Amount.String(),
	//	Destination: intEncoded.String(),
	//	BlockNumber: deposit.BlockNumber,
	//	BlockTime:   deposit.Time.Unix(),
	//}
	//
	//err = c.Store.InsertEthereumDeposit(r.Context(), incomingTx)
	//if err != nil {
	//	log.WithField("error", err).Error("Error inserting incoming ethereum transaction")
	//	problem.Render(r.Context(), w, err)
	//	return
	//}
	//
	//w.WriteHeader(http.StatusOK)
	//_, _ = w.Write([]byte(incomingTx.Hash))
}

func getConcordiumDeposit(txHash string, depositStore *store.DB) (store.ConcordiumDeposit, error) {
	storeDeposit, err := depositStore.GetConcordiumDeposit(context.Background(), txHash)
	if err != nil {
		return store.ConcordiumDeposit{}, err
	}
	return storeDeposit, nil
}

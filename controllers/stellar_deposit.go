package controllers

import (
	"database/sql"
	"math/big"
	"net/http"
	"strings"

	"github.com/stellar/go/strkey"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/store"
)

var (
	InvalidStellarTxHash = problem.P{
		Type:   "invalid_stellar_tx_hash",
		Title:  "Invalid Stellar Transaction Hash",
		Status: http.StatusBadRequest,
		Detail: "The transaction hash of the Stellar transaction is invalid.",
	}
	StellarTxHashNotFound = problem.P{
		Type:   "stellar_tx_hash_not_found",
		Title:  "Stellar Transaction Hash Not Found",
		Status: http.StatusNotFound,
		Detail: "The stellar transaction cannot be found.",
	}
)

func getStellarDeposit(depositStore *store.DB, r *http.Request) (store.StellarDeposit, error) {
	txHash := strings.TrimPrefix(r.PostFormValue("transaction_hash"), "0x")
	if !validTxHash.MatchString(txHash) {
		return store.StellarDeposit{}, InvalidStellarTxHash
	}

	deposit, err := depositStore.GetStellarDeposit(r.Context(), txHash)
	if err == sql.ErrNoRows {
		return store.StellarDeposit{}, StellarTxHashNotFound
	}
	return deposit, err
}

// TODO remove after prototype demo
type TestDeposit struct {
	Store *store.DB
	Token string
}

func (c *TestDeposit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hash := r.PostFormValue("hash")
	stellarAddress := r.PostFormValue("stellar_address")

	decoded, err := strkey.Decode(strkey.VersionByteAccountID, stellarAddress)
	if err != nil {
		log.WithField("error", err).Error("Error strkey.Decode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var intEncoded big.Int
	intEncoded.SetBytes(decoded)

	incomingTx := store.EthereumDeposit{
		ID:          ethereum.DepositID(hash, 1),
		Token:       c.Token,
		Hash:        hash,
		LogIndex:    1,
		Amount:      "1",
		Destination: intEncoded.String(),
	}

	err = c.Store.InsertEthereumDeposit(r.Context(), incomingTx)
	if err != nil {
		log.WithField("error", err).Error("Error inserting incoming ethereum transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(incomingTx.Hash))
}

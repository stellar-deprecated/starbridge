package store

import (
	"context"
	"math/big"
	"time"
)

type IncomingEthereumTransaction struct {
	Hash               string
	ValueWei           *big.Int
	StellarAddress     string
	WithdrawExpiration time.Time

	TransactionBlob string
}

func (m *DB) GetIncomingEthereumTransactionByHash(ctx context.Context, hash string) (IncomingEthereumTransaction, error) {
	return IncomingEthereumTransaction{
		Hash:           hash,
		ValueWei:       big.NewInt(627836782638726),
		StellarAddress: "GD6DBPOQ5EFMDEJ6X2PTRTKM2ZNGJCINC3U3BE7BHMM3C6D75JDLP2KX",
	}, nil
}

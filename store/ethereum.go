package store

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type IncomingEthereumTransaction struct {
	Hash               string `db:"hash"`
	ValueWei           int64  `db:"value_wei"` // TODO change to big.Int
	StellarAddress     string `db:"stellar_address"`
	WithdrawExpiration time.Time
}

func (m *DB) GetIncomingEthereumTransactionByHash(ctx context.Context, hash string) (IncomingEthereumTransaction, error) {
	sql := sq.Select("*").From("incoming_ethereum_transactions").Where(sq.Eq{"hash": hash})

	var result IncomingEthereumTransaction
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) InsertIncomingEthereumTransaction(ctx context.Context, newtx IncomingEthereumTransaction) error {
	query := sq.Insert("incoming_ethereum_transactions").
		SetMap(map[string]interface{}{
			"hash":            newtx.Hash,
			"value_wei":       50,
			"stellar_address": newtx.StellarAddress,
		})

	_, err := m.Session.Exec(ctx, query)
	return err
}

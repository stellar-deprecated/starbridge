package store

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type IncomingEthereumTransaction struct {
	Hash               string    `db:"hash"`
	ValueWei           string    `db:"value_wei"`
	StellarAddress     string    `db:"stellar_address"`
	WithdrawExpiration time.Time `db:"withdraw_expiration"`
	Withdrawn          bool      `db:"withdrawn"`
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
			"hash":                newtx.Hash,
			"value_wei":           newtx.ValueWei,
			"stellar_address":     newtx.StellarAddress,
			"withdraw_expiration": newtx.WithdrawExpiration,
			"withdrawn":           false,
		})

	_, err := m.Session.Exec(ctx, query)
	return err
}

func (m *DB) MarkIncomingEthereumTransactionAsWithdrawn(ctx context.Context, hash string) error {
	query := sq.Update("incoming_ethereum_transactions").
		Set("withdrawn", true).
		Where(sq.Eq{"hash": hash})

	_, err := m.Session.Exec(ctx, query)
	return err
}

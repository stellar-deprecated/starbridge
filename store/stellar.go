package store

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type HistoryStellarTransaction struct {
	Hash     string `db:"hash"`
	Envelope string `db:"envelope"`
	MemoHash string `db:"memo_hash"`
}

type OutgoingStellarTransaction struct {
	Hash      string `db:"hash"`
	Envelope  string `db:"envelope"`
	Sequence  int64  `db:"sequence"`
	Action    Action `db:"requested_action"`
	DepositID string `db:"deposit_id"`
}

func (m *DB) InsertHistoryStellarTransaction(ctx context.Context, tx HistoryStellarTransaction) error {
	query := sq.Insert("history_stellar_transactions").
		SetMap(map[string]interface{}{
			"hash":      tx.Hash,
			"envelope":  tx.Envelope,
			"memo_hash": tx.MemoHash,
		})

	_, err := m.Session.Exec(ctx, query)
	return err
}

func (m *DB) HistoryStellarTransactionExists(ctx context.Context, memoHash string) (bool, error) {
	stmt := sq.Select("1").From("history_stellar_transactions").
		Where(sq.Eq{"memo_hash": memoHash})

	var result int
	err := m.Session.Get(ctx, &result, stmt)
	if err == nil {
		return true, nil
	} else if err == sql.ErrNoRows {
		return false, nil
	}
	return false, err
}

func (m *DB) GetOutgoingStellarTransaction(ctx context.Context, action Action, depositID string) (OutgoingStellarTransaction, error) {
	sql := sq.Select("*").From("outgoing_stellar_transactions").Where(map[string]interface{}{
		"requested_action": action,
		"deposit_id":       depositID,
	})

	var result OutgoingStellarTransaction
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) UpsertOutgoingStellarTransaction(ctx context.Context, newtx OutgoingStellarTransaction) error {
	query := sq.Insert("outgoing_stellar_transactions").
		SetMap(map[string]interface{}{
			"hash":             newtx.Hash,
			"envelope":         newtx.Envelope,
			"requested_action": newtx.Action,
			"deposit_id":       newtx.DepositID,
			"sequence":         newtx.Sequence,
		}).
		Suffix("ON CONFLICT (requested_action, deposit_id) " +
			"DO UPDATE SET " +
			"sequence=EXCLUDED.sequence, hash=EXCLUDED.hash, envelope=EXCLUDED.envelope",
		)

	_, err := m.Session.Exec(ctx, query)
	return err
}

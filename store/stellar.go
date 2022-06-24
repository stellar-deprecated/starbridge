package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

type OutgoingStellarTransactionState string

const (
	PendingState OutgoingStellarTransactionState = "pending"
	SuccessState OutgoingStellarTransactionState = "success"
	// InvalidState is either expired, failed or irrevocably invalid tx
	InvalidState OutgoingStellarTransactionState = "invalid"
)

type HistoryStellarTransaction struct {
	Hash     string `db:"hash"`
	Envelope string `db:"envelope"`
	MemoHash string `db:"memo_hash"`
}

type OutgoingStellarTransaction struct {
	Hash     string                          `db:"hash"`
	State    OutgoingStellarTransactionState `db:"state"`
	Envelope string                          `db:"envelope"`

	IncomingType            NetworkType `db:"incoming_type"`
	IncomingTransactionHash string      `db:"incoming_transaction_hash"`
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

func (m *DB) GetHistoryStellarTransactionByMemoHash(ctx context.Context, memoHash string) (HistoryStellarTransaction, error) {
	sql := sq.Select("*").From("history_stellar_transactions").
		Where(sq.Eq{"memo_hash": memoHash})

	var result HistoryStellarTransaction
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return HistoryStellarTransaction{}, err
	}

	return result, nil
}

func (m *DB) GetPendingOutgoingStellarTransactions(ctx context.Context) ([]OutgoingStellarTransaction, error) {
	sql := sq.Select("*").From("outgoing_stellar_transactions").
		Where(sq.Eq{"state": PendingState})

	var results []OutgoingStellarTransaction
	if err := m.Session.Select(ctx, &results, sql); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *DB) GetOutgoingStellarTransactionForEthereumByHash(ctx context.Context, hash string) (OutgoingStellarTransaction, error) {
	sql := sq.Select("*").From("outgoing_stellar_transactions").Where(map[string]interface{}{
		"incoming_type":             Ethereum,
		"incoming_transaction_hash": hash,
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
			"state":                     newtx.State,
			"hash":                      newtx.Hash,
			"envelope":                  newtx.Envelope,
			"incoming_type":             newtx.IncomingType,
			"incoming_transaction_hash": newtx.IncomingTransactionHash,
		}).
		Suffix("ON CONFLICT (incoming_type,incoming_transaction_hash) DO UPDATE SET state=EXCLUDED.state, hash=EXCLUDED.hash, envelope=EXCLUDED.envelope")

	_, err := m.Session.Exec(ctx, query)
	return err
}

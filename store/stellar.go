package store

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type OutgoingStellarTransactionState string

const (
	PendingState OutgoingStellarTransactionState = "pending"
	SuccessState OutgoingStellarTransactionState = "success"
	FailedState  OutgoingStellarTransactionState = "failed"
	ExpiredState OutgoingStellarTransactionState = "expired"
)

type OutgoingStellarTransaction struct {
	State      OutgoingStellarTransactionState `db:"state"`
	Hash       string                          `db:"hash"`
	Envelope   string                          `db:"envelope"`
	Expiration time.Time                       `db:"expiration"`

	IncomingType                    NetworkType `db:"incoming_type"`
	IncomingEthereumTransactionHash *string     `db:"incoming_ethereum_transaction_hash"`
}

// TODO: this should select loaded transactions for update so other go routines wait
// but will be fixed in another PR by running worker and observer in the same go routine.
func (m *DB) GetOutgoingStellarTransactions(ctx context.Context) ([]OutgoingStellarTransaction, error) {
	sql := sq.Select("*").From("outgoing_stellar_transactions")

	var results []OutgoingStellarTransaction
	if err := m.Session.Select(ctx, &results, sql); err != nil {
		return nil, err
	}

	return results, nil
}

// TODO: this should select loaded transactions for update so other go routines wait
// but will be fixed in another PR by running worker and observer in the same go routine.
func (m *DB) GetOutgoingStellarTransactionForEthereumByHash(ctx context.Context, hash string) (OutgoingStellarTransaction, error) {
	sql := sq.Select("*").From("outgoing_stellar_transactions").Where(map[string]interface{}{
		"incoming_type":                      Ethereum,
		"incoming_ethereum_transaction_hash": hash,
	})

	var result OutgoingStellarTransaction
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) MarkOutgoingStellarTransactionExpired(ctx context.Context, expiredBefore time.Time) (int64, error) {
	sql := sq.Update("outgoing_stellar_transactions").
		Set("state", ExpiredState).
		Where(sq.NotEq{"state": ExpiredState}).
		Where(sq.Lt{"expiration": expiredBefore})
	result, err := m.Session.Exec(ctx, sql)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (m *DB) UpsertOutgoingStellarTransaction(ctx context.Context, newtx OutgoingStellarTransaction) error {
	query := sq.Insert("outgoing_stellar_transactions").
		SetMap(map[string]interface{}{
			"state":                              newtx.State,
			"hash":                               newtx.Hash,
			"envelope":                           newtx.Envelope,
			"expiration":                         newtx.Expiration,
			"incoming_type":                      newtx.IncomingType,
			"incoming_ethereum_transaction_hash": newtx.IncomingEthereumTransactionHash,
		}).
		Suffix("ON CONFLICT (hash) DO UPDATE SET state=EXCLUDED.state")

	_, err := m.Session.Exec(ctx, query)
	return err
}

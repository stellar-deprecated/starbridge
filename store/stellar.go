package store

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type StellarDeposit struct {
	// ID is the globally unique id for this deposit
	// and is equal to deposit transaction hash
	ID string `db:"id"`
	// Asset is the string encoding of the Stellar assets
	// which were deposited to the bridge
	Asset string `db:"asset"`
	// Sender is the address of the account which deposited the tokens
	Sender string `db:"sender"`
	// Destination is the intended recipient of the bridge transfer
	Destination string `db:"destination"`
	// Amount is the amount of tokens which were deposited to the bridge
	// contract
	Amount string `db:"amount"`
	// LedgerTime is the unix timestamp of the deposit
	LedgerTime int64 `db:"ledger_time"`
}

type HistoryStellarTransaction struct {
	Hash     string `db:"hash"`
	Envelope string `db:"envelope"`
	// MemoHash represents:
	//   - Ethereum deposit ID in case of withdrawals in Ethereum->Stellar flow
	//   - Stellar transaction hash in case of refunds in Stellar->Ethereum flow
	MemoHash string `db:"memo_hash"`
}

type OutgoingStellarTransaction struct {
	Envelope      string `db:"envelope"`
	SourceAccount string `db:"source_account"`
	Sequence      int64  `db:"sequence"`
	Action        Action `db:"requested_action"`
	DepositID     string `db:"deposit_id"`
}

func (m *DB) GetStellarDeposit(ctx context.Context, id string) (StellarDeposit, error) {
	sql := sq.Select("*").From("stellar_deposits").Where(
		sq.Eq{"id": id},
	)

	var result StellarDeposit
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) InsertStellarDeposit(ctx context.Context, deposit StellarDeposit) error {
	query := sq.Insert("stellar_deposits").
		SetMap(map[string]interface{}{
			"id":          deposit.ID,
			"ledger_time": deposit.LedgerTime,
			"amount":      deposit.Amount,
			"destination": deposit.Destination,
			"sender":      deposit.Sender,
			"asset":       deposit.Asset,
		})

	_, err := m.Session.Exec(ctx, query)
	return err
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
			"envelope":         newtx.Envelope,
			"requested_action": newtx.Action,
			"deposit_id":       newtx.DepositID,
			"sequence":         newtx.Sequence,
			"source_account":   newtx.SourceAccount,
		}).
		Suffix("ON CONFLICT (requested_action, deposit_id) " +
			"DO UPDATE SET " +
			"sequence=EXCLUDED.sequence, source_account=EXCLUDED.source_account, envelope=EXCLUDED.envelope",
		)

	_, err := m.Session.Exec(ctx, query)
	return err
}

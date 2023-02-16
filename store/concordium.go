package store

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type ConcordiumDeposit struct {
	ID          string `db:"id"`
	Sender      string `db:"sender"`
	Destination string `db:"destination"`
	Amount      string `db:"amount"`
	BlockHash   string `db:"block_hash"`
	BlockTime   int64  `db:"block_time"`
}

type ConcordiumSignature struct {
	Address    string `db:"address"`
	Token      string `db:"token"`
	Amount     string `db:"amount"`
	Signature  string `db:"signature"`
	Expiration int64  `db:"expiration"`
	Action     Action `db:"requested_action"`
	DepositID  string `db:"deposit_id"`
}

func (m *DB) GetConcordiumDeposit(ctx context.Context, id string) (ConcordiumDeposit, error) {
	sql := sq.Select("*").From("concordium_deposits").Where(
		sq.Eq{"id": strings.ToLower(id)},
	)

	var result ConcordiumDeposit
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) InsertConcordiumDeposit(ctx context.Context, deposit ConcordiumDeposit) error {
	query := sq.Insert("concordium_deposits").
		SetMap(map[string]interface{}{
			"id":          strings.ToLower(deposit.ID),
			"amount":      deposit.Amount,
			"destination": deposit.Destination,
			"sender":      deposit.Sender,
			"block_hash":  deposit.BlockHash,
			"block_time":  deposit.BlockTime,
		})

	_, err := m.Session.Exec(ctx, query)
	return err
}

func (m *DB) GetConcordiumSignature(ctx context.Context, action Action, depositID string) (ConcordiumSignature, error) {
	sql := sq.Select("*").From("concordium_signatures").Where(map[string]interface{}{
		"requested_action": action,
		"deposit_id":       strings.ToLower(depositID),
	})

	var result ConcordiumSignature
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) UpsertConcordiumSignature(ctx context.Context, newSig ConcordiumSignature) error {
	query := sq.Insert("concordium_signatures").
		SetMap(map[string]interface{}{
			"address":          newSig.Address,
			"signature":        newSig.Signature,
			"expiration":       newSig.Expiration,
			"requested_action": newSig.Action,
			"deposit_id":       strings.ToLower(newSig.DepositID),
			"token":            newSig.Token,
			"amount":           newSig.Amount,
		}).
		Suffix("ON CONFLICT (requested_action, deposit_id) " +
			"DO UPDATE SET " +
			"signature=EXCLUDED.signature, address=EXCLUDED.address, " +
			"expiration=EXCLUDED.expiration, token=EXCLUDED.token, amount=EXCLUDED.amount",
		)

	_, err := m.Session.Exec(ctx, query)
	return err
}

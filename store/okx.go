package store

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type OkxDeposit struct {
	// ID is the globally unique id for this deposit
	ID string `db:"id"`
	// Token is the address (0x0 in the case that okx was deposited)
	// of the tokens which were deposited to the bridge
	Token string `db:"token"`
	// Sender is the address of the account which deposited the tokens
	Sender string `db:"sender"`
	// Destination is the intended recipient of the bridge transfer
	Destination string `db:"destination"`
	// Amount is the amount of tokens which were deposited to the bridge
	// contract
	Amount string `db:"amount"`
	// Hash is the hash of the transaction containing the deposit
	Hash string `db:"hash"`
	// LogIndex is the log index within the okx block of the deposit event
	// emitted by the bridge contract
	LogIndex uint `db:"log_index"`
	// BlockNumber is the sequence number of the block containing the deposit
	// transaction
	BlockNumber uint64 `db:"block_number"`
	// BlockTime is the unix timestamp of the deposit
	BlockTime int64 `db:"block_time"`
}

// OkxSignature represents a signature for a withdrawal / refund against the
// bridge okx smart contract
type OkxSignature struct {
	Address    string `db:"address"`
	Token      string `db:"token"`
	Amount     string `db:"amount"`
	Signature  string `db:"signature"`
	Expiration int64  `db:"expiration"`
	Action     Action `db:"requested_action"`
	DepositID  string `db:"deposit_id"`
}

func (m *DB) GetOkxDeposit(ctx context.Context, id string) (OkxDeposit, error) {
	sql := sq.Select("*").From("okx_deposits").Where(
		sq.Eq{"id": strings.ToLower(id)},
	)

	var result OkxDeposit
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) InsertOkxDeposit(ctx context.Context, deposit OkxDeposit) error {
	query := sq.Insert("okx_deposits").
		SetMap(map[string]interface{}{
			"id":           strings.ToLower(deposit.ID),
			"hash":         deposit.Hash,
			"log_index":    deposit.LogIndex,
			"block_number": deposit.BlockNumber,
			"block_time":   deposit.BlockTime,
			"amount":       deposit.Amount,
			"destination":  deposit.Destination,
			"sender":       deposit.Sender,
			"token":        deposit.Token,
		})

	_, err := m.Session.Exec(ctx, query)
	return err
}

func (m *DB) GetOkxSignature(ctx context.Context, action Action, depositID string) (OkxSignature, error) {
	sql := sq.Select("*").From("okx_signatures").Where(map[string]interface{}{
		"requested_action": action,
		"deposit_id":       strings.ToLower(depositID),
	})

	var result OkxSignature
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) UpsertOkxSignature(ctx context.Context, newSig OkxSignature) error {
	query := sq.Insert("okx_signatures").
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

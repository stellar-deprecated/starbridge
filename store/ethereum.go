package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

type EthereumDeposit struct {
	// ID is the globally unique id for this deposit
	ID string `db:"id"`
	// Token is the address (0x0 in the case that eth was deposited)
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
	// LogIndex is the log index within the ethereum block of the deposit event
	// emitted by the bridge contract
	LogIndex uint `db:"log_index"`
	// BlockNumber is the sequence number of the block containing the deposit
	// transaction
	BlockNumber uint64 `db:"block_number"`
	// BlockTime is the unix timestamp of the deposit
	BlockTime int64 `db:"block_time"`
}

// EthereumSignature represents a signature for a withdrawal / refund against the
// bridge ethereum smart contract
type EthereumSignature struct {
	Address    string `db:"address"`
	Signature  string `db:"signature"`
	Expiration int64  `db:"expiration"`
	Action     Action `db:"requested_action"`
	DepositID  string `db:"deposit_id"`
}

func (m *DB) GetEthereumDeposit(ctx context.Context, id string) (EthereumDeposit, error) {
	sql := sq.Select("*").From("ethereum_deposits").Where(
		sq.Eq{"id": id},
	)

	var result EthereumDeposit
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) InsertEthereumDeposit(ctx context.Context, deposit EthereumDeposit) error {
	query := sq.Insert("ethereum_deposits").
		SetMap(map[string]interface{}{
			"id":           deposit.ID,
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

func (m *DB) GetEthereumSignature(ctx context.Context, action Action, depositID string) (EthereumSignature, error) {
	sql := sq.Select("*").From("ethereum_signatures").Where(map[string]interface{}{
		"requested_action": action,
		"deposit_id":       depositID,
	})

	var result EthereumSignature
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) UpsertEthereumSignature(ctx context.Context, newSig EthereumSignature) error {
	query := sq.Insert("ethereum_signatures").
		SetMap(map[string]interface{}{
			"address":          newSig.Address,
			"signature":        newSig.Signature,
			"expiration":       newSig.Expiration,
			"requested_action": newSig.Action,
			"deposit_id":       newSig.DepositID,
		}).
		Suffix("ON CONFLICT (requested_action, deposit_id) " +
			"DO UPDATE SET " +
			"signature=EXCLUDED.signature, address=EXCLUDED.address, expiration=EXCLUDED.expiration",
		)

	_, err := m.Session.Exec(ctx, query)
	return err
}

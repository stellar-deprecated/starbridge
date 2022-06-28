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
	// Timestamp is the unix timestamp of the deposit
	Timestamp int64 `db:"block_time"`
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
			"block_time":   deposit.Timestamp,
			"amount":       deposit.Amount,
			"destination":  deposit.Destination,
			"sender":       deposit.Sender,
			"token":        deposit.Token,
		})

	_, err := m.Session.Exec(ctx, query)
	return err
}

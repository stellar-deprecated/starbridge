package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

type NetworkType string

const (
	Ethereum NetworkType = "ethereum"
)

type SignatureRequest struct {
	IncomingType            NetworkType `db:"incoming_type"`
	IncomingTransactionHash string      `db:"incoming_transaction_hash"`
	TxExpirationTimestamp   int64       `db:"tx_expiration_timestamp"`
}

func (m *DB) InsertSignatureRequestForIncomingEthereumTransaction(ctx context.Context, hash string, expirationTimestamp int64) error {
	sql := sq.Insert("signature_requests").SetMap(map[string]interface{}{
		"incoming_type":             Ethereum,
		"incoming_transaction_hash": hash,
		`tx_expiration_timestamp`:   expirationTimestamp,
	})
	_, err := m.Session.Exec(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}

func (m *DB) GetSignatureRequests(ctx context.Context) ([]SignatureRequest, error) {
	sql := sq.Select("*").From("signature_requests")

	var results []SignatureRequest
	if err := m.Session.Select(ctx, &results, sql); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *DB) GetSignatureRequestForIncomingEthereumTransaction(ctx context.Context, hash string) (SignatureRequest, error) {
	sql := sq.Select("*").From("signature_requests").Where(map[string]interface{}{
		"incoming_type":             Ethereum,
		"incoming_transaction_hash": hash,
	})

	var result SignatureRequest
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) DeleteSignatureRequestForIncomingEthereumTransaction(ctx context.Context, hash string) error {
	del := sq.Delete("signature_requests").Where(map[string]interface{}{
		"incoming_type":             Ethereum,
		"incoming_transaction_hash": hash,
	})
	_, err := m.Session.Exec(ctx, del)
	return err
}

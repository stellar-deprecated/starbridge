package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/guregu/null"
)

type NetworkType string

const (
	Ethereum NetworkType = "ethereum"
)

type SignatureRequest struct {
	IncomingType                    NetworkType `db:"incoming_type"`
	IncomingEthereumTransactionHash null.String `db:"incoming_ethereum_transaction_hash"`
}

func (m *DB) InsertSignatureRequestForIncomingEthereumTransaction(ctx context.Context, hash string) error {
	sql := sq.Insert("signature_requests").SetMap(map[string]interface{}{
		"incoming_type":                      Ethereum,
		"incoming_ethereum_transaction_hash": hash,
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
		"incoming_type":                      Ethereum,
		"incoming_ethereum_transaction_hash": hash,
	})

	var result SignatureRequest
	if err := m.Session.Get(ctx, &result, sql); err != nil {
		return result, err
	}

	return result, nil
}

func (m *DB) DeleteSignatureRequestForIncomingEthereumTransaction(ctx context.Context, hash string) error {
	del := sq.Delete("signature_requests").Where(map[string]interface{}{
		"incoming_type":                      Ethereum,
		"incoming_ethereum_transaction_hash": hash,
	})
	_, err := m.Session.Exec(ctx, del)
	return err
}

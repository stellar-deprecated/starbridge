package store

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type (
	Blockchain string
	Action     string
)

const (
	Stellar  Blockchain = "stellar"
	Ethereum Blockchain = "ethereum"
)

const (
	Withdraw Action = "withdraw"
	Refund   Action = "refund"
)

type SignatureRequest struct {
	DepositChain Blockchain `db:"deposit_chain"`
	Action       Action     `db:"requested_action"`
	DepositID    string     `db:"deposit_id"`
}

func (m *DB) InsertSignatureRequest(ctx context.Context, request SignatureRequest) error {
	sql := sq.Insert("signature_requests").SetMap(map[string]interface{}{
		"deposit_chain":    request.DepositChain,
		"requested_action": request.Action,
		"deposit_id":       request.DepositID,
	})
	_, err := m.Session.Exec(ctx, sql)
	// Ignore duplicate violations
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
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

func (m *DB) DeleteSignatureRequest(ctx context.Context, request SignatureRequest) error {
	del := sq.Delete("signature_requests").Where(map[string]interface{}{
		"deposit_chain":    request.DepositChain,
		"deposit_id":       request.DepositID,
		"requested_action": request.Action,
	})
	_, err := m.Session.Exec(ctx, del)
	return err
}

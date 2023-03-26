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
	Stellar    Blockchain = "stellar"
	Ethereum   Blockchain = "ethereum"
	Okx        Blockchain = "okx"
	Concordium Blockchain = "concordium"
)

const (
	Withdraw Action = "withdraw"
	Refund   Action = "refund"
)

type SignatureRequest struct {
	WithdrawChain Blockchain `db:"withdraw_chain"`
	DepositChain  Blockchain `db:"deposit_chain"`
	Action        Action     `db:"requested_action"`
	DepositID     string     `db:"deposit_id"`
}

func (m *DB) InsertSignatureRequest(ctx context.Context, request SignatureRequest) error {
	sql := sq.Insert("signature_requests").SetMap(map[string]interface{}{
		"withdraw_chain":   request.WithdrawChain,
		"deposit_chain":    request.DepositChain,
		"requested_action": request.Action,
		"deposit_id":       strings.ToLower(request.DepositID),
	})
	_, err := m.Session.Exec(ctx, sql)
	// Ignore duplicate violations
	if err != nil && !IsDuplicateError(err) {
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
		"withdraw_chain":   request.WithdrawChain,
		"deposit_chain":    request.DepositChain,
		"deposit_id":       strings.ToLower(request.DepositID),
		"requested_action": request.Action,
	})
	_, err := m.Session.Exec(ctx, del)
	return err
}

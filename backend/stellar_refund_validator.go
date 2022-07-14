package backend

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stellar/go/support/db"
	"github.com/stellar/go/support/render/problem"

	"github.com/stellar/go/support/errors"
	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/store"
)

var RefundAlreadyExecuted = problem.P{
	Type:   "refund_already_executed",
	Title:  "Refund Already Executed",
	Status: http.StatusBadRequest,
	Detail: "The refund has already been executed.",
}

// StellarRefundValidator checks if it is possible to
// refund a deposit to depositor's Stellar account.
type StellarRefundValidator struct {
	Session                db.SessionInterface
	WithdrawalWindow       time.Duration
	Observer               ethereum.Observer
	EthereumFinalityBuffer uint64
}

// StellarRefundDetails includes metadata about the
// validation result.
type StellarRefundDetails struct {
	// LedgerSequence is the sequence number of the Stellar ledger
	// for which the validation result is accurate.
	LedgerSequence uint32
}

func (s StellarRefundValidator) CanRefund(ctx context.Context, deposit store.StellarDeposit) (StellarRefundDetails, error) {
	dbStore := store.DB{Session: s.Session.Clone()}
	err := dbStore.Session.BeginTx(&sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		return StellarRefundDetails{}, errors.Wrap(err, "error starting repeatable read transaction")
	}
	defer func() {
		// explicitly ignore return value to make the linter happy
		_ = dbStore.Session.Rollback()
	}()

	lastLedgerSequence, err := dbStore.GetLastLedgerSequence(context.Background())
	if err != nil {
		return StellarRefundDetails{}, errors.Wrap(err, "error getting last ledger sequence")
	}

	// Check if refund tx was seen without signature request
	exists, err := dbStore.HistoryStellarTransactionExists(ctx, deposit.ID)
	if err != nil {
		return StellarRefundDetails{}, errors.Wrap(err, "error getting history stellar transaction by memo hash")
	}
	if exists {
		return StellarRefundDetails{}, RefundAlreadyExecuted
	}

	withdrawalDeadline := time.Unix(deposit.LedgerTime, 0).Add(s.WithdrawalWindow)

	// rollback to release used DB connection because further checks
	// do not involve DB
	_ = dbStore.Session.Rollback()

	// Checks on Ethereum side:
	// - Ensure that there was no withdrawal to Ethereum account
	// - The response from the client is after the withdrawal deadline
	depositID := common.HexToHash(deposit.ID)
	requestStatus, err := s.Observer.GetRequestStatus(ctx, depositID)
	if err != nil {
		return StellarRefundDetails{}, errors.Wrap(err, "error getting request status from ethereum observer")
	}

	if requestStatus.BlockNumber <= s.EthereumFinalityBuffer {
		return StellarRefundDetails{}, EthereumNodeBehind
	}

	block, err := s.Observer.GetBlockByNumber(ctx, requestStatus.BlockNumber-s.EthereumFinalityBuffer)
	if err != nil {
		return StellarRefundDetails{}, errors.Wrap(err, "error getting block from ethereum observer")
	}

	if !block.Time.After(withdrawalDeadline) {
		return StellarRefundDetails{}, WithdrawalWindowStillActive
	}

	if requestStatus.Fulfilled {
		return StellarRefundDetails{}, WithdrawalAlreadyExecuted
	}

	return StellarRefundDetails{
		LedgerSequence: lastLedgerSequence,
	}, nil
}

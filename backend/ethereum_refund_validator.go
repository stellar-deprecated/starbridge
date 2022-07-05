package backend

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/stellar/go/support/db"

	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/store"
)

var WithdrawalWindowStillActive = problem.P{
	Type:   "withdrawal_window_still_active",
	Title:  "Withdrawal Window Still Active",
	Status: http.StatusBadRequest,
	Detail: "The withdrawal window is still active." +
		" Wait until the withdrawal window has closed before attempting a refund.",
}

// EthereumRefundValidator checks if it is possible to
// refund a deposit to the ethereum bridge smart contract.
type EthereumRefundValidator struct {
	Session          db.SessionInterface
	WithdrawalWindow time.Duration
}

func (s EthereumRefundValidator) CanRefund(ctx context.Context, deposit store.EthereumDeposit) error {
	dbStore := store.DB{Session: s.Session.Clone()}
	err := dbStore.Session.BeginTx(&sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		return errors.Wrap(err, "error starting repeatable read transaction")
	}
	defer func() {
		// explicitly ignore return value to make the linter happy
		_ = dbStore.Session.Rollback()
	}()

	lastLedgerCloseTime, err := dbStore.GetLastLedgerCloseTime(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting last ledger close time")
	}
	withdrawalDeadline := time.Unix(deposit.BlockTime, 0).Add(s.WithdrawalWindow)
	log.Info(lastLedgerCloseTime)
	log.Info(deposit.BlockTime)
	log.Info(withdrawalDeadline)
	if !lastLedgerCloseTime.After(withdrawalDeadline) {
		return WithdrawalWindowStillActive
	}

	// Check if withdrawal tx was seen without signature request
	exists, err := dbStore.HistoryStellarTransactionExists(ctx, deposit.ID)
	if err != nil {
		return errors.Wrap(err, "error getting history stellar transaction by memo hash")
	}
	if exists {
		return WithdrawalAlreadyExecuted
	}

	return nil
}

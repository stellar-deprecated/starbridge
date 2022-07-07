package backend

import (
	"context"
	"database/sql"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stellar/go/strkey"
	"github.com/stellar/go/support/db"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/store"
)

var (
	WithdrawalWindowExpired = problem.P{
		Type:   "withdrawal_window_expired",
		Title:  "Withdrawal Window Expired",
		Status: http.StatusBadRequest,
		Detail: "The withdrawal window has expired. Only refunds are allowed at this point.",
	}
	WithdrawalAlreadyExecuted = problem.P{
		Type:   "withdrawal_already_executed",
		Title:  "Withdrawal Already Executed",
		Status: http.StatusBadRequest,
		Detail: "The withdrawal has already been executed.",
	}
	WithdrawalAssetInvalid = problem.P{
		Type:   "withdrawal_asset_invalid",
		Title:  "Withdrawal Asset Invalid",
		Status: http.StatusBadRequest,
		Detail: "Withdrawing the requested asset is not supported by the bridge." +
			"Refund the deposit once the withdrawal period has expired.",
	}
	WithdrawalAmountInvalid = problem.P{
		Type:   "withdrawal_amount_invalid",
		Title:  "Withdrawal Amount Invalid",
		Status: http.StatusBadRequest,
		Detail: "Withdrawing the requested amount is not supported by the bridge." +
			"Refund the deposit once the withdrawal period has expired.",
	}
	InvalidStellarRecipient = problem.P{
		Type:   "invalid_stellar_recipient",
		Title:  "Invalid Stellar Recipient",
		Status: http.StatusBadRequest,
		Detail: "The recipient of the deposit is not a valid Stellar address.",
	}

	ethereumTokenAddress = common.Address{}
)

// StellarWithdrawalValidator checks if it is possible to
// withdraw a deposit to the ethereum bridge smart contract on
// Stellar.
type StellarWithdrawalValidator struct {
	Session          db.SessionInterface
	WithdrawalWindow time.Duration
}

// StellarWithdrawalDetails includes metadata about the
// validation result.
type StellarWithdrawalDetails struct {
	// Deadline is the deadline for executing the withdrawal
	// transaction on Stellar.
	Deadline time.Time
	// Recipient is the Stellar account which should receive the
	// withdrawal.
	Recipient string
	// LedgerSequence is the sequence number of the Stellar ledger
	// for which the validation result is accurate.
	LedgerSequence uint32
}

func (s StellarWithdrawalValidator) CanWithdraw(ctx context.Context, deposit store.EthereumDeposit) (StellarWithdrawalDetails, error) {
	// TODO: add support for erc20 transfers
	if !common.IsHexAddress(deposit.Token) ||
		common.HexToAddress(deposit.Token) != ethereumTokenAddress {
		return StellarWithdrawalDetails{}, WithdrawalAssetInvalid
	}

	// TODO: implement amount validation which is specific to the type of token
	amount := &big.Int{}
	_, ok := amount.SetString(deposit.Amount, 10)
	if !ok || !amount.IsInt64() || amount.Cmp(big.NewInt(0)) <= 0 {
		return StellarWithdrawalDetails{}, WithdrawalAmountInvalid
	}

	destination, ok := new(big.Int).SetString(deposit.Destination, 10)
	if !ok {
		return StellarWithdrawalDetails{}, InvalidStellarRecipient
	}
	destinationAccountID, err := strkey.Encode(
		strkey.VersionByteAccountID,
		destination.Bytes(),
	)
	if err != nil {
		return StellarWithdrawalDetails{}, InvalidStellarRecipient
	}

	dbStore := store.DB{Session: s.Session.Clone()}
	err = dbStore.Session.BeginTx(&sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		return StellarWithdrawalDetails{}, errors.Wrap(err, "error starting repeatable read transaction")
	}
	defer func() {
		// explicitly ignore return value to make the linter happy
		_ = dbStore.Session.Rollback()
	}()

	lastLedgerSequence, err := dbStore.GetLastLedgerSequence(context.Background())
	if err != nil {
		return StellarWithdrawalDetails{}, errors.Wrap(err, "error getting last ledger sequence")
	}

	lastLedgerCloseTime, err := dbStore.GetLastLedgerCloseTime(ctx)
	if err != nil {
		return StellarWithdrawalDetails{}, errors.Wrap(err, "error getting last ledger close time")
	}
	withdrawalDeadline := time.Unix(deposit.BlockTime, 0).Add(s.WithdrawalWindow)
	if lastLedgerCloseTime.After(withdrawalDeadline) {
		return StellarWithdrawalDetails{}, WithdrawalWindowExpired
	}

	// Check if withdrawal tx was seen without signature request
	exists, err := dbStore.HistoryStellarTransactionExists(ctx, deposit.ID)
	if err != nil {
		return StellarWithdrawalDetails{}, errors.Wrap(err, "error getting history stellar transaction by memo hash")
	}
	if exists {
		return StellarWithdrawalDetails{}, WithdrawalAlreadyExecuted
	}

	return StellarWithdrawalDetails{
		Deadline:       withdrawalDeadline,
		Recipient:      destinationAccountID,
		LedgerSequence: lastLedgerSequence,
	}, nil
}

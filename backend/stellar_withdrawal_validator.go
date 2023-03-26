package backend

import (
	"context"
	"database/sql"
	"math/big"
	"net/http"
	"time"

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
	InvalidStellarRecipient = problem.P{
		Type:   "invalid_stellar_recipient",
		Title:  "Invalid Stellar Recipient",
		Status: http.StatusBadRequest,
		Detail: "The recipient of the deposit is not a valid Stellar address.",
	}
)

// StellarWithdrawalValidator checks if it is possible to
// withdraw a deposit to the ethereum bridge smart contract on
// Stellar.
type StellarWithdrawalValidator struct {
	Session          db.SessionInterface
	WithdrawalWindow time.Duration
	Converter        AssetConverter
	CcdToken         string
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
	// Asset is the Stellar asset which will be transferred to the
	// recipient.
	Asset string
	// Amount is the amount which will be transferred to the recipient.
	Amount int64
}

func (s StellarWithdrawalValidator) CanWithdraw(ctx context.Context, deposit store.EthereumDeposit) (StellarWithdrawalDetails, error) {
	stellarAsset, stellarAmount, err := s.Converter.ToStellar(deposit.Token, deposit.Amount, true, false, false)
	if err != nil {
		return StellarWithdrawalDetails{}, err
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

	lastLedgerSequence, err := dbStore.GetLastLedgerSequence(ctx)
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
		Asset:          stellarAsset,
		Amount:         stellarAmount,
	}, nil
}

func (s StellarWithdrawalValidator) CanWithdrawOkx(ctx context.Context, deposit store.OkxDeposit) (StellarWithdrawalDetails, error) {
	stellarAsset, stellarAmount, err := s.Converter.ToStellar(deposit.Token, deposit.Amount, false, false, true)
	if err != nil {
		return StellarWithdrawalDetails{}, err
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

	lastLedgerSequence, err := dbStore.GetLastLedgerSequence(ctx)
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
		Asset:          stellarAsset,
		Amount:         stellarAmount,
	}, nil
}

func (s StellarWithdrawalValidator) CanWithdrawConcordium(ctx context.Context, deposit store.ConcordiumDeposit) (StellarWithdrawalDetails, error) {
	stellarAsset, stellarAmount, err := s.Converter.ToStellar(s.CcdToken, deposit.Amount, false, true, false)
	if err != nil {
		return StellarWithdrawalDetails{}, err
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

	lastLedgerSequence, err := dbStore.GetLastLedgerSequence(ctx)
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
		Asset:          stellarAsset,
		Amount:         stellarAmount,
	}, nil
}

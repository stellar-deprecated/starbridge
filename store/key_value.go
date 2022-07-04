package store

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/stellar/go/support/errors"
)

const (
	lastLedgerSequenceKey  = "last_ledger_sequence"
	lastLedgerCloseTimeKey = "last_ledger_close_time"
)

func (m *DB) GetLastLedgerSequence(ctx context.Context) (uint32, error) {
	lastLedgerSequence, err := m.getValueFromStore(ctx, lastLedgerSequenceKey)
	if err != nil {
		return 0, err
	}

	if lastLedgerSequence == "" {
		return 0, nil
	} else {
		ledgerSequence, err := strconv.ParseUint(lastLedgerSequence, 10, 32)
		if err != nil {
			return 0, errors.Wrap(err, "Error converting lastLedgerSequence value")
		}

		return uint32(ledgerSequence), nil
	}
}

func (m *DB) UpdateLastLedgerSequence(ctx context.Context, ledgerSequence uint32) error {
	return m.updateValueInStore(
		ctx,
		lastLedgerSequenceKey,
		strconv.FormatUint(uint64(ledgerSequence), 10),
	)
}

func (m *DB) GetLastLedgerCloseTime(ctx context.Context) (time.Time, error) {
	lastLedgerCloseTime, err := m.getValueFromStore(ctx, lastLedgerCloseTimeKey)
	if err != nil {
		return time.Now(), err
	}

	if lastLedgerCloseTime == "" {
		return time.Now(), errors.Errorf("no value for key: %s", lastLedgerCloseTimeKey)
	} else {
		ledgerCloseTime, err := strconv.ParseInt(lastLedgerCloseTime, 10, 64)
		if err != nil {
			return time.Now(), errors.Wrap(err, "Error converting lastLedgerCloseTime value")
		}

		return time.Unix(ledgerCloseTime, 0), nil
	}
}

func (m *DB) UpdateLastLedgerCloseTime(ctx context.Context, closeTime time.Time) error {
	return m.updateValueInStore(
		ctx,
		lastLedgerCloseTimeKey,
		strconv.FormatInt(closeTime.Unix(), 10),
	)
}

// getValueFromStore returns a value for a given key from KV store. If value
// is not present in the key value store "" will be returned.
func (m *DB) getValueFromStore(ctx context.Context, key string) (string, error) {
	query := sq.Select("key_value_store.value").
		From("key_value_store").
		Where("key_value_store.key = ?", key)

	var value string
	if err := m.Session.Get(ctx, &value, query); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return "", nil
		}
		return "", errors.Wrap(err, "could not get value")
	}

	return value, nil
}

// updateValueInStore updates a value for a given key in KV store
func (m *DB) updateValueInStore(ctx context.Context, key, value string) error {
	query := sq.Insert("key_value_store").
		Columns("key", "value").
		Values(key, value).
		Suffix("ON CONFLICT (key) DO UPDATE SET value=EXCLUDED.value")

	_, err := m.Session.Exec(ctx, query)
	return err
}

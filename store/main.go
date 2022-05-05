package store

import (
	"github.com/stellar/go/support/db"
)

type DB struct {
	Session *db.Session
}

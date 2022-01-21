package main

import (
	"fmt"

	"github.com/stellar/go/txnbuild"
)

// TODO: Clear out transactions that are not needed anymore in the in-memory store
// TODO: Support fee bump transactions.

type Store struct {
	transactions map[[32]byte]*txnbuild.Transaction
}

func NewStore() *Store {
	return &Store{
		transactions: map[[32]byte]*txnbuild.Transaction{},
	}
}

func (s *Store) StoreAndUpdate(txHash [32]byte, tx *txnbuild.Transaction) (*txnbuild.Transaction, error) {
	storedTx := s.transactions[txHash]
	if storedTx != nil {
		sigsSeen := map[string]bool{}
		for _, s := range tx.Signatures() {
			b, err := s.MarshalBinary()
			if err != nil {
				return nil, fmt.Errorf("unexpected error marshaling sig %x: %w", txHash, err)
			}
			sigsSeen[string(b)] = true
		}
		for _, s := range storedTx.Signatures() {
			b, err := s.MarshalBinary()
			if err != nil {
				return nil, fmt.Errorf("unexpected error marshaling sig %x: %w", txHash, err)
			}
			if sigsSeen[string(b)] {
				continue
			}
			tx, err = tx.AddSignatureDecorated(s)
			if err != nil {
				return nil, fmt.Errorf("adding signature to tx %x: %w", txHash, err)
			}
		}
	}
	s.transactions[txHash] = tx
	return tx, nil
}

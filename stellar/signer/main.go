package signer

import (
	"encoding/hex"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/xdr"
)

type Signer struct {
	NetworkPassphrase string
	SecretKey         string

	kp *keypair.Full
}

// Sign signs an envelope.
func (s *Signer) Sign(envelope xdr.TransactionEnvelope) (string, xdr.DecoratedSignature, error) {
	hash, err := network.HashTransactionInEnvelope(envelope, s.NetworkPassphrase)
	if err != nil {
		return "", xdr.DecoratedSignature{}, errors.Wrap(err, "failed to hash transaction")
	}

	if s.kp == nil {
		s.kp = keypair.MustParseFull(s.SecretKey)
	}

	sig, err := s.kp.SignDecorated(hash[:])
	if err != nil {
		return "", xdr.DecoratedSignature{}, errors.Wrap(err, "failed to sign transaction")
	}

	return hex.EncodeToString(hash[:]), sig, nil
}

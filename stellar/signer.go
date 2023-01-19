package stellar

import (
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/xdr"
)

type Signer struct {
	NetworkPassphrase string
	Signer            *keypair.Full
}

// Sign signs an envelope.
func (s *Signer) Sign(envelope xdr.TransactionEnvelope) (xdr.DecoratedSignature, error) {
	hash, err := network.HashTransactionInEnvelope(envelope, s.NetworkPassphrase)
	if err != nil {
		return xdr.DecoratedSignature{}, errors.Wrap(err, "failed to hash transaction")
	}

	sig, err := s.Signer.SignDecorated(hash[:])
	if err != nil {
		return xdr.DecoratedSignature{}, errors.Wrap(err, "failed to sign transaction")
	}

	return sig, nil
}

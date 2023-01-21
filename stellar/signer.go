package stellar

import (
	"github.com/pkg/errors"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

type Signer struct {
	BridgeAccount     string
	BridgeContractID  [32]byte
	NetworkPassphrase string
	Signer            *keypair.Full
}

// NewWithdrawalTransaction builds and signs a transaction. It does not check if expirationTimestamp is valid.
func (s *Signer) NewWithdrawalTransaction(assetContractID [32]byte, isWrappedAsset bool, txSource, destination, amount string, sequence, expirationTimestamp int64, id [32]byte) (xdr.TransactionEnvelope, error) {
	if txSource == s.BridgeAccount {
		return xdr.TransactionEnvelope{}, errors.New("bridge account cannot be used as a transaction source")
	}

	sourceAccount := txnbuild.SimpleAccount{
		AccountID: txSource,
		Sequence:  sequence,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &sourceAccount,
			Operations: []txnbuild.Operation{
				// TODO construct invoke host function op

			},
			BaseFee: txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewTimebounds(0, expirationTimestamp),
			},
		},
	)
	if err != nil {
		return xdr.TransactionEnvelope{}, errors.Wrap(err, "error building transaction")
	}

	tx, err = tx.Sign(s.NetworkPassphrase, s.Signer)
	if err != nil {
		return xdr.TransactionEnvelope{}, errors.Wrap(err, "failed to sign transaction")
	}

	return tx.ToXDR(), nil
}

package txbuilder

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

type Builder struct {
	BridgeAccount string
}

// BuildTransaction builds a transaction. It does not check if expirationTimestamp is valid.
func (b *Builder) BuildTransaction(asset, txSource, destination, amount string, sequence, expirationTimestamp int64, memoHash []byte) (xdr.TransactionEnvelope, error) {
	if txSource == b.BridgeAccount {
		return xdr.TransactionEnvelope{}, errors.New("bridge account cannot be used as a transaction source")
	}

	sourceAccount := txnbuild.SimpleAccount{
		AccountID: txSource,
		Sequence:  sequence,
	}

	var memoHashArray txnbuild.MemoHash
	copy(memoHashArray[:], memoHash)

	var txAsset txnbuild.Asset
	if asset == "native" {
		txAsset = txnbuild.NativeAsset{}
	} else {
		parts := strings.Split(asset, ":")
		txAsset = txnbuild.CreditAsset{
			Code:   parts[0],
			Issuer: parts[1],
		}
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &sourceAccount,
			Memo:          memoHashArray,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					SourceAccount: b.BridgeAccount,
					Amount:        amount,
					Destination:   destination,
					Asset:         txAsset,
				},
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

	return tx.ToXDR(), nil
}

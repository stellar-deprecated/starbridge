package txbuilder

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

type Builder struct {
	HorizonURL    string
	BridgeAccount string
}

// BuildTransaction builds a transaction. It does not check if expirationTimestamp is valid.
func (b *Builder) BuildTransaction(txSource, destination, amount string, expirationTimestamp int64) (xdr.TransactionEnvelope, error) {
	// TODO remove seqnum fetch from here. it should be provided by the user
	client := &horizonclient.Client{
		HorizonURL: b.HorizonURL,
	}

	if txSource == b.BridgeAccount {
		return xdr.TransactionEnvelope{}, errors.New("bridge account cannot be used as a transaction source")
	}

	sourceAccount, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: txSource})
	if err != nil {
		return xdr.TransactionEnvelope{}, errors.Wrap(err, "error getting account details")
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			IncrementSequenceNum: true,

			SourceAccount: &sourceAccount,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					SourceAccount: b.BridgeAccount,
					Amount:        amount,
					Destination:   destination,
					Asset: txnbuild.CreditAsset{
						Code:   "ETH",
						Issuer: b.BridgeAccount,
					},
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

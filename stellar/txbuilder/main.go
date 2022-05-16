package txbuilder

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

type Builder struct {
	BridgeAccount string
}

func (b *Builder) BuildTransaction(txSource, destination, amount string) (xdr.TransactionEnvelope, error) {
	// TODO remove seqnum fetch from here. it should be provided by the user
	client := &horizonclient.Client{
		HorizonURL: "http://localhost:8000",
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
			// TODO: one minute for faster debugging, change do 5m/10m
			Preconditions: txnbuild.Preconditions{TimeBounds: txnbuild.NewTimeout(60)},
		},
	)
	if err != nil {
		return xdr.TransactionEnvelope{}, errors.Wrap(err, "error building transaction")
	}

	return tx.ToXDR(), nil
}

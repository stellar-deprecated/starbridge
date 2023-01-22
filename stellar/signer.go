package stellar

import (
	"github.com/pkg/errors"

	"github.com/stellar/go/clients/stellarcore"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"

	soroban_bridge "github.com/stellar/starbridge/soroban-bridge"
)

type Signer struct {
	BridgeAccount     string
	BridgeContractID  [32]byte
	NetworkPassphrase string
	Signer            *keypair.Full
	CoreClient        *stellarcore.Client
}

// NewWithdrawalTransaction builds and signs a transaction. It does not check if expirationTimestamp is valid.
func (s *Signer) NewWithdrawalTransaction(assetContractID [32]byte, isWrappedAsset bool, txSource, destination, assetAmount string, sequence, expirationTimestamp int64, id [32]byte) (xdr.TransactionEnvelope, error) {
	if txSource == s.BridgeAccount {
		return xdr.TransactionEnvelope{}, errors.New("bridge account cannot be used as a transaction source")
	}

	sourceAccount := txnbuild.SimpleAccount{
		AccountID: txSource,
		Sequence:  sequence,
	}

	amountParam, err := soroban_bridge.AmountContractParam(assetAmount)
	if err != nil {
		return xdr.TransactionEnvelope{}, err
	}

	invokeWithdraw := &txnbuild.InvokeHostFunction{
		Function: xdr.HostFunction{
			Type: xdr.HostFunctionTypeHostFunctionTypeInvokeContract,
			InvokeArgs: &xdr.ScVec{
				soroban_bridge.BytesContractParam(s.BridgeContractID[:]),
				soroban_bridge.FunctionNameParam("withdraw"),
				soroban_bridge.BytesContractParam(assetContractID[:]),
				soroban_bridge.BoolContractParam(isWrappedAsset),
				soroban_bridge.AccountIDEnumParam(destination),
				soroban_bridge.BytesContractParam(id[:]),
				amountParam,
			},
		},
		SourceAccount: s.BridgeAccount,
	}
	if _, err := soroban_bridge.Preflight(s.CoreClient, invokeWithdraw); err != nil {
		return xdr.TransactionEnvelope{}, err
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &sourceAccount,
			Operations: []txnbuild.Operation{
				invokeWithdraw,
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

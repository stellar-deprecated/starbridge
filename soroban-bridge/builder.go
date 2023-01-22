package soroban_bridge

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/pkg/errors"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/clients/stellarcore"
	stellarcoreproto "github.com/stellar/go/protocols/stellarcore"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func StellarAssetContractID(passPhrase string, asset xdr.Asset) (xdr.Hash, error) {
	networkId := xdr.Hash(sha256.Sum256([]byte(passPhrase)))
	preImage := xdr.HashIdPreimage{
		Type: xdr.EnvelopeTypeEnvelopeTypeContractIdFromAsset,
		FromAsset: &xdr.HashIdPreimageFromAsset{
			NetworkId: networkId,
			Asset:     asset,
		},
	}
	xdrPreImageBytes, err := preImage.MarshalBinary()
	if err != nil {
		return xdr.Hash{}, err
	}
	return sha256.Sum256(xdrPreImageBytes), nil
}

func AmountContractParam(assetAmount string) (xdr.ScVal, error) {
	parsedAmount, err := amount.ParseInt64(assetAmount)
	if err != nil || parsedAmount <= 0 {
		return xdr.ScVal{}, errors.Wrap(err, "asset amount is invalid")
	}
	obj := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoI128,
		I128: &xdr.Int128Parts{
			Lo: xdr.Uint64(parsedAmount),
			Hi: 0,
		},
	}
	return xdr.ScVal{
		Type: xdr.ScValTypeScvObject,
		Obj:  &obj,
	}, nil
}

func BoolContractParam(b bool) xdr.ScVal {
	ic := xdr.ScStaticScsFalse
	if b {
		ic = xdr.ScStaticScsTrue
	}
	return xdr.ScVal{
		Type: xdr.ScValTypeScvStatic,
		Ic:   &ic,
	}
}

func AccountIDEnumParam(accountID string) xdr.ScVal {
	accountObj := &xdr.ScObject{
		Type:      xdr.ScObjectTypeScoAccountId,
		AccountId: xdr.MustAddressPtr(accountID),
	}
	accountSym := xdr.ScSymbol("Account")
	accountEnum := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoVec,
		Vec: &xdr.ScVec{
			xdr.ScVal{
				Type: xdr.ScValTypeScvSymbol,
				Sym:  &accountSym,
			},
			xdr.ScVal{
				Type: xdr.ScValTypeScvObject,
				Obj:  &accountObj,
			},
		},
	}
	return xdr.ScVal{
		Type: xdr.ScValTypeScvObject,
		Obj:  &accountEnum,
	}
}

func ContractIDEnumParam(contractID xdr.Hash) xdr.ScVal {
	contractIDBytes := contractID[:]
	contractIDObj := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoBytes,
		Bin:  &contractIDBytes,
	}
	contractSym := xdr.ScSymbol("Contract")
	contractEnum := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoVec,
		Vec: &xdr.ScVec{
			xdr.ScVal{
				Type: xdr.ScValTypeScvSymbol,
				Sym:  &contractSym,
			},
			xdr.ScVal{
				Type: xdr.ScValTypeScvObject,
				Obj:  &contractIDObj,
			},
		},
	}
	return xdr.ScVal{
		Type: xdr.ScValTypeScvObject,
		Obj:  &contractEnum,
	}
}

func FunctionNameParam(name string) xdr.ScVal {
	contractFnParameterSym := xdr.ScSymbol(name)
	return xdr.ScVal{
		Type: xdr.ScValTypeScvSymbol,
		Sym:  &contractFnParameterSym,
	}
}

func BytesContractParam(bytes []byte) xdr.ScVal {
	parameterObj := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoBytes,
		Bin:  &bytes,
	}
	return xdr.ScVal{
		Type: xdr.ScValTypeScvObject,
		Obj:  &parameterObj,
	}
}

func Preflight(coreClient *stellarcore.Client, invokeHostFn *txnbuild.InvokeHostFunction) (stellarcoreproto.PreflightResponse, error) {
	opXDR, err := invokeHostFn.BuildXDR()
	if err != nil {
		return stellarcoreproto.PreflightResponse{}, err
	}

	invokeHostFunctionOp := opXDR.Body.MustInvokeHostFunctionOp()

	response, err := coreClient.Preflight(
		context.Background(),
		invokeHostFn.SourceAccount,
		invokeHostFunctionOp,
	)
	if err != nil {
		return response, err
	}

	if response.Status != stellarcoreproto.PreflightStatusOk {
		return response, fmt.Errorf("status is not ok: %v", response.Detail)
	}

	if response.Status != stellarcoreproto.PreflightStatusOk {
		return response, errors.New("preflight did not succeed")
	}
	if err := xdr.SafeUnmarshalBase64(response.Footprint, &invokeHostFn.Footprint); err != nil {
		return response, errors.Wrap(err, "could not decode footprint")
	}
	return response, nil
}

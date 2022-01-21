package main

import (
	"errors"
	"fmt"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

var ErrNotAuthorized = errors.New("not authorized")

func AuthorizedTransaction(horizonClient horizonclient.ClientInterface, txHash [32]byte, tx *txnbuild.Transaction) (*txnbuild.Transaction, error) {
	// Get accounts and thresholds required.
	txThresholds := TransactionThresholds(tx)

	// Get the thresholds and signers of the accounts required by the tx.
	accountThresholds := map[string]byte{}
	type accountSigner struct {
		Account string
		Signer  horizon.Signer
	}
	accountSignersByHint := map[[4]byte][]accountSigner{}
	for account, threshold := range txThresholds {
		horizonAccount, err := horizonClient.AccountDetail(horizonclient.AccountRequest{AccountID: account})
		if err == nil {
			accountThreshold := byte(0)
			switch threshold {
			case LowThreshold:
				accountThreshold = horizonAccount.Thresholds.LowThreshold
			case MedThreshold:
				accountThreshold = horizonAccount.Thresholds.MedThreshold
			case HighThreshold:
				accountThreshold = horizonAccount.Thresholds.HighThreshold
			}
			if accountThreshold == 0 {
				accountThreshold = 1
			}
			accountThresholds[account] = accountThreshold
			for _, signer := range horizonAccount.Signers {
				hint, err := GetHint(signer.Key)
				if err != nil {
					return nil, err
				}
				accountSignersByHint[hint] = append(accountSignersByHint[hint], accountSigner{
					Account: account,
					Signer:  signer,
				})
			}
		} else if horizonclient.IsNotFoundError(err) {
			accountThresholds[account] = 1
			hint, err := GetHint(account)
			if err != nil {
				return nil, err
			}
			signer := horizon.Signer{Type: "ed25519_public_key", Key: account, Weight: 1}
			accountSignersByHint[hint] = append(accountSignersByHint[hint], accountSigner{
				Account: account,
				Signer:  signer,
			})
		} else {
			return nil, err
		}
	}

	// Collect the weights that signatures provide, and a list of signatures
	// need to authorize.
	sigWeights := map[string]int32{}
	authorizingSigs := []xdr.DecoratedSignature{}
	for _, sig := range tx.Signatures() {
		accountSigners := accountSignersByHint[sig.Hint]
		addSig := false
		for _, accountSigner := range accountSigners {
			account := accountSigner.Account
			if sigWeights[account] >= int32(accountThresholds[account]) {
				continue
			}

			signer := accountSigner.Signer
			switch signer.Type {
			case "ed25519_public_key":
				kp, err := keypair.ParseAddress(signer.Key)
				if err != nil {
					return nil, err
				}
				err = kp.Verify(txHash[:], sig.Signature)
				if err == keypair.ErrInvalidSignature {
					continue
				} else if err != nil {
					return nil, err
				}
			default:
				// TODO: Support all signer types.
				return nil, fmt.Errorf("unsupported signer: %s", signer.Type)
			}

			sigWeights[account] += signer.Weight
			addSig = true
		}
		if addSig {
			authorizingSigs = append(authorizingSigs, sig)
		}
	}

	// Determine if authorized.
	for account, threshold := range accountThresholds {
		if sigWeights[account] < int32(threshold) {
			return nil, ErrNotAuthorized
		}
	}

	// Build transaction with authorizing signatures.
	tx, err := tx.ClearSignatures()
	if err != nil {
		return nil, err
	}
	tx, err = tx.AddSignatureDecorated(authorizingSigs...)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

type Threshold int

const (
	LowThreshold Threshold = iota
	MedThreshold
	HighThreshold
)

func (t Threshold) And(t2 Threshold) Threshold {
	if t2 > t {
		return t2
	}
	return t
}

var operationThresholds = map[xdr.OperationType]Threshold{
	xdr.OperationTypeBumpSequence:      LowThreshold,
	xdr.OperationTypeAllowTrust:        LowThreshold,
	xdr.OperationTypeSetTrustLineFlags: LowThreshold,

	xdr.OperationTypeSetOptions:   HighThreshold,
	xdr.OperationTypeAccountMerge: HighThreshold,

	xdr.OperationTypeCreateAccount:                 MedThreshold,
	xdr.OperationTypePayment:                       MedThreshold,
	xdr.OperationTypePathPaymentStrictReceive:      MedThreshold,
	xdr.OperationTypeManageSellOffer:               MedThreshold,
	xdr.OperationTypeCreatePassiveSellOffer:        MedThreshold,
	xdr.OperationTypeChangeTrust:                   MedThreshold,
	xdr.OperationTypeInflation:                     MedThreshold,
	xdr.OperationTypeManageData:                    MedThreshold,
	xdr.OperationTypeManageBuyOffer:                MedThreshold,
	xdr.OperationTypePathPaymentStrictSend:         MedThreshold,
	xdr.OperationTypeBeginSponsoringFutureReserves: MedThreshold,
	xdr.OperationTypeEndSponsoringFutureReserves:   MedThreshold,
	xdr.OperationTypeCreateClaimableBalance:        MedThreshold,
	xdr.OperationTypeClaimClaimableBalance:         MedThreshold,
	xdr.OperationTypeRevokeSponsorship:             MedThreshold,
	xdr.OperationTypeClawback:                      MedThreshold,
	xdr.OperationTypeClawbackClaimableBalance:      MedThreshold,
	xdr.OperationTypeLiquidityPoolDeposit:          MedThreshold,
	xdr.OperationTypeLiquidityPoolWithdraw:         MedThreshold,
}

func OperationThreshold(op xdr.OperationType) Threshold {
	t, ok := operationThresholds[op]
	if !ok {
		panic(fmt.Errorf("unexpected operation type %d is unknown", op))
	}
	return t
}

func FeeBumpTransactionThresholds(tx *txnbuild.FeeBumpTransaction) map[string]Threshold {
	return map[string]Threshold{
		tx.FeeAccount(): LowThreshold,
	}
}

func TransactionThresholds(tx *txnbuild.Transaction) map[string]Threshold {
	txEnv := tx.ToXDR()

	thresholds := map[string]Threshold{}
	add := func(account string, t Threshold) {
		thresholds[account] = thresholds[account].And(t)
	}

	txSource := txEnv.SourceAccount()
	add(txSource.ToAccountId().Address(), LowThreshold)

	for _, op := range txEnv.Operations() {
		opSource := op.SourceAccount
		if opSource == nil {
			opSource = &txSource
		}
		opThreshold := OperationThreshold(op.Body.Type)
		add(opSource.ToAccountId().Address(), opThreshold)
	}

	return thresholds
}

func GetHint(key string) (hint [4]byte, err error) {
	_, raw, err := strkey.DecodeAny(key)
	if err != nil {
		return
	}
	if len(raw) < 4 {
		err = fmt.Errorf("key too short to calculate hint: %s", key)
		return
	}
	copy(hint[:], raw[len(raw)-4:])
	return
}

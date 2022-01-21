package integrations

import (
	"fmt"
	"strings"

	"github.com/stellar/go/txnbuild"
	"github.com/stellar/starbridge/model"
)

// TODO need to set the contract account, source account
var sourceAccount = "GDXLLPH23EHIFOQLO46X2WQQPBBRJ6YPV7JOXEWS7V3AU74Z4EY7PGCS" // var sourceSecretKey = "SAJOMEU6AAHWIUSF43Z7BGFEEB4VUCZVTG56U4DU6RR3UGAZRFSHEYEQ"

var baseFee int64 = 100

func getStellarAsset(assetInfo *model.AssetInfo) txnbuild.Asset {
	if assetInfo.ContractAddress == "native" {
		return txnbuild.NativeAsset{}
	}
	return txnbuild.CreditAsset{
		Code:   assetInfo.Code,
		Issuer: assetInfo.ContractAddress,
	}
}

func Transaction2Stellar(tx *model.Transaction) (*txnbuild.Transaction, error) {
	if tx.Chain != model.ChainStellar {
		return nil, fmt.Errorf("cannot convert transaction from a different chain ('%s') to Stellar, need to convert the transaction to the Stellar chain first", tx.Chain.Name)
	}

	ops := []txnbuild.Operation{}
	ops = append(ops, &txnbuild.CreateClaimableBalance{
		Destinations: []txnbuild.Claimant{
			txnbuild.NewClaimant(tx.To, &txnbuild.UnconditionalPredicate),
		},
		Asset:  getStellarAsset(tx.AssetInfo),
		Amount: fmt.Sprintf("%d", tx.Amount),
		// SourceAccount: nil,
	})

	return txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: sourceAccount,
				Sequence:  int64(tx.SeqNum),
			},
			BaseFee:              baseFee,
			IncrementSequenceNum: true,
			Operations:           ops,
			// TODO need to set timebounds correctly
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
}

// Stellar2String is a String converter
func Stellar2String(tx *txnbuild.Transaction) string {
	memoString := ""
	if tx.Memo() != nil {
		memoXdr, _ := tx.Memo().ToXDR()
		memoString = memoXdr.GoString()
	}

	sb := strings.Builder{}
	sb.WriteString("StellarTx[")
	sb.WriteString(fmt.Sprintf("SourceAccount=%s", tx.SourceAccount().AccountID))
	sb.WriteString(fmt.Sprintf(", SeqNum=%d", tx.SequenceNumber()))
	sb.WriteString(fmt.Sprintf(", BaseFee=%d", tx.BaseFee()))
	sb.WriteString(fmt.Sprintf(", MaxFee=%d", tx.MaxFee()))
	sb.WriteString(fmt.Sprintf(", TimeBounds.MinTime=%d", tx.Timebounds().MinTime))
	sb.WriteString(fmt.Sprintf(", TimeBounds.MaxTime=%d", tx.Timebounds().MaxTime))
	sb.WriteString(fmt.Sprintf(", Memo=%s", memoString))
	sb.WriteString(fmt.Sprintf(", Operations=%s", stellarOps2String(tx.Operations())))
	sb.WriteString("]")
	return sb.String()
}

func stellarOps2String(ops []txnbuild.Operation) string {
	sb := strings.Builder{}
	sb.WriteString("Array[")

	for i := 0; i < len(ops); i++ {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(stellarOp2String(ops[i]))
	}

	sb.WriteString("]")
	return sb.String()
}

func stellarOp2String(op txnbuild.Operation) string {
	sb := strings.Builder{}
	switch o := op.(type) {
	case *txnbuild.CreateClaimableBalance:
		sb.WriteString("Operation[")
		sb.WriteString(fmt.Sprintf("SourceAccount=%s", op.GetSourceAccount()))
		sb.WriteString(fmt.Sprintf(", Type=%s", "CreateClaimableBalance"))
		sb.WriteString(fmt.Sprintf(", Destinations=%s", getDestinationsString(o.Destinations)))
		sb.WriteString(fmt.Sprintf(", Amount=%s", o.Amount))

		asset := ""
		if o.Asset.IsNative() {
			asset = "native"
		} else {
			asset = fmt.Sprintf("%s:%s", o.Asset.GetCode(), o.Asset.GetIssuer())
		}
		sb.WriteString(fmt.Sprintf(", Asset=%s", asset))
		sb.WriteString("]")
	default:
		sb.WriteString(fmt.Sprintf("unrecognized_operation_type__%T", o))
	}
	return sb.String()
}

func getDestinationsString(destinations []txnbuild.Claimant) string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for i, d := range destinations {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("Claimant[")
		sb.WriteString(fmt.Sprintf("Desination=%s", d.Destination))
		sb.WriteString(fmt.Sprintf(", Predicate=%s", d.Predicate.Type.String()))
		sb.WriteString("]")
	}
	sb.WriteString("]")
	return sb.String()
}

package transform

import (
	"fmt"
	"strings"

	"github.com/stellar/go/txnbuild"
	"github.com/stellar/starbridge/model"
)

const (
	stellarDecimals               = 7
	ethereumNativeContractAddress = "0x0000000000000000000000000000000000000000"
)

// TODO need to set the contract account, source account, and seq numbers properly
var assetCode_WETH = "WETH"
var contractAccount = "GB42JR56FDOVUR75LN2J2F6DARS7SDUYMYPETQ24TDGRBCCQCHS2M2Y7"
var sourceAccount = contractAccount
var sourceLastSeqNum int64 = 0

// var contractSecretKey = "SAZGAQANN6UB3SM3GM7SF4PDF5EMC67LOHOYACK4O7VECYI2WTDI4F4P"
// var sourceSecretKey = contractSecretKey
// var destinationSecretKey = "SALR2RNJG55BBWTML2MKO5CXG5QDI4ZTSVDIA53XDWOU7QPOAEQNYUE2"

// TODO this destination needs to be input from the input transaction on the remote chain from the memo or similar
var destinationAccount = "GCBAA5476KARHPDSU6WFQTPXQOWX3QMXU4LF7JVZ2ZMWJ4OQEL7ZMV6G"
var baseFee int64 = 100

func getStellarAsset(tx *model.Transaction) txnbuild.CreditAsset {
	if tx.ContractAddress == ethereumNativeContractAddress {
		return txnbuild.CreditAsset{
			Code:   assetCode_WETH,
			Issuer: contractAccount,
		}
	}
	panic(fmt.Sprintf("unsupported contract address '%s' on tx", tx.ContractAddress))
}

func Transaction2Stellar(tx *model.Transaction) (*txnbuild.Transaction, error) {
	ops := []txnbuild.Operation{}
	ops = append(ops, &txnbuild.Payment{
		Destination: destinationAccount,
		Amount:      fmt.Sprintf("%d", tx.AmountUsingDecimals(stellarDecimals)),
		Asset:       getStellarAsset(tx),
		// SourceAccount: nil,
	})

	return txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: sourceAccount,
				Sequence:  sourceLastSeqNum,
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
	case *txnbuild.Payment:
		sb.WriteString("Operation[")
		sb.WriteString(fmt.Sprintf("SourceAccount=%s", op.GetSourceAccount()))
		sb.WriteString(fmt.Sprintf(", Type=%s", "Payment"))
		sb.WriteString(fmt.Sprintf(", Destination=%s", o.Destination))
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
		sb.WriteString("unrecognized_operation")
	}
	return sb.String()
}

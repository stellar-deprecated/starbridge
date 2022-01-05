package model

import (
	"fmt"
	"math"
	"strings"
)

type Transaction struct {
	Chain                ChainID
	Hash                 string
	Block                uint64
	SeqNum               uint64
	IsPending            bool
	ContractAddress      string
	From                 string
	To                   string
	Amount               uint64
	Decimals             int
	OriginalTx           interface{}
	AdditionalOriginalTx []interface{}
}

// String is the Stringer method
func (t Transaction) String() string {
	sb := strings.Builder{}
	sb.WriteString("Transaction[")
	sb.WriteString(fmt.Sprintf("ChainID=%s", t.Chain.String()))
	sb.WriteString(fmt.Sprintf(", Hash=%s", t.Hash))
	sb.WriteString(fmt.Sprintf(", Block=%d", t.Block))
	sb.WriteString(fmt.Sprintf(", SeqNum=%d", t.SeqNum))
	sb.WriteString(fmt.Sprintf(", IsPending=%v", t.IsPending))
	sb.WriteString(fmt.Sprintf(", ContractAddress=%s", t.ContractAddress))
	sb.WriteString(fmt.Sprintf(", From=%s", t.From))
	sb.WriteString(fmt.Sprintf(", To=%s", t.To))
	sb.WriteString(fmt.Sprintf(", Amount=%d", t.Amount))
	sb.WriteString(fmt.Sprintf(", Decimals=%d", t.Decimals))
	sb.WriteString(fmt.Sprintf(", HasOriginalTx=%v", t.OriginalTx != nil))
	sb.WriteString(fmt.Sprintf(", HasAdditionalOriginalTx=%v", t.AdditionalOriginalTx != nil && len(t.AdditionalOriginalTx) > 0))
	sb.WriteString("]")
	return sb.String()
}

// AmountUsingDecimals is a helper that converts decimal values for us
func (t Transaction) AmountUsingDecimals(newDecimals int) uint64 {
	var amountExponent int = newDecimals - t.Decimals
	if amountExponent > 0 {
		return t.Amount * uint64(math.Pow10(amountExponent))
	} else if amountExponent < 0 {
		return uint64(t.Amount / uint64(math.Pow10(-amountExponent)))
	}
	return t.Amount
}

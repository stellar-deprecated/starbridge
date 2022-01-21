package model

import (
	"fmt"
	"strings"
)

type Transaction struct {
	Chain                *Chain
	Hash                 string
	Block                uint64
	SeqNum               uint64
	IsPending            bool
	From                 string
	To                   string
	AssetInfo            *AssetInfo
	Amount               uint64
	OriginalTx           interface{}
	AdditionalOriginalTx []interface{}
	// TODO add fee information to the transaction
}

// String is the Stringer method
func (t Transaction) String() string {
	sb := strings.Builder{}
	sb.WriteString("Transaction[")
	sb.WriteString(fmt.Sprintf("Chain=%s", t.Chain.String()))
	sb.WriteString(fmt.Sprintf(", Hash=%s", t.Hash))
	sb.WriteString(fmt.Sprintf(", Block=%d", t.Block))
	sb.WriteString(fmt.Sprintf(", SeqNum=%d", t.SeqNum))
	sb.WriteString(fmt.Sprintf(", IsPending=%v", t.IsPending))
	sb.WriteString(fmt.Sprintf(", From=%s", t.From))
	sb.WriteString(fmt.Sprintf(", To=%s", t.To))
	sb.WriteString(fmt.Sprintf(", AssetInfo=%s", t.AssetInfo.String()))
	sb.WriteString(fmt.Sprintf(", Amount=%d", t.Amount))
	sb.WriteString(fmt.Sprintf(", HasOriginalTx=%v", t.OriginalTx != nil))
	sb.WriteString(fmt.Sprintf(", HasAdditionalOriginalTx=%v", t.AdditionalOriginalTx != nil && len(t.AdditionalOriginalTx) > 0))
	sb.WriteString("]")
	return sb.String()
}

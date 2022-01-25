package model

import (
	"fmt"
	"strings"
)

type ContractData struct {
	EventName                             string
	TargetDestinationChain                *Chain
	TargetDestinationAddressOnRemoteChain string
	AssetInfo                             *AssetInfo
	Amount                                uint64
}

// String is the Stringer method
func (c ContractData) String() string {
	sb := strings.Builder{}
	sb.WriteString("ContractData[")
	sb.WriteString(fmt.Sprintf("EventName=%s", c.EventName))
	sb.WriteString(fmt.Sprintf(", TargetDestinationChain=%s", c.TargetDestinationChain.String()))
	sb.WriteString(fmt.Sprintf(", TargetDestinationAddressOnRemoteChain=%s", c.TargetDestinationAddressOnRemoteChain))
	sb.WriteString(fmt.Sprintf(", AssetInfo=%s", c.AssetInfo.String()))
	sb.WriteString(fmt.Sprintf(", Amount=%d", c.Amount))
	sb.WriteString("]")
	return sb.String()
}

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
	Data                 ContractData // this is where the actual information lives about what needs to be transferred
	OriginalTx           interface{}
	AdditionalOriginalTx []interface{}
	// TODO NS add fee information to the transaction
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
	sb.WriteString(fmt.Sprintf(", Data=%s", t.Data.String()))
	sb.WriteString(fmt.Sprintf(", HasOriginalTx=%v", t.OriginalTx != nil))
	sb.WriteString(fmt.Sprintf(", HasAdditionalOriginalTx=%v", t.AdditionalOriginalTx != nil && len(t.AdditionalOriginalTx) > 0))
	sb.WriteString("]")
	return sb.String()
}

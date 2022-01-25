package model

import (
	"fmt"
	"log"

	"github.com/stellar/go/clients/horizonclient"
)

// Enum for Chain
type Chain struct {
	Name            string
	NativeAsset     *AssetInfo
	AddressMappings map[string]*AssetInfo // maps from a fixed set of assets from the remote chain to another fixed set of assets on the native chain (for now hard-coded, later on load from db)
	nextNonceFn     func(sourceAccount string) (uint64, error)
}

var (
	ChainStellar *Chain = &Chain{
		Name:        "Stellar",
		NativeAsset: AssetStellar_XLM,
		AddressMappings: map[string]*AssetInfo{
			AssetEthereum_ETH.MapKey():  AssetStellar_ETH,
			AssetEthereum_USDC.MapKey(): AssetStellar_WUSDC,
			AssetEthereum_WXLM.MapKey(): AssetStellar_XLM,
		},
		nextNonceFn: nextStellarNonceFn,
	}
	ChainEthereum *Chain = &Chain{
		Name:        "Ethereum",
		NativeAsset: AssetEthereum_ETH,
		AddressMappings: map[string]*AssetInfo{
			AssetStellar_ETH.MapKey():   AssetEthereum_ETH,
			AssetStellar_WUSDC.MapKey(): AssetEthereum_USDC,
			AssetStellar_XLM.MapKey():   AssetEthereum_WXLM,
		},
		nextNonceFn: unsupportedNonceForChain, // TODO we haven't had the time to add the logic to go from Stellar to Ethereum yet
	}
)

// String is the Stringer method
func (c *Chain) String() string {
	return c.Name
}

// NextNonce
func (c *Chain) NextNonce(sourceAccount string) (uint64, error) {
	return c.nextNonceFn(sourceAccount)
}

func nextStellarNonceFn(sourceAccount string) (uint64, error) {
	sdexAPI := horizonclient.DefaultTestNetClient

	log.Println("loading sequence number for Stellar")
	acctReq := horizonclient.AccountRequest{AccountID: sourceAccount}
	accountDetail, e := sdexAPI.AccountDetail(acctReq)
	if e != nil {
		return 0, fmt.Errorf("error loading account detail: %s", e)
	}
	seqNum, e := accountDetail.GetSequenceNumber()
	if e != nil {
		return 0, fmt.Errorf("error getting seq num: %s", e)
	}
	incrementedSeqNum := uint64(seqNum) + 1
	return incrementedSeqNum, nil
}

func unsupportedNonceForChain(sourceAccount string) (uint64, error) {
	return 0, fmt.Errorf("unsupported chain")
}

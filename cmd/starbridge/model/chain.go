package model

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"
)

// global for now
var logger *supportlog.Entry

// Enum for Chain
type Chain struct {
	Name                         string
	NativeAsset                  *AssetInfo
	AllAssetMap                  map[string]*AssetInfo // maps it's own assets (native + issued) so we can know what asset was used in a transaction being sent from this chain
	AddressMappings              map[string]*AssetInfo // maps from a fixed set of assets from the remote chain to another fixed set of assets on the native chain (for now hard-coded, later on load from db)
	nextNonceFn                  func(sourceAccount string) (uint64, error)
	ValidateDestinationAddressFn func(addr string) error
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
		nextNonceFn:                  nextStellarNonceFn,
		ValidateDestinationAddressFn: validateDestinationAddressFnStellar,
	}
	ChainEthereum *Chain = &Chain{
		Name:        "Ethereum",
		NativeAsset: AssetEthereum_ETH,
		AddressMappings: map[string]*AssetInfo{
			AssetStellar_ETH.MapKey():   AssetEthereum_ETH,
			AssetStellar_WUSDC.MapKey(): AssetEthereum_USDC,
			AssetStellar_XLM.MapKey():   AssetEthereum_WXLM,
		},
		nextNonceFn:                  unsupportedNonceForChain, // TODO NS we haven't had the time to add the logic to go from Stellar to Ethereum yet
		ValidateDestinationAddressFn: unsupportedValidateDestinationAddressFn,
	}
)

func computeAllAssetMap(chain *Chain) map[string]*AssetInfo {
	m := map[string]*AssetInfo{}
	for _, v := range chain.AddressMappings {
		m[v.MapKey()] = v
	}

	// native asset is included in the above list since it will be in the AddressMappings but add explicitly here too
	m[chain.NativeAsset.MapKey()] = chain.NativeAsset

	logger.Debugf("added %d items when creating AllAssetMap for chain=%s", len(m), chain.Name)
	return m
}

func init() {
	logger = supportlog.New()
	logger.SetLevel(logrus.InfoLevel)

	ChainStellar.AllAssetMap = computeAllAssetMap(ChainStellar)
	ChainEthereum.AllAssetMap = computeAllAssetMap(ChainEthereum)
}

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

	logger.Infof("loading sequence number for Stellar")
	acctReq := horizonclient.AccountRequest{AccountID: sourceAccount}
	accountDetail, err := sdexAPI.AccountDetail(acctReq)
	if err != nil {
		return 0, fmt.Errorf("error loading account detail: %s", err)
	}
	seqNum, err := accountDetail.GetSequenceNumber()
	if err != nil {
		return 0, fmt.Errorf("error getting seq num: %s", err)
	}
	return uint64(seqNum), nil
}

func unsupportedNonceForChain(sourceAccount string) (uint64, error) {
	return 0, fmt.Errorf("unsupported chain")
}

func validateDestinationAddressFnStellar(addr string) error {
	_, err := xdr.AddressToAccountId(addr)
	if err != nil {
		return fmt.Errorf("error parsing Stellar address %w", err)
	}
	return nil
}

func unsupportedValidateDestinationAddressFn(addr string) error {
	return fmt.Errorf("unsupported chain")
}

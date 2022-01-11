package model

import "fmt"

// Enum for Chain
type Chain struct {
	Name            string
	NativeAsset     *AssetInfo
	AddressMappings map[*AssetInfo]*AssetInfo // maps from a fixed set of assets from the remote chain to another fixed set of assets on the native chain (for now hard-coded, later on load from db)
}

var (
	ChainStellar *Chain = &Chain{Name: "Stellar", NativeAsset: AssetStellar_XLM, AddressMappings: map[*AssetInfo]*AssetInfo{
		AssetEthereum_ETH:  AssetStellar_ETH,
		AssetEthereum_WXLM: AssetStellar_XLM,
	}}
	ChainEthereum *Chain = &Chain{Name: "Ethereum", NativeAsset: AssetEthereum_ETH, AddressMappings: map[*AssetInfo]*AssetInfo{
		AssetStellar_ETH: AssetEthereum_ETH,
		AssetStellar_XLM: AssetEthereum_WXLM,
	}}
)

// String is the Stringer method
func (c *Chain) String() string {
	return c.Name
}

// NextNonce
func (c *Chain) NextNonce() uint64 {
	if c != ChainStellar {
		panic(fmt.Errorf("unsupported chain %s", c.Name))
	}
	// TODO need to set the seq numbers properly
	return 0
}

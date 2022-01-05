package model

import "fmt"

// Enum for Chain
type Chain struct {
	Name            string
	NativeAsset     *AssetInfo
	AddressMappings map[*AssetInfo]*AssetInfo // maps from a fixed set of assets to another fixed set of assets (for now hard-coded, later on load from db)
}

var (
	ChainStellar *Chain = &Chain{Name: "Stellar", NativeAsset: AssetXLM, AddressMappings: map[*AssetInfo]*AssetInfo{
		AssetETH:  AssetWETH,
		AssetWXLM: AssetXLM,
	}}
	ChainEthereum *Chain = &Chain{Name: "Ethereum", NativeAsset: AssetETH, AddressMappings: map[*AssetInfo]*AssetInfo{
		AssetWETH: AssetETH,
		AssetXLM:  AssetWXLM,
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

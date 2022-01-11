package model

import "fmt"

// AssetInfo represents an asset with all information needed to use it
// For now, we consciously do not encode the Chain it is native to so we can avoid a circular dependency
type AssetInfo struct {
	Code            string
	ContractAddress string
	Decimals        int
}

var (
	AssetEthereum_ETH *AssetInfo = &AssetInfo{
		Code:            "ETH",
		ContractAddress: "0x0000000000000000000000000000000000000000",
		Decimals:        18,
	}
	AssetStellar_XLM *AssetInfo = &AssetInfo{
		Code:            "XLM",
		ContractAddress: "native",
		Decimals:        7,
	}

	// wrapped assets
	// TODO need to set contract account and key correctly
	AssetStellar_ETH *AssetInfo = &AssetInfo{
		Code:            "ETH",                                                      // Stellar assets will retain the original assetCode and will not have the W prefix
		ContractAddress: "GB42JR56FDOVUR75LN2J2F6DARS7SDUYMYPETQ24TDGRBCCQCHS2M2Y7", // contractSecret: SAZGAQANN6UB3SM3GM7SF4PDF5EMC67LOHOYACK4O7VECYI2WTDI4F4P
		Decimals:        7,
	}
	// TODO need to set contract account and key correctly
	AssetEthereum_WXLM *AssetInfo = &AssetInfo{
		Code:            "WXLM",
		ContractAddress: "0x0000000000000000000000000000000123456789", // contractSecret ??
		Decimals:        7,
	}
)

// String is the Stringer method
func (a AssetInfo) String() string {
	return fmt.Sprintf("AssetInfo[Code=%s, ContractAddress=%s, Decimals=%d]", a.Code, a.ContractAddress, a.Decimals)
}

package model

import "fmt"

var stellarEscrowAccount = "GBNV7CTQAJTSYJE4BTA76JF2GZ6UA6WRN3KN55GUY2K62XIYD4ZADID5" // var escrowSecretKey = "SABNONONIEROOG7JURODU56QHBBI4SYGYEZ7I432QPV4THZPHRSACIBF"

// AssetInfo represents an asset with all information needed to use it
// For now, we consciously do not encode the Chain it is native to so we can avoid a circular dependency
type AssetInfo struct {
	Code            string
	ContractAddress string
	Decimals        int
	mapKey          func(a *AssetInfo) string
}

var (
	// native assets on Ethereum chain
	AssetEthereum_USDC *AssetInfo = &AssetInfo{
		Code:            "USDC",
		ContractAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		Decimals:        6,
		mapKey:          ethereumMapKey,
	}
	AssetEthereum_ETH *AssetInfo = &AssetInfo{
		Code:            "ETH",
		ContractAddress: "0x0000000000000000000000000000000000000000",
		Decimals:        18,
		mapKey:          ethereumMapKey,
	}

	// native assets on Stellar chain
	AssetStellar_XLM *AssetInfo = &AssetInfo{
		Code:            "XLM",
		ContractAddress: "native",
		Decimals:        7,
		mapKey:          stellarMapKey,
	}

	// wrapped assets on Ethereum
	AssetEthereum_WXLM *AssetInfo = &AssetInfo{
		// TODO need to set contract account and key correctly
		Code:            "WXLM",
		ContractAddress: "0x0000000000000000000000000000000123456789", // contractSecret ??
		Decimals:        7,
		mapKey:          ethereumMapKey,
	}

	// wrapped assets on Stellar
	AssetStellar_ETH *AssetInfo = &AssetInfo{
		// TODO need to set contract account and key correctly
		Code:            "ETH",                // Stellar assets will retain the original assetCode and will not have the W prefix
		ContractAddress: stellarEscrowAccount, // this is the escrow account
		Decimals:        7,
		mapKey:          stellarMapKey,
	}
	AssetStellar_WUSDC *AssetInfo = &AssetInfo{
		// TODO need to set contract account and key correctly
		Code:            "WUSDC",              // since we have a USDC on Stellar, we add the "W" prefix for now
		ContractAddress: stellarEscrowAccount, // this is the escrow account
		Decimals:        7,
		mapKey:          stellarMapKey,
	}
)

// String is the Stringer method
func (a AssetInfo) String() string {
	return fmt.Sprintf("AssetInfo[Code=%s, ContractAddress=%s, Decimals=%d]", a.Code, a.ContractAddress, a.Decimals)
}

// MapKey returns a string that can be used to uniquely identify this asset in a map as the key
func (a AssetInfo) MapKey() string {
	return a.mapKey(&a)
}

func ethereumMapKey(a *AssetInfo) string {
	return a.ContractAddress
}

func stellarMapKey(a *AssetInfo) string {
	if a.Code == "native" {
		return a.Code
	}
	return fmt.Sprintf("%s:%s", a.Code, a.ContractAddress)
}

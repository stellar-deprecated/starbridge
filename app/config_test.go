package app

import (
	"testing"

	"github.com/stellar/go/support/config"
	"github.com/stellar/starbridge/backend"

	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	var cfg Config
	err := config.Read("./testdata/example.cfg", &cfg)
	require.NoError(t, err)

	expected := Config{
		Port:                  8000,
		AdminPort:             6060,
		PostgresDSN:           "dbname=starbridge user=starbridge",
		HorizonURL:            "https://horizon-testnet.stellar.org",
		NetworkPassphrase:     "Test SDF Network ; September 2015",
		StellarBridgeAccount:  "GAJKCRY6CIOXRIVK55ALOOJA327XN4JZ5KKN7YCTT3WM5W6BMFXMVQC2",
		StellarPrivateKey:     "SCSTO3PMPM2BNLR2MYKVHWCJ2FNHQGFWKPOFH6UX4N3HO6HMK4JBSJ6F",
		EthereumRPCURL:        "https://ethereum-goerli-rpc.allthatnode.com",
		EthereumBridgeAddress: "0xD0675839A6C2c3412a3026Aa5F521Ea1e948E526",
		EthereumPrivateKey:    "2aecee1800342bae06228ed990a152563b8dedf5fe15e3eab4b44854c9e001e5",
		AssetMapping: []backend.AssetMappingConfigEntry{
			{
				StellarAsset:      "ETH:GAJKCRY6CIOXRIVK55ALOOJA327XN4JZ5KKN7YCTT3WM5W6BMFXMVQC2",
				EthereumToken:     "0x0000000000000000000000000000000000000000",
				StellarToEthereum: "100000000000",
			},
			{
				StellarAsset:      "native",
				EthereumToken:     "0x23896e5E10363e4a90573abDd405Ab9761E6cCE2",
				StellarToEthereum: "1",
			},
		},
	}
	require.Equal(t, expected, cfg)
}

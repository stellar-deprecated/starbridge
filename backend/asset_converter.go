package backend

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/stellar/go/amount"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"

	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/go/xdr"
)

var (
	WithdrawalAssetInvalid = problem.P{
		Type:   "withdrawal_asset_invalid",
		Title:  "Withdrawal Asset Invalid",
		Status: http.StatusBadRequest,
		Detail: "Withdrawing the requested asset is not supported by the bridge." +
			"Refund the deposit once the withdrawal period has expired.",
	}
	WithdrawalAmountInvalid = problem.P{
		Type:   "withdrawal_amount_invalid",
		Title:  "Withdrawal Amount Invalid",
		Status: http.StatusBadRequest,
		Detail: "Withdrawing the requested amount is not supported by the bridge." +
			"Refund the deposit once the withdrawal period has expired.",
	}
)

// AssetMappingConfigEntry is the toml representation of
// a mapping between a Stellar asset and an Ethereum token
type AssetMappingConfigEntry struct {
	StellarAsset        string `toml:"stellar_asset" valid:"-"`
	EthereumToken       string `toml:"ethereum_token" valid:"-"`
	ConcordiumToken     string `toml:"concordium_token" valid:"-"`
	OkxToken            string `toml:"okx_token" valid:"-"`
	StellarToEthereum   string `toml:"stellar_to_ethereum" valid:"-"`
	StellarToOkx        string `toml:"stellar_to_okx" valid:"-"`
	StellarToConcordium string `toml:"stellar_to_concordium" valid:"-"`
}

type stellarRate struct {
	asset string
	rate  *big.Rat
}

type ethereumRate struct {
	token common.Address
	rate  *big.Rat
}

type okxRate struct {
	token common.Address
	rate  *big.Rat
}

type concordiumRate struct {
	token string
	rate  *big.Rat
}

// AssetConverter maps assets from Stellar to their
// equivalent tokens on Ethereum and vice versa.
type AssetConverter struct {
	ethereumToStellar   map[common.Address]stellarRate
	stellarToEthereum   map[string]ethereumRate
	okxToStellar        map[common.Address]stellarRate
	stellarToOkx        map[string]okxRate
	concordiumToStellar map[string]stellarRate
	stellarToConcordium map[string]concordiumRate
}

func isAsset(assetString string) bool {
	var asset xdr.Asset

	if strings.ToLower(assetString) == "native" {
		return true
	} else {

		parts := strings.Split(assetString, ":")
		if len(parts) != 2 {
			return false
		}

		code := parts[0]
		if !xdr.ValidAssetCode.MatchString(code) {
			return false
		}

		issuer, err := xdr.AddressToAccountId(parts[1])
		if err != nil {
			return false
		}

		if err := asset.SetCredit(code, issuer); err != nil {
			return false
		}
	}

	return true
}

// NewAssetConverter constructs a new instance of AssetConverter
func NewAssetConverter(configEntries []AssetMappingConfigEntry) (AssetConverter, error) {
	converter := AssetConverter{
		ethereumToStellar:   map[common.Address]stellarRate{},
		stellarToEthereum:   map[string]ethereumRate{},
		okxToStellar:        map[common.Address]stellarRate{},
		stellarToOkx:        map[string]okxRate{},
		concordiumToStellar: map[string]stellarRate{},
		stellarToConcordium: map[string]concordiumRate{},
	}

	if len(configEntries) == 0 {
		return converter, fmt.Errorf("config entries are empty")
	}

	for _, entry := range configEntries {
		if !isAsset(entry.StellarAsset) {
			return converter, fmt.Errorf("%s is not a valid stellar asset", entry.StellarAsset)
		}
		if !common.IsHexAddress(entry.EthereumToken) {
			return converter, fmt.Errorf("%s is not a valid ethereum address", entry.EthereumToken)
		}
		multiplierEth, ok := new(big.Int).SetString(entry.StellarToEthereum, 10)
		if !ok {
			return converter, fmt.Errorf("%s is not a valid multiplier", entry.StellarToEthereum)
		}
		if !common.IsHexAddress(entry.OkxToken) {
			return converter, fmt.Errorf("%s is not a valid okx address", entry.OkxToken)
		}
		multiplierOkx, ok := new(big.Int).SetString(entry.StellarToOkx, 10)
		if !ok {
			return converter, fmt.Errorf("%s is not a valid multiplier", entry.StellarToOkx)
		}
		multiplierCcd, ok := new(big.Int).SetString(entry.StellarToConcordium, 10)
		if !ok {
			return converter, fmt.Errorf("%s is not a valid multiplier", entry.StellarToConcordium)
		}
		token := common.HexToAddress(entry.EthereumToken)
		okxToken := common.HexToAddress(entry.OkxToken)
		_, exists := converter.stellarToEthereum[entry.StellarAsset]
		if exists {
			return converter, fmt.Errorf("asset %v is repeated in the asset mapping ", entry.StellarAsset)
		}
		_, exists = converter.ethereumToStellar[token]
		if exists {
			return converter, fmt.Errorf("token %v is repeated in the asset mapping ", entry.EthereumToken)
		}
		_, exists = converter.stellarToOkx[entry.StellarAsset]
		if exists {
			return converter, fmt.Errorf("asset %v is repeated in the asset mapping ", entry.StellarAsset)
		}
		_, exists = converter.okxToStellar[okxToken]
		if exists {
			return converter, fmt.Errorf("token %v is repeated in the asset mapping ", entry.OkxToken)
		}
		_, exists = converter.stellarToConcordium[entry.StellarAsset]
		if exists {
			return converter, fmt.Errorf("asset %v is repeated in the asset mapping ", entry.StellarAsset)
		}
		_, exists = converter.concordiumToStellar[entry.ConcordiumToken]
		if exists {
			return converter, fmt.Errorf("token %v is repeated in the asset mapping ", entry.ConcordiumToken)
		}
		converter.stellarToEthereum[entry.StellarAsset] = ethereumRate{
			token: token,
			rate:  new(big.Rat).SetFrac(multiplierEth, big.NewInt(1)),
		}
		converter.ethereumToStellar[token] = stellarRate{
			asset: entry.StellarAsset,
			rate:  new(big.Rat).SetFrac(big.NewInt(1), multiplierEth),
		}
		converter.stellarToOkx[entry.StellarAsset] = okxRate{
			token: okxToken,
			rate:  new(big.Rat).SetFrac(multiplierOkx, big.NewInt(1)),
		}
		converter.okxToStellar[okxToken] = stellarRate{
			asset: entry.StellarAsset,
			rate:  new(big.Rat).SetFrac(big.NewInt(1), multiplierOkx),
		}
		converter.stellarToConcordium[entry.StellarAsset] = concordiumRate{
			token: entry.ConcordiumToken,
			rate:  new(big.Rat).SetFrac(multiplierCcd, big.NewInt(1)),
		}
		converter.concordiumToStellar[entry.ConcordiumToken] = stellarRate{
			asset: entry.StellarAsset,
			rate:  new(big.Rat).SetFrac(big.NewInt(1), multiplierCcd),
		}
	}

	return converter, nil
}

// ToStellar returns the Stellar asset and amount for the given Ethereum token
func (c AssetConverter) ToStellar(token string, tokenAmount string, fromEth bool, fromCcd bool, fromOkx bool) (string, int64, error) {
	//if !common.IsHexAddress(token) {
	//	return "", 0, WithdrawalAssetInvalid
	//}

	parsedAmount := &big.Int{}
	_, ok := parsedAmount.SetString(tokenAmount, 10)
	if !ok || parsedAmount.Cmp(big.NewInt(0)) <= 0 {
		return "", 0, WithdrawalAmountInvalid
	}

	entryEth, okEth := c.ethereumToStellar[common.HexToAddress(token)]
	entryOkx, okOkx := c.okxToStellar[common.HexToAddress(token)]
	entryCcd, okCcd := c.concordiumToStellar[token]
	entry := stellarRate{}
	if !okEth && !okCcd && !okOkx {
		return "", 0, WithdrawalAssetInvalid
	}
	if fromEth {
		entry = entryEth
	} else if fromCcd {
		entry = entryCcd
	} else if fromOkx {
		entry = entryOkx
	} else {
		return "", 0, WithdrawalAssetInvalid
	}

	product := new(big.Rat).Mul(new(big.Rat).SetFrac(parsedAmount, big.NewInt(1)), entry.rate)
	if product.IsInt() {
		val := product.Num().Int64()
		if product.Num().IsInt64() && val > 0 {
			return entry.asset, val, nil
		}
	}

	return entry.asset, 0, WithdrawalAmountInvalid
}

// ToEthereum returns the Ethereum token and amount for the given Stellar asset
func (c AssetConverter) ToEthereum(asset string, assetAmount string) (common.Address, *big.Int, error) {
	entry, ok := c.stellarToEthereum[asset]
	if !ok {
		return common.Address{}, nil, WithdrawalAssetInvalid
	}

	parsedAmount, err := amount.ParseInt64(assetAmount)
	if err != nil {
		return common.Address{}, nil, WithdrawalAssetInvalid
	}

	product := new(big.Rat).Mul(new(big.Rat).SetFrac(big.NewInt(parsedAmount), big.NewInt(1)), entry.rate)
	if product.IsInt() {
		val := product.Num()
		if val.Cmp(big.NewInt(0)) <= 0 || val.Cmp(math.MaxBig256) > 0 {
			return entry.token, nil, WithdrawalAmountInvalid
		}
		return entry.token, val, nil
	}

	return entry.token, nil, WithdrawalAmountInvalid
}

// ToOkx returns the Okx token and amount for the given Stellar asset
func (c AssetConverter) ToOkx(asset string, assetAmount string) (common.Address, *big.Int, error) {
	entry, ok := c.stellarToOkx[asset]
	if !ok {
		return common.Address{}, nil, WithdrawalAssetInvalid
	}

	parsedAmount, err := amount.ParseInt64(assetAmount)
	if err != nil {
		return common.Address{}, nil, WithdrawalAssetInvalid
	}

	product := new(big.Rat).Mul(new(big.Rat).SetFrac(big.NewInt(parsedAmount), big.NewInt(1)), entry.rate)
	if product.IsInt() {
		val := product.Num()
		if val.Cmp(big.NewInt(0)) <= 0 || val.Cmp(math.MaxBig256) > 0 {
			return entry.token, nil, WithdrawalAmountInvalid
		}
		return entry.token, val, nil
	}

	return entry.token, nil, WithdrawalAmountInvalid
}

func (c AssetConverter) ToConcordium(asset string, assetAmount string) (string, *big.Int, error) {
	entry, ok := c.stellarToConcordium[asset]
	if !ok {
		return "", nil, WithdrawalAssetInvalid
	}

	parsedAmount, err := amount.ParseInt64(assetAmount)
	if err != nil {
		return "", nil, WithdrawalAssetInvalid
	}

	product := new(big.Rat).Mul(new(big.Rat).SetFrac(big.NewInt(parsedAmount), big.NewInt(1)), entry.rate)
	if product.IsInt() {
		val := product.Num()
		if val.Cmp(big.NewInt(0)) <= 0 || val.Cmp(math.MaxBig256) > 0 {
			return entry.token, nil, WithdrawalAmountInvalid
		}
		return entry.token, val, nil
	}

	return entry.token, nil, WithdrawalAmountInvalid
}

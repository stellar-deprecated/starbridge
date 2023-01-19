package controllers

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"

	"github.com/stellar/go/amount"
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
	StellarAssetContractID string `toml:"stellar_asset_contract_id" valid:"-"`
	EthereumToken          string `toml:"ethereum_token" valid:"-"`
	StellarToEthereum      string `toml:"stellar_to_ethereum" valid:"-"`
}

type stellarRate struct {
	asset [32]byte
	rate  *big.Rat
}

type ethereumRate struct {
	token common.Address
	rate  *big.Rat
}

// AssetConverter maps assets from Stellar to their
// equivalent tokens on Ethereum and vice versa.
type AssetConverter struct {
	ethereumToStellar map[common.Address]stellarRate
	stellarToEthereum map[[32]byte]ethereumRate
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
		ethereumToStellar: map[common.Address]stellarRate{},
		stellarToEthereum: map[[32]byte]ethereumRate{},
	}

	if len(configEntries) == 0 {
		return converter, fmt.Errorf("config entries are empty")
	}

	for _, entry := range configEntries {
		contractIDBytes, err := hex.DecodeString(entry.StellarAssetContractID)
		if err != nil || len(contractIDBytes) != 32 {
			return converter, fmt.Errorf("%s is not a valid stellar asset", entry.StellarAssetContractID)
		}
		var contractID [32]byte
		copy(contractID[:], contractIDBytes)
		if !common.IsHexAddress(entry.EthereumToken) {
			return converter, fmt.Errorf("%s is not a valid ethereum address", entry.EthereumToken)
		}
		multiplier, ok := new(big.Int).SetString(entry.StellarToEthereum, 10)
		if !ok {
			return converter, fmt.Errorf("%s is not a valid multiplier", entry.StellarToEthereum)
		}
		token := common.HexToAddress(entry.EthereumToken)
		_, exists := converter.stellarToEthereum[contractID]
		if exists {
			return converter, fmt.Errorf("asset %v is repeated in the asset mapping ", entry.StellarAssetContractID)
		}
		_, exists = converter.ethereumToStellar[token]
		if exists {
			return converter, fmt.Errorf("token %v is repeated in the asset mapping ", entry.EthereumToken)
		}
		converter.stellarToEthereum[contractID] = ethereumRate{
			token: token,
			rate:  new(big.Rat).SetFrac(multiplier, big.NewInt(1)),
		}
		converter.ethereumToStellar[token] = stellarRate{
			asset: contractID,
			rate:  new(big.Rat).SetFrac(big.NewInt(1), multiplier),
		}
	}

	return converter, nil
}

// ToStellar returns the Stellar asset and amount for the given Ethereum token
func (c AssetConverter) ToStellar(token string, tokenAmount string) ([32]byte, int64, error) {
	if !common.IsHexAddress(token) {
		return [32]byte{}, 0, WithdrawalAssetInvalid
	}

	parsedAmount := &big.Int{}
	_, ok := parsedAmount.SetString(tokenAmount, 10)
	if !ok || parsedAmount.Cmp(big.NewInt(0)) <= 0 {
		return [32]byte{}, 0, WithdrawalAmountInvalid
	}

	entry, ok := c.ethereumToStellar[common.HexToAddress(token)]
	if !ok {
		return [32]byte{}, 0, WithdrawalAssetInvalid
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
func (c AssetConverter) ToEthereum(assetContractID [32]byte, assetAmount string) (common.Address, *big.Int, error) {
	entry, ok := c.stellarToEthereum[assetContractID]
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

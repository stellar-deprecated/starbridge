package controllers

import (
	"crypto/sha256"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
	"net/http"
	"strings"

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
	StellarAsset      string `toml:"stellar_asset" valid:"-"`
	EthereumToken     string `toml:"ethereum_token" valid:"-"`
	StellarToEthereum string `toml:"stellar_to_ethereum" valid:"-"`
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
	ethereumToStellar   map[common.Address]stellarRate
	stellarToEthereum   map[[32]byte]ethereumRate
	contractIDIsWrapped map[[32]byte]bool
}

func assetToXDR(assetString string) (xdr.Asset, error) {
	if strings.ToLower(assetString) == "native" {
		return xdr.MustNewNativeAsset(), nil
	}

	var asset xdr.Asset
	parts := strings.Split(assetString, ":")
	if len(parts) != 2 {
		return asset, fmt.Errorf("asset has too many colons")
	}

	code := parts[0]
	if !xdr.ValidAssetCode.MatchString(code) {
		return asset, fmt.Errorf("asset code is invalid")
	}

	issuer, err := xdr.AddressToAccountId(parts[1])
	if err != nil {
		return asset, fmt.Errorf("asset issure is invalid")
	}

	if err := asset.SetCredit(code, issuer); err != nil {
		return asset, err
	}
	return asset, nil
}

func stellarAssetContractID(passPhrase string, asset xdr.Asset) ([32]byte, error) {
	networkId := xdr.Hash(sha256.Sum256([]byte(passPhrase)))
	preImage := xdr.HashIdPreimage{
		Type: xdr.EnvelopeTypeEnvelopeTypeContractIdFromAsset,
		FromAsset: &xdr.HashIdPreimageFromAsset{
			NetworkId: networkId,
			Asset:     asset,
		},
	}
	xdrPreImageBytes, err := preImage.MarshalBinary()
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(xdrPreImageBytes), nil
}

// NewAssetConverter constructs a new instance of AssetConverter
func NewAssetConverter(passPhrase string, bridgeAccount string, configEntries []AssetMappingConfigEntry) (AssetConverter, error) {
	converter := AssetConverter{
		ethereumToStellar:   map[common.Address]stellarRate{},
		stellarToEthereum:   map[[32]byte]ethereumRate{},
		contractIDIsWrapped: map[[32]byte]bool{},
	}

	if len(configEntries) == 0 {
		return converter, fmt.Errorf("config entries are empty")
	}

	for _, entry := range configEntries {
		assetXDR, err := assetToXDR(entry.StellarAsset)
		if err != nil {
			return converter, err
		}
		assetContractID, err := stellarAssetContractID(passPhrase, assetXDR)
		if err != nil {
			return converter, err
		}
		converter.contractIDIsWrapped[assetContractID] = assetXDR.GetIssuer() == bridgeAccount
		if !common.IsHexAddress(entry.EthereumToken) {
			return converter, fmt.Errorf("%s is not a valid ethereum address", entry.EthereumToken)
		}
		multiplier, ok := new(big.Int).SetString(entry.StellarToEthereum, 10)
		if !ok {
			return converter, fmt.Errorf("%s is not a valid multiplier", entry.StellarToEthereum)
		}
		token := common.HexToAddress(entry.EthereumToken)
		_, exists := converter.stellarToEthereum[assetContractID]
		if exists {
			return converter, fmt.Errorf("asset %v is repeated in the asset mapping ", entry.StellarAsset)
		}
		_, exists = converter.ethereumToStellar[token]
		if exists {
			return converter, fmt.Errorf("token %v is repeated in the asset mapping ", entry.EthereumToken)
		}
		converter.stellarToEthereum[assetContractID] = ethereumRate{
			token: token,
			rate:  new(big.Rat).SetFrac(multiplier, big.NewInt(1)),
		}
		converter.ethereumToStellar[token] = stellarRate{
			asset: assetContractID,
			rate:  new(big.Rat).SetFrac(big.NewInt(1), multiplier),
		}
	}

	return converter, nil
}

// ToStellar returns the Stellar asset and amount for the given Ethereum token
func (c AssetConverter) ToStellar(token string, tokenAmount string) ([32]byte, bool, int64, error) {
	if !common.IsHexAddress(token) {
		return [32]byte{}, false, 0, WithdrawalAssetInvalid
	}

	parsedAmount := &big.Int{}
	_, ok := parsedAmount.SetString(tokenAmount, 10)
	if !ok || parsedAmount.Cmp(big.NewInt(0)) <= 0 {
		return [32]byte{}, false, 0, WithdrawalAmountInvalid
	}

	entry, ok := c.ethereumToStellar[common.HexToAddress(token)]
	if !ok {
		return [32]byte{}, false, 0, WithdrawalAssetInvalid
	}

	product := new(big.Rat).Mul(new(big.Rat).SetFrac(parsedAmount, big.NewInt(1)), entry.rate)
	if product.IsInt() {
		val := product.Num().Int64()
		if product.Num().IsInt64() && val > 0 {
			return entry.asset, c.contractIDIsWrapped[entry.asset], val, nil
		}
	}

	return entry.asset, false, 0, WithdrawalAmountInvalid
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

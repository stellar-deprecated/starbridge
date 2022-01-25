package integrations

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stellar/starbridge/cmd/starbridge/model"
)

// temporarily used this contract address, which has a tx (0xf08debd774d2ecf0a18f62d593f0ec3af2aabd373139d87b58e3d28e088c2b59) through which
// a USDC payment was sent from one address to another by interacting with this contract directly (the USDC contract).
// i.e. this contract address is the USDC contract on Ethereum
// var MY_ETHEREUM_CONTRACT_ADDRESS = strings.ToLower("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

// this is something I deployed on the ropsten test network using the solidity Remix IDE - see the contract in ../../../contracts/SimpleEscrowEvents/simple_escrow_events.sol
var MY_ETHEREUM_CONTRACT_ADDRESS = strings.ToLower("0xa3f92C4Af3eA1189777001b654a7a9E8ae4C90b3")

// IsMyContractAddress returns true only iff the passed in address is the same as MY_ETHEREUM_CONTRACT_ADDRESS
func IsMyContractAddress(otherAddress string) bool {
	return strings.ToLower(otherAddress) == MY_ETHEREUM_CONTRACT_ADDRESS
}

var MY_CONTRACT_METHOD_HASH_SELECTORS = map[string]string{}

type PaymentEvent struct {
	DestinationStellarAddress string
	TokenContractAddress      string
	TokenAmount               *big.Int
}

const eventName = "Payment"

var ethContractAddressHash = fmt.Sprintf("0x%s", common.Bytes2Hex(crypto.Keccak256([]byte("0x0000000000000000000000000000000000000000"))))
var usdcContractAddressHash = fmt.Sprintf("0x%s", common.Bytes2Hex(crypto.Keccak256([]byte("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"))))
var payableAsset = model.AssetEthereum_ETH // TODO this is hard-coded to ETH right now since our smart-contract only allows sending that asset to our smart-contract for now

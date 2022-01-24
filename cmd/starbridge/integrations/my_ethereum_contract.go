package integrations

import (
	"encoding/hex"
	"log"
	"strings"

	"golang.org/x/crypto/sha3"
)

// temporarily used this contract address, which has a tx (0xf08debd774d2ecf0a18f62d593f0ec3af2aabd373139d87b58e3d28e088c2b59) through which
// a USDC payment was sent from one address to another by interacting with this contract directly (the USDC contract).
// i.e. this contract address is the USDC contract on Ethereum
// var MY_ETHEREUM_CONTRACT_ADDRESS = strings.ToLower("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

// this is something I deployed on the ropsten test network using the solidity Remix IDE - see the contract in ../../../contracts/SimpleEscrowEvents/simple_escrow_events.sol
var MY_ETHEREUM_CONTRACT_ADDRESS = strings.ToLower("0x3F1F0b8Bc207F3a94A7Fc72be07B30fB28319D5d")

// IsMyContractAddress returns true only iff the passed in address is the same as MY_ETHEREUM_CONTRACT_ADDRESS
func IsMyContractAddress(otherAddress string) bool {
	return strings.ToLower(otherAddress) == MY_ETHEREUM_CONTRACT_ADDRESS
}

var MY_CONTRACT_METHOD_HASH_SELECTORS = map[string]string{}

func init() {
	fnSigs := []string{
		"incrementCounter()",
		"decrementCounter()",
		"getCount()",
	}
	for _, fnSig := range fnSigs {
		h := sha3.New512()
		h.Write([]byte(fnSig))
		sha3Sum := h.Sum(nil)

		key := hex.EncodeToString([]byte(sha3Sum))[:8]
		MY_CONTRACT_METHOD_HASH_SELECTORS[key] = fnSig

		log.Printf("added value to MY_CONTRACT_METHOD_HASH_SELECTORS for key '%s' = '%s'", key, fnSig)
	}
}

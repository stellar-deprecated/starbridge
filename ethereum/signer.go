package ethereum

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/solidity-go"
)

var (
	bytes32           = mustType("bytes32")
	withdrawERC20Type = mustTupleType([]abi.ArgumentMarshaling{
		{Name: "id", Type: "bytes32"},
		{Name: "expiration", Type: "uint256"},
		{Name: "recipient", Type: "address"},
		{Name: "token", Type: "address"},
		{Name: "amount", Type: "uint256"},
	})
	withdrawETHType = mustTupleType([]abi.ArgumentMarshaling{
		{Name: "id", Type: "bytes32"},
		{Name: "expiration", Type: "uint256"},
		{Name: "recipient", Type: "address"},
		{Name: "amount", Type: "uint256"},
	})
)

func mustType(t string) abi.Type {
	ty, err := abi.NewType(t, t, nil)
	if err != nil {
		log.Fatalf("%v invalid type %v", t, err)
	}
	return ty
}

func mustTupleType(components []abi.ArgumentMarshaling) abi.Type {
	ty, err := abi.NewType("tuple", "", components)
	if err != nil {
		log.Fatalf("invalid type %v", err)
	}
	return ty
}

// Signer represents an ethereum validator account which is
// authorized to approve withdrawals from the bridge smart contract.
type Signer struct {
	privateKey      *ecdsa.PrivateKey
	domainSeparator [32]byte
	address         common.Address
}

// NewSigner constructs a new Signer instance
func NewSigner(privateKey string, domainSeparator [32]byte) (Signer, error) {
	parsed, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return Signer{}, err
	}
	return Signer{
		privateKey:      parsed,
		domainSeparator: domainSeparator,
		address:         crypto.PubkeyToAddress(parsed.PublicKey),
	}, nil
}

// Address returns the ethereum address corresponding to the public
// key of the signer
func (s Signer) Address() common.Address {
	return s.address
}

// SignWithdrawal returns a signature for the given withdrawal request
func (s Signer) SignWithdrawal(
	id common.Hash,
	expiration int64,
	recipient,
	token common.Address, // an address of 0x0 indicates an ETH transfer
	amount *big.Int,
) ([]byte, error) {
	if token == (common.Address{}) {
		return s.signWithdrawETHRequest(solidity.WithdrawETHRequest{
			Id:         id,
			Expiration: big.NewInt(expiration),
			Recipient:  recipient,
			Amount:     amount,
		})
	} else {
		return s.signWithdrawERC20Request(solidity.WithdrawERC20Request{
			Id:         id,
			Expiration: big.NewInt(expiration),
			Recipient:  recipient,
			Token:      token,
			Amount:     amount,
		})
	}
}

func (s Signer) signWithdrawERC20Request(request solidity.WithdrawERC20Request) ([]byte, error) {
	arguments := abi.Arguments{
		{Type: bytes32},
		{Type: bytes32},
		{Type: withdrawERC20Type},
	}

	abiEncoded, err := arguments.Pack(
		s.domainSeparator,
		crypto.Keccak256Hash([]byte("withdrawERC20")),
		request,
	)
	if err != nil {
		return nil, err
	}
	return s.signPayload(abiEncoded)
}

func (s Signer) signWithdrawETHRequest(request solidity.WithdrawETHRequest) ([]byte, error) {
	arguments := abi.Arguments{
		{Type: bytes32},
		{Type: bytes32},
		{Type: withdrawETHType},
	}

	abiEncoded, err := arguments.Pack(
		s.domainSeparator,
		crypto.Keccak256Hash([]byte("withdrawETH")),
		request,
	)
	if err != nil {
		return nil, err
	}
	return s.signPayload(abiEncoded)
}

func (s Signer) signPayload(abiEncoded []byte) ([]byte, error) {
	sig, err := crypto.Sign(accounts.TextHash(crypto.Keccak256(abiEncoded)), s.privateKey)
	if err != nil {
		return nil, err
	}
	// The ECDSA solidity library used by the bridge smart contract expects the v
	// value to be 27 or 28, see:
	// https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.7.0/contracts/utils/cryptography/ECDSA.sol#L41-L43
	// However, crypto.Sign() encodes the v value as 0 or 1, see:
	// https://github.com/ethereum/go-ethereum/blob/v1.10.20/crypto/signature_cgo.go#L54
	// That is why we need to transform the signature to be compatible with the
	// ECDSA solidity library.
	sig[len(sig)-1] += 27
	return sig, nil
}

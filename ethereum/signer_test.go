package ethereum

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func createSigner(t *testing.T) Signer {
	signer, err := NewSigner(
		"51138e68e8a5fa906d38c5b42bc01b805d7adb3fce037743fb406bb10aa83307",
		[32]byte{1, 2, 3},
	)
	assert.NoError(t, err)
	return signer
}

func TestSigner_Address(t *testing.T) {
	signer := createSigner(t)
	assert.Equal(
		t,
		common.HexToAddress("0xAe1B35129e5924C3a7313EE579878f829f3e8495"),
		signer.Address(),
	)
}

func TestSigner_SignWithdrawal(t *testing.T) {
	signer := createSigner(t)
	for _, withdrawal := range []struct {
		id         common.Hash
		expiration int64
		recipient  common.Address
		token      common.Address
		amount     *big.Int
		expected   string
	}{
		{
			id:         common.HexToHash("0x99"),
			expiration: 100,
			recipient:  common.HexToAddress("0x123"),
			token:      common.HexToAddress("0x456"),
			amount:     big.NewInt(200),
			expected:   "d668c6d190f0a1dcb03b5540794479659a5d46cfa741d2e52f65b9d5e4afae420dae19570a2dd871486c886c0e19511eb1bb39299baa99baf7b73ef30190e0d91c",
		},
		{
			id:         common.HexToHash("0x55"),
			expiration: 200,
			recipient:  common.HexToAddress("0x321"),
			amount:     big.NewInt(100),
			expected:   "2f7c8c16b1cec4063b9df791ead00e4db539072963b07bec12957c8680dc1eee6d41b2e906746b9b619af2ffb5792d794c81a90381d8b26e40b913e5a0fe0f821b",
		},
	} {
		signature, err := signer.SignWithdrawal(
			withdrawal.id,
			withdrawal.expiration,
			withdrawal.recipient,
			withdrawal.token,
			withdrawal.amount,
		)
		assert.NoError(t, err)
		assert.Equal(
			t,
			strings.ToLower(withdrawal.expected),
			strings.ToLower(hex.EncodeToString(signature)),
		)
	}
}

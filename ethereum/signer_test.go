package ethereum

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
)

func createSigner(t *testing.T) Signer {
	signer, err := NewSigner(
		"51138e68e8a5fa906d38c5b42bc01b805d7adb3fce037743fb406bb10aa83307",
		0,
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
			expected:   "2a4fead286732e459cc4f167aa34f6ca6b83fa4e7c993429582a048020a4c2840f47a9903257266605de9975d2122d1db8697b7398474889aafaa3c9d1b4bd6c01",
		},
		{
			id:         common.HexToHash("0x55"),
			expiration: 200,
			recipient:  common.HexToAddress("0x321"),
			amount:     big.NewInt(100),
			expected:   "40cd596f0d1683bbe42c6fc57220ecfda15d78d5ce9ccdf68c4da4d9a4d1a6bf63d249d59b17ead7c59b4b4da9755aad6eb9eefe656685322de074128370d31e01",
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

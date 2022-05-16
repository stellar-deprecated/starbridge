package store

import (
	"math/big"
	"time"
)

type IncomingEthereumTransaction struct {
	Hash               string
	ValueWei           *big.Int
	StellarAddress     string
	WithdrawExpiration time.Time

	TransactionBlob string
}

func (m *Memory) GetIncomingEthereumTransactionByHash(hash string) (IncomingEthereumTransaction, error) {
	return IncomingEthereumTransaction{
		Hash:     hash,
		ValueWei: big.NewInt(627836782638726),
		// Public Key	GATBFH6GV7GMWNI5RXH546BB2MDSNO3DPLGPT4EAFS5ICLRZT3D7F4YS
		// Secret Key	SBEICGMVMPF2WWIYV34IP7ON2Q6BUOT7F7IGHOTUMYUIG5K4IWIOUQC3
		StellarAddress: "GATBFH6GV7GMWNI5RXH546BB2MDSNO3DPLGPT4EAFS5ICLRZT3D7F4YS",
	}, nil
}

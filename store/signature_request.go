package store

import (
	"database/sql"
	"errors"
)

type NetworkType string

const (
	Ethereum NetworkType = "ethereum"
)

type SignatureRequest struct {
	IncomingType                    NetworkType
	IncomingEthereumTransactionHash *string
}

func (m *Memory) InsertSignatureRequestForIncomingEthereumTransaction(hash string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, sr := range m.signatureRequests {
		if sr.IncomingEthereumTransactionHash != nil &&
			*sr.IncomingEthereumTransactionHash == hash {
			return errors.New("signature request for this transaction already exists")
		}
	}

	m.signatureRequests = append(m.signatureRequests, SignatureRequest{
		IncomingType:                    Ethereum,
		IncomingEthereumTransactionHash: &hash,
	})
	return nil
}

func (m *Memory) GetSignatureRequests() ([]SignatureRequest, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.signatureRequests) == 0 {
		return nil, sql.ErrNoRows
	}

	return m.signatureRequests, nil
}

func (m *Memory) GetSignatureRequestForIncomingEthereumTransaction(hash string) (SignatureRequest, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, sr := range m.signatureRequests {
		if sr.IncomingEthereumTransactionHash != nil &&
			*sr.IncomingEthereumTransactionHash == hash {
			return sr, nil
		}
	}

	return SignatureRequest{}, sql.ErrNoRows
}

func (m *Memory) DeleteSignatureRequestForIncomingEthereumTransaction(hash string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, sr := range m.signatureRequests {
		if sr.IncomingEthereumTransactionHash != nil &&
			*sr.IncomingEthereumTransactionHash == hash {
			m.signatureRequests = append(m.signatureRequests[:i], m.signatureRequests[i+1:]...)
		}
	}

	return nil
}

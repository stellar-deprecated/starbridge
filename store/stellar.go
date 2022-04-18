package store

import "database/sql"

type OutgoingStellarTransactionState string

const (
	PendingState OutgoingStellarTransactionState = "pending"
	SuccessState OutgoingStellarTransactionState = "success"
	FailedState  OutgoingStellarTransactionState = "failed"
	ExpiredState OutgoingStellarTransactionState = "expired"
)

type OutgoingStellarTransaction struct {
	// State can only be updated by transaction streamer
	State    OutgoingStellarTransactionState
	Hash     string
	Envelope string

	IncomingType                    NetworkType
	IncomingEthereumTransactionHash *string
}

func (m *Memory) GetOutgoingStellarTransactionForEthereumByHash(hash string) (OutgoingStellarTransaction, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, tx := range m.outgoingStellarTransactions {
		if tx.IncomingType == Ethereum && *tx.IncomingEthereumTransactionHash == hash {
			return tx, nil
		}
	}

	return OutgoingStellarTransaction{}, sql.ErrNoRows
}

func (m *Memory) UpsertOutgoingStellarTransaction(newtx OutgoingStellarTransaction) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, tx := range m.outgoingStellarTransactions {
		if tx.IncomingType == newtx.IncomingType &&
			*tx.IncomingEthereumTransactionHash == *newtx.IncomingEthereumTransactionHash {
			m.outgoingStellarTransactions[i] = newtx
			return nil
		}
	}

	// If not found insert
	m.outgoingStellarTransactions = append(m.outgoingStellarTransactions, newtx)
	return nil
}

package store

import "sync"

type Memory struct {
	mutex sync.Mutex

	signatureRequests           []SignatureRequest
	outgoingStellarTransactions []OutgoingStellarTransaction
}

package main

import (
	"crypto/sha256"

	"github.com/stellar/starbridge/p2p"
)

// TODO: Clear out messages that are not needed anymore in the in-memory store

type Store struct {
	msgs map[[32]byte]p2p.MessageV0
}

func NewStore() *Store {
	return &Store{
		msgs: map[[32]byte]p2p.MessageV0{},
	}
}

func (s *Store) StoreAndUpdate(m p2p.MessageV0) (p2p.MessageV0, error) {
	hash := sha256.Sum256(m.Body)
	storedMsg, ok := s.msgs[hash]
	if ok {
		sigsSeen := map[string]bool{}
		for _, s := range m.Signatures {
			sigsSeen[string(s)] = true
		}
		for _, s := range storedMsg.Signatures {
			if sigsSeen[string(s)] {
				continue
			}
			m.Signatures = append(m.Signatures, s)
		}
	}
	s.msgs[hash] = m
	return m, nil
}

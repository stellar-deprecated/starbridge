package sigsharestellar

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/stellar/go/txnbuild"
)

type SigShareStellar struct {
	networkPassphrase string
	topic             *pubsub.Topic
}

func NewSigShareStellar(networkPassphrase string, pubSub *pubsub.PubSub) (*SigShareStellar, error) {
	topic, err := pubSub.Join("starbridge-stellar-transactions-signed")
	if err != nil {
		return nil, err
	}
	s := &SigShareStellar{
		networkPassphrase: networkPassphrase,
		topic: topic,
	}
	return s, nil
}

func (s *SigShareStellar) Share(ctx context.Context, tx *txnbuild.GenericTransaction) error {
	hash, err := tx.HashHex(s.networkPassphrase)
	if err != nil {
		return fmt.Errorf("hashing tx: %w", err)
	}

	txBytes, err := tx.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshaling tx %s: %w", hash, err)
	}

	err = s.topic.Publish(ctx, txBytes)
	if err != nil {
		return fmt.Errorf("publishing tx %s to topic: %w", hash, err)
	}
	return nil
}

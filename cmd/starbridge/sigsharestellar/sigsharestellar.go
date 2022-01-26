package sigsharestellar

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/go/txnbuild"
)

type SigShareStellarConfig struct {
	NetworkPassphrase string
	Logger            *supportlog.Entry
	PubSub            *pubsub.PubSub
}

type SigShareStellar struct {
	networkPassphrase string
	logger            *supportlog.Entry
	topic             *pubsub.Topic
}

func NewSigShareStellar(config SigShareStellarConfig) (*SigShareStellar, error) {
	topic, err := config.PubSub.Join("starbridge-stellar-transactions-signed")
	if err != nil {
		return nil, err
	}
	s := &SigShareStellar{
		networkPassphrase: config.NetworkPassphrase,
		logger:            config.Logger,
		topic:             topic,
	}
	return s, nil
}

func (s *SigShareStellar) Close() error {
	return s.topic.Close()
}

func (s *SigShareStellar) Share(ctx context.Context, tx *txnbuild.Transaction) error {
	logger := s.logger.Ctx(ctx)

	hash, err := tx.HashHex(s.networkPassphrase)
	if err != nil {
		return fmt.Errorf("hashing tx: %w", err)
	}
	logger = logger.WithField("tx", hash)
	logger = logger.WithField("sigcount", len(tx.Signatures()))

	txBytes, err := tx.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshaling tx %s: %w", hash, err)
	}

	logger = logger.WithField("topic", s.topic.String())
	err = s.topic.Publish(ctx, txBytes)
	if err != nil {
		return fmt.Errorf("publishing tx %s to topic: %w", hash, err)
	}
	logger.Infof("Tx published")
	return nil
}

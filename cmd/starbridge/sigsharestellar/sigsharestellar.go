package sigsharestellar

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"github.com/stellar/starbridge/p2p"
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
	topic, err := config.PubSub.Join("starbridge-messages-signed")
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

func (s *SigShareStellar) Share(ctx context.Context, tx *txnbuild.Transaction, sig xdr.DecoratedSignature) error {
	logger := s.logger.Ctx(ctx)

	txHash, err := tx.HashHex(s.networkPassphrase)
	if err != nil {
		return fmt.Errorf("hashing tx: %w", err)
	}
	logger = logger.WithField("txhash", txHash)

	txBytes, err := tx.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshaling: %w", err)
	}

	sigBytes, err := sig.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshaling sig: %w", err)
	}

	msg := p2p.Message{
		V: 0,
		V0: p2p.MessageV0{
			Chain:      p2p.ChainStellar,
			Body:       txBytes,
			Signatures: [][]byte{sigBytes},
		},
	}
	msgBytes, err := msg.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshaling message: %w", err)
	}

	hash := sha256.Sum256(msgBytes)
	hashHex := hex.EncodeToString(hash[:])
	logger = logger.WithField("msghash", hashHex)

	bodyHash := sha256.Sum256(msg.V0.Body)
	bodyHashHex := hex.EncodeToString(bodyHash[:])
	logger = logger.WithField("msgbodyhash", bodyHashHex)

	logger = logger.WithField("topic", s.topic.String())
	err = s.topic.Publish(ctx, msgBytes)
	if err != nil {
		return fmt.Errorf("publishing msg: %w", err)
	}
	logger.Infof("Msg published")
	return nil
}

package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/p2p"
)

type CollectorConfig struct {
	Logger    *supportlog.Entry
	PubSub    *pubsub.PubSub
	EthClient *ethclient.Client
	SecretKey *ecdsa.PrivateKey
	PublicKey *ecdsa.PublicKey
}

type Collector struct {
	logger    *supportlog.Entry
	topic     *pubsub.Topic
	ethClient *ethclient.Client
	secretKey *ecdsa.PrivateKey
	publicKey *ecdsa.PublicKey
	address   common.Address
}

func NewCollector(config CollectorConfig) (*Collector, error) {
	config.Logger.Warnf("TODO: Subscribe to topic relevant to this collector's address only.")
	topic, err := config.PubSub.Join("starbridge-messages-signed-aggregated-ethereum")
	if err != nil {
		return nil, err
	}
	c := &Collector{
		logger:    config.Logger,
		topic:     topic,
		ethClient: config.EthClient,
		secretKey: config.SecretKey,
		publicKey: &config.SecretKey.PublicKey,
		address:   crypto.PubkeyToAddress(config.SecretKey.PublicKey),
	}
	return c, nil
}

func (c *Collector) Collect() error {
	sub, err := c.topic.Subscribe()
	if err != nil {
		return err
	}
	logger := c.logger.WithField("topic", c.topic.String())
	logger.Info("Subscribed")
	ctx := context.Background()
	for {
		logger := logger

		raw, err := sub.Next(ctx)
		if err != nil {
			return err
		}

		hash := sha256.Sum256(raw.Data)
		hashHex := hex.EncodeToString(hash[:])
		logger = logger.WithField("msghash", hashHex)

		logger.Infof("Msg received")

		msg := p2p.Message{}
		err = msg.UnmarshalBinary(raw.Data)
		if err != nil {
			logger.Errorf("Unmarshaling message: %s", err)
			continue
		}

		if msg.V != 0 {
			logger.Errorf("Dropping message with unsupported version %d", msg.V)
			continue
		}

		logger = logger.WithField("msgbodysize", len(msg.V0.Body))
		logger = logger.WithField("msgsigcount", len(msg.V0.Signatures))

		bodyHash := sha256.Sum256(msg.V0.Body)
		bodyHashHex := hex.EncodeToString(bodyHash[:])
		logger = logger.WithField("msgbodyhash", bodyHashHex)

		logger.Infof("Msg unpacked")

		c.logger.Warnf("TODO: Check if message is sufficiently signed to be posted to chain.")
		c.logger.Warnf("TODO: Call contract with message.")

		// chainID, err := c.ethClient.ChainID(ctx)
		// if err != nil {
		// 	return err
		// }
		// nonce, err := c.ethClient.PendingNonceAt(ctx, c.address)
		// if err != nil {
		// 	return err
		// }
		// gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
		// if err != nil {
		// 	return err
		// }

		// transactOpts, err := bind.NewKeyedTransactorWithChainID(c.secretKey, chainID)
		// if err != nil {
		// 	return err
		// }
		// transactOpts.Nonce = big.NewInt(int64(nonce))
		// transactOpts.Value = big.NewInt(0)
		// transactOpts.GasLimit = 300000
		// transactOpts.GasPrice = gasPrice

		// contract, err := soliditybridge.
		// if err != nil {
		// 	return err
		// }
	}
}

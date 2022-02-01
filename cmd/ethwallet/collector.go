package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/p2p"
)

type CollectorConfig struct {
	Logger    *supportlog.Entry
	PubSub    *pubsub.PubSub
	EthClient *ethclient.Client
	Addr      common.Address
}

type Collector struct {
	logger    *supportlog.Entry
	topic     *pubsub.Topic
	ethClient *ethclient.Client
}

func NewCollector(config CollectorConfig) (*Collector, error) {
	config.Logger.Warnf("TODO: Subscribe to topic relevant to this collector's address only.")
	topic, err := config.PubSub.Join("starbridge-messages-signed-aggregated-ethereum")
	if err != nil {
		return nil, err
	}
	c := &Collector{
		logger:    config.Logger,
		ethClient: config.EthClient,
		topic:     topic,
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

		logger.Warnf("TODO: Send transaction to Ethereum")

		// n, err := client.BlockNumber(ctx)
		// if err != nil {
		// 	return err
		// }
		// spew.Dump("n", n)

		// sk, err := crypto.HexToECDSA("da1da8a2bb731e77b295acbf3f4a9e4a9eae9ea0735e8ca14334f1e31ad22ab8")
		// if err != nil {
		// 	return err
		// }
		// pk := sk.Public()
		// pkECDSA, ok := pk.(*ecdsa.PublicKey)
		// if !ok {
		// 	return err
		// }
		// addr := crypto.PubkeyToAddress(*pkECDSA)
		// nonce, err := client.PendingNonceAt(ctx, addr)
		// if err != nil {
		// 	return err
		// }
		// spew.Dump("nonce", nonce)
		// gasPrice, err := client.SuggestGasPrice(ctx)
		// if err != nil {
		// 	return err
		// }
		// spew.Dump("gasPrice", gasPrice)
		// tx := types.NewTx(types.DynamicFeeTx{
		// 	ChainID: big.NewInt(1337),
		// 	Nonce:   nonce,
		// })
		// client.SendTransaction(ctx, tx)

	}
}

package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/stellar/go/clients/horizonclient"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"github.com/stellar/starbridge/p2p"
)

type CollectorConfig struct {
	NetworkPassphrase string
	Logger            *supportlog.Entry
	HorizonClient     horizonclient.ClientInterface
	PubSub            *pubsub.PubSub
}

type Collector struct {
	networkPassphrase string
	logger            *supportlog.Entry
	horizonClient     horizonclient.ClientInterface
	topic             *pubsub.Topic
}

func NewCollector(config CollectorConfig) (*Collector, error) {
	topic, err := config.PubSub.Join("starbridge-stellar-transactions-signed-aggregated")
	if err != nil {
		return nil, err
	}
	c := &Collector{
		networkPassphrase: config.NetworkPassphrase,
		logger:            config.Logger,
		horizonClient:     config.HorizonClient,
		topic:             topic,
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

		bodyHash := sha256.Sum256(raw.Data)
		bodyHashHex := hex.EncodeToString(bodyHash[:])
		logger = logger.WithField("msgbodyhash", bodyHashHex)

		logger.Infof("Msg unpacked")

		txB64 := base64.StdEncoding.EncodeToString(msg.V0.Body)
		txGeneric, err := txnbuild.TransactionFromXDR(txB64)
		if err != nil {
			return fmt.Errorf("unmarshaling tx %s", txB64)
		}
		tx, ok := txGeneric.Transaction()
		if !ok {
			return fmt.Errorf("unsupported tx type")
		}

		txHash, err := tx.Hash(c.networkPassphrase)
		if err != nil {
			return err
		}
		logger = logger.WithField("txhash", hex.EncodeToString(txHash[:]))

		logger.Infof("Tx unpacked")

		sigs := make([]xdr.DecoratedSignature, len(msg.V0.Signatures))
		for i, sigBytes := range msg.V0.Signatures {
			err = sigs[i].UnmarshalBinary(sigBytes)
			if err != nil {
				return err
			}
		}
		tx, err = tx.ClearSignatures()
		if err != nil {
			return err
		}
		tx, err = tx.AddSignatureDecorated(sigs...)
		if err != nil {
			return err
		}

		logger = logger.WithField("txsigcount", len(tx.Signatures()))
		logger.Infof("Tx updated with all known sigs")

		tx, err = AuthorizedTransaction(c.horizonClient, txHash, tx)
		if errors.Is(err, ErrNotAuthorized) {
			logger.Infof("Tx not yet authorized")
			continue
		} else if err != nil {
			return err
		}
		logger = logger.WithField("txsigcount", len(tx.Signatures()))
		logger.Infof("Tx authorized")

		// Submit transaction.
		go func() {
			logger.Infof("Tx submitting")
			txResp, err := c.horizonClient.SubmitTransaction(tx)
			if err != nil {
				logger.Error(err)
				return
			}
			logger.WithField("txsuccessful", txResp.Successful).Infof("Tx submitted")
		}()
	}
}

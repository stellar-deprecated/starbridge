package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/go/txnbuild"
)

type CollectorConfig struct {
	NetworkPassphrase string
	Logger            *supportlog.Entry
	HorizonClient     horizonclient.ClientInterface
	PubSub            *pubsub.PubSub
}

type Collector struct {
	address           *keypair.FromAddress
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
	c.logger.Infof("Subscribed to topic %s", c.topic.String())
	ctx := context.Background()
	for {
		logger := c.logger

		msg, err := sub.Next(ctx)
		if err != nil {
			return err
		}

		txBytes := msg.Data
		txB64 := base64.StdEncoding.EncodeToString(txBytes)
		txGeneric, err := txnbuild.TransactionFromXDR(txB64)
		if err != nil {
			return err
		}
		// TODO: Support all tx types.
		tx, ok := txGeneric.Transaction()
		if !ok {
			return fmt.Errorf("unsupported tx type")
		}

		hash, err := tx.Hash(c.networkPassphrase)
		if err != nil {
			return err
		}
		logger = logger.WithField("tx", hex.EncodeToString(hash[:]))

		tx, err = AuthorizedTransaction(c.horizonClient, hash, tx)
		if errors.Is(err, ErrNotAuthorized) {
			logger.Infof("tx not yet authorized")
			continue
		} else if err != nil {
			return err
		}
		logger.Infof("tx authorized: sig count: %d", len(tx.Signatures()))

		// Submit transaction.
		go func() {
			// TODO: Wrap in a fee bump transaction.
			logger.Infof("tx submitting: sig count: %d", len(tx.Signatures()))
			txResp, err := c.horizonClient.SubmitTransaction(tx)
			if err != nil {
				logger.Error(err)
				return
			}
			logger.Infof("tx submitted: successful: %t", txResp.Successful)
		}()
	}
}

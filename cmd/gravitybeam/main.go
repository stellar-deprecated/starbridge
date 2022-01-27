package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	ff "github.com/peterbourgon/ff/v3"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/p2p"
)

func main() {
	logger := supportlog.New()
	logger.SetLevel(logrus.InfoLevel)
	err := run(os.Args[1:], logger)
	if err != nil {
		logger.WithStack(err).Error(err)
		os.Exit(1)
	}
}

func run(args []string, logger *supportlog.Entry) error {
	ctx := context.Background()

	fs := flag.NewFlagSet("gravitybeam", flag.ExitOnError)

	portP2P := "0"
	peers := ""
	horizonURL := "https://horizon-testnet.stellar.org"

	fs.StringVar(&portP2P, "port-p2p", portP2P, "Port to accept P2P requests on (also via PORT_P2P)")
	fs.StringVar(&peers, "peers", peers, "Comma-separated list of addresses of peers to connect to on start (also via PEERS)")
	fs.StringVar(&horizonURL, "horizon", horizonURL, "Horizon URL (also via HORIZON_URL)")

	err := ff.Parse(fs, args, ff.WithEnvVarNoPrefix())
	if err != nil {
		return err
	}

	logger.Info("Starting")

	horizonClient := &horizonclient.Client{HorizonURL: horizonURL}

	networkDetails, err := horizonClient.Root()
	if err != nil {
		return err
	}

	pubSub, err := p2p.New(ctx, p2p.Config{
		Logger: logger,
		Port:   portP2P,
		Peers:  strings.Split(peers, ","),
	})
	if err != nil {
		return err
	}

	store := NewStore()
	collector, err := NewCollector(CollectorConfig{
		NetworkPassphrase: networkDetails.NetworkPassphrase,
		Logger:            logger,
		HorizonClient:     horizonClient,
		PubSub:            pubSub,
		Store:             store,
	})
	if err != nil {
		return fmt.Errorf("creating collector: %v", err)
	}
	err = collector.Collect()
	if err != nil {
		return fmt.Errorf("starting collecting: %v", err)
	}

	return nil
}

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ff "github.com/peterbourgon/ff/v3"
	"github.com/sirupsen/logrus"
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

	fs := flag.NewFlagSet("ethwallet", flag.ExitOnError)

	portP2P := "0"
	peers := ""
	rpc := "http://localhost:8545"
	addrStr := ""

	fs.StringVar(&portP2P, "port-p2p", portP2P, "Port to accept P2P requests on (also via PORT_P2P)")
	fs.StringVar(&peers, "peers", peers, "Comma-separated list of addresses of peers to connect to on start (also via PEERS)")
	fs.StringVar(&rpc, "rpc", rpc, "Ethereum client RPC (also via RPC)")
	fs.StringVar(&addrStr, "addr", addrStr, "Ethereum address this wallet is listening to transactions about (also via ADDR)")

	err := ff.Parse(fs, args, ff.WithEnvVarNoPrefix())
	if err != nil {
		return err
	}

	logger.Info("Starting")

	pubSub, err := p2p.New(ctx, p2p.Config{
		Logger: logger,
		Port:   portP2P,
		Peers:  strings.Split(peers, ","),
	})
	if err != nil {
		return err
	}

	addr := common.HexToAddress(addrStr)

	client, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}

	collector, err := NewCollector(CollectorConfig{
		Logger:    logger,
		PubSub:    pubSub,
		EthClient: client,
		Addr:      addr,
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

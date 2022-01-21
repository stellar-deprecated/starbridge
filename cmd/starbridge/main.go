package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	ff "github.com/peterbourgon/ff/v3"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/cmd/starbridge/sigsharestellar"
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
	fs := flag.NewFlagSet("starbridge", flag.ExitOnError)

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

	logger.Info("Starting...")

	horizonClient := &horizonclient.Client{HorizonURL: horizonURL}

	networkDetails, err := horizonClient.Root()
	if err != nil {
		return err
	}

	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/" + portP2P))
	if err != nil {
		return err
	}
	hostAddrInfo := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	hostAddrs, err := peer.AddrInfoToP2pAddrs(&hostAddrInfo)
	if err != nil {
		return err
	}
	for _, a := range hostAddrs {
		logger.Infof("Listening for p2p on... %v", a)
	}

	if peers == "" {
		logger.Info("Using mdns to discover local peers...")
		mdnsService := mdns.NewMdnsService(host, "gravitybeam", &mdnsNotifee{Host: host, Logger: logger})
		err = mdnsService.Start()
		if err != nil {
			return err
		}
	} else {
		peersArr := strings.Split(peers, ",")
		for _, p := range peersArr {
			p := p
			go func() {
				logger := logger.WithField("peer", p)
				logger.Info("Connecting to peer...")
				peerAddrInfo, err := peer.AddrInfoFromString(p)
				if err != nil {
					logger.Errorf("Error parsing peer address: %v", err)
					return
				}
				err = host.Connect(context.Background(), *peerAddrInfo)
				if err != nil {
					logger.Errorf("Error connecting to peer: %v", err)
					return
				}
				logger.Info("Connected to peer")
			}()
		}
	}

	pubSub, err := pubsub.NewGossipSub(context.Background(), host)
	if err != nil {
		return fmt.Errorf("configuring pubsub: %v", err)
	}

	_, err = sigsharestellar.NewSigShareStellar(networkDetails.NetworkPassphrase, pubSub)
	if err != nil {
		return fmt.Errorf("setting up sharing stellar signatures: %v", err)
	}
	// TODO: When generating and signing a Stellar transaction, call sigShareStellar.Share(ctx, tx).

	return nil
}

type mdnsNotifee struct {
	Host   host.Host
	Logger *supportlog.Entry
}

func (n *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if pi.ID == n.Host.ID() {
		return
	}
	n.Logger.Infof("Connecting to peer discovered via mdns: %s", pi.ID.Pretty())
	err := n.Host.Connect(context.Background(), pi)
	if err != nil {
		n.Logger.WithStack(err).Error(fmt.Errorf("Error connecting to peer %s: %w", pi.ID.Pretty(), err))
	}
}

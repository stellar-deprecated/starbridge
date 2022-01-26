package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	libp2pnetwork "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	ff "github.com/peterbourgon/ff/v3"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/go/txnbuild"

	"github.com/stellar/starbridge/cmd/starbridge/integrations"
	"github.com/stellar/starbridge/cmd/starbridge/sigsharestellar"
	"github.com/stellar/starbridge/cmd/starbridge/transform"
)

func main() {
	logger := supportlog.New()
	logger.SetLevel(logrus.InfoLevel)
	integrations.SetLogger(logger)
	err := run(os.Args[1:], logger)
	if err != nil {
		logger.WithStack(err).Error(err)
		os.Exit(1)
	}
}

func run(args []string, logger *supportlog.Entry) error {
	fs := flag.NewFlagSet("starbridge", flag.ExitOnError)

	txHash := ""
	seed := ""
	portP2P := "0"
	peers := ""
	horizonURL := "https://horizon-testnet.stellar.org"

	fs.StringVar(&txHash, "txHash", "", "txHash on Ethereum to be queried for conversion")
	fs.StringVar(&seed, "seed", "", "Seed secret key for Stellar with which to sign transactions for this node")
	fs.StringVar(&portP2P, "port-p2p", portP2P, "Port to accept P2P requests on (also via PORT_P2P)")
	fs.StringVar(&peers, "peers", peers, "Comma-separated list of addresses of peers to connect to on start (also via PEERS)")
	fs.StringVar(&horizonURL, "horizon", horizonURL, "Horizon URL (also via HORIZON_URL)")

	err := ff.Parse(fs, args, ff.WithEnvVarNoPrefix())
	if err != nil {
		return err
	}

	if txHash == "" {
		return fmt.Errorf("needs a valid 'txHash' command line option")
	}
	if seed == "" {
		return fmt.Errorf("needs a valid 'seed' command line option")
	}

	logger.Info("Starting")

	horizonClient := &horizonclient.Client{HorizonURL: horizonURL}

	networkDetails, err := horizonClient.Root()
	if err != nil {
		return err
	}

	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/" + portP2P))
	if err != nil {
		return err
	}
	host.Network().Notify(&libp2pnetwork.NotifyBundle{
		ConnectedF: func(n libp2pnetwork.Network, c libp2pnetwork.Conn) {
			logger := logger.WithField("peer", c.RemotePeer().Pretty())
			logger.Info("Connected to peer")
		},
	})
	hostAddrInfo := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	hostAddrs, err := peer.AddrInfoToP2pAddrs(&hostAddrInfo)
	if err != nil {
		return err
	}
	for _, a := range hostAddrs {
		logger.WithField("addr", a).Info("Listening")
	}

	if peers != "" {
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
			}()
		}
	}

	mdnsService := mdns.NewMdnsService(host, "starbridge", &mdnsNotifee{Host: host, Logger: logger})
	err = mdnsService.Start()
	if err != nil {
		return err
	}

	pubSub, err := pubsub.NewGossipSub(context.Background(), host)
	if err != nil {
		return fmt.Errorf("configuring pubsub: %v", err)
	}

	sigShareStellar, err := sigsharestellar.NewSigShareStellar(sigsharestellar.SigShareStellarConfig{
		NetworkPassphrase: networkDetails.NetworkPassphrase,
		Logger:            logger,
		PubSub:            pubSub,
	})
	if err != nil {
		return fmt.Errorf("setting up sharing stellar signatures: %v", err)
	}
	// TODO: When generating and signing a Stellar transaction, call sigShareStellar.Share(ctx, tx).

	time.Sleep(2 * time.Second)

	modelTxEth, err := integrations.FetchEthTxByHash(txHash)
	if err != nil {
		return fmt.Errorf("fetching eth tx %s: %w", txHash, err)
	}
	logger.Infof("transaction fetched as modelTxEth: %s", modelTxEth)

	modelTxStellar, err := transform.MapTxToChain(modelTxEth)
	if err != nil {
		return fmt.Errorf("mapping model eth tx to model stellar tx: %w", err)
	}
	logger.Infof("transaction converted to modelTxStellar: %s", modelTxStellar)
	if modelTxStellar.To != modelTxStellar.Data.TargetDestinationAddressOnRemoteChain {
		return fmt.Errorf("incorrect mapping since To value of converted transaction should match TargetDestinationAddressOnRemoteChain from event data")
	}

	stellarTx, err := integrations.Transaction2Stellar(modelTxStellar)
	if err != nil {
		return fmt.Errorf("building stellar tx: %w", err)
	}
	logger.Infof("transaction as an unsigned stellarTx: %s", integrations.Stellar2String(stellarTx))

	logger.Infof("signing Stellar tx...")
	signedStellarTx, err := signTxForStellar(stellarTx, seed)
	if err != nil {
		return fmt.Errorf("signing tx: %w", err)
	}
	logger.Infof("transaction as a signed stellarTx: %s", integrations.Stellar2String(signedStellarTx))

	signedStellarTxBase64String, err := signedStellarTx.Base64()
	if err != nil {
		return fmt.Errorf("converting to base64 string: %w", err)
	}
	logger.Infof("stellar tx base64 encoded: %s", signedStellarTxBase64String)

	err = sigShareStellar.Share(context.Background(), signedStellarTx)
	if err != nil {
		return fmt.Errorf("sharing stellar tx: %w", err)
	}

	time.Sleep(2 * time.Second)

	return nil
}

func signTxForStellar(tx *txnbuild.Transaction, seed string) (*txnbuild.Transaction, error) {
	networkPassphrase := network.TestNetworkPassphrase

	kp, err := keypair.Parse(seed)
	if err != nil {
		return nil, fmt.Errorf("cannot parse seed into keypair: %s", err)
	}

	// keep adding signatures
	signedTx, err := tx.Sign(networkPassphrase, kp.(*keypair.Full))
	if err != nil {
		return nil, fmt.Errorf("cannot sign tx with keypair (pubKey: %s): %s", kp.Address(), err)
	}

	return signedTx, nil
}

type mdnsNotifee struct {
	Host   host.Host
	Logger *supportlog.Entry
}

func (n *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if pi.ID == n.Host.ID() {
		return
	}
	err := n.Host.Connect(context.Background(), pi)
	if err != nil {
		n.Logger.WithStack(err).Error(fmt.Errorf("Error connecting to peer %s: %w", pi.ID.Pretty(), err))
	}
}

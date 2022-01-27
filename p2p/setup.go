package p2p

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	supportlog "github.com/stellar/go/support/log"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Logger *supportlog.Entry
	Port   string
	Peers  []string
}

func New(ctx context.Context, c Config) (*pubsub.PubSub, error) {
	logger := c.Logger

	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/" + c.Port))
	if err != nil {
		return nil, err
	}
	host.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(n network.Network, c network.Conn) {
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
		return nil, err
	}
	for _, a := range hostAddrs {
		logger.WithField("addr", a).Info("Listening")
	}

	connected := uint64(0)
	g := errgroup.Group{}
	for _, p := range c.Peers {
		p := strings.TrimSpace(p)
		if p == "" {
			continue
		}
		g.Go(func() error {
			logger := logger.WithField("peer", p)
			logger.Info("Connecting to peer...")
			peerAddrInfo, err := peer.AddrInfoFromString(p)
			if err != nil {
				logger.Errorf("Error parsing peer address: %v", err)
				return nil
			}
			err = host.Connect(ctx, *peerAddrInfo)
			if err != nil {
				logger.Errorf("Error connecting to peer: %v", err)
				return nil
			}
			atomic.AddUint64(&connected, 1)
			return nil
		})
	}
	_ = g.Wait()
	logger.Errorf("Connected to %d peers", connected)

	mdnsService := mdns.NewMdnsService(host, "starbridge", &mdnsNotifee{Host: host, Logger: logger})
	err = mdnsService.Start()
	if err != nil {
		return nil, err
	}

	pubSub, err := pubsub.NewGossipSub(context.Background(), host)
	if err != nil {
		return nil, fmt.Errorf("configuring pubsub: %v", err)
	}

	return pubSub, nil
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

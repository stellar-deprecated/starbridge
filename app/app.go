package app

import (
	"context"
	"encoding/hex"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/clients/stellarcore"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/support/log"

	"github.com/stellar/starbridge/controllers"
	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/httpx"
	"github.com/stellar/starbridge/stellar"
)

type App struct {
	appCtx    context.Context
	cancelCtx context.CancelFunc

	httpServer *httpx.Server

	prometheusRegistry *prometheus.Registry
}

type Config struct {
	Port      uint16 `toml:"port" valid:"-"`
	AdminPort uint16 `toml:"admin_port" valid:"-"`

	PostgresDSN string `toml:"postgres_dsn" valid:"-"`

	HorizonURL string `toml:"horizon_url" valid:"-"`
	CoreURL    string `toml:"core_url" valid:"-"`

	NetworkPassphrase       string `toml:"network_passphrase" valid:"-"`
	StellarBridgeAccount    string `toml:"stellar_bridge_account" valid:"stellar_accountid"`
	StellarBridgeContractID string `toml:"stellar_bridge_contract_id" valid:"-"`
	StellarPrivateKey       string `toml:"stellar_private_key" valid:"stellar_seed"`

	EthereumRPCURL        string `toml:"ethereum_rpc_url" valid:"-"`
	EthereumBridgeAddress string `toml:"ethereum_bridge_address" valid:"-"`
	EthereumPrivateKey    string `toml:"ethereum_private_key" valid:"-"`

	AssetMapping []controllers.AssetMappingConfigEntry `toml:"asset_mapping" valid:"-"`

	EthereumFinalityBuffer uint64        `toml:"-" valid:"-"`
	WithdrawalWindow       time.Duration `toml:"-" valid:"-"`
}

func NewApp(config Config) *App {
	app := &App{
		prometheusRegistry: prometheus.NewRegistry(),
	}
	app.initGracefulShutdown()

	app.initHTTP(config)
	app.initLogger()
	app.initPrometheus()

	return app
}

func (a *App) initGracefulShutdown() {
	a.appCtx, a.cancelCtx = context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-signalChan:
			log.Info("Shutdown signal received...")
			a.Close()
		case <-a.appCtx.Done():
			return
		}
	}()
}

func (a *App) initPrometheus() {
	a.httpServer.RegisterMetrics(a.prometheusRegistry)
}

func (a *App) initLogger() {
	log.SetLevel(log.InfoLevel)
}

func (a *App) initHTTP(config Config) {
	var (
		signerKey *keypair.Full
		err       error
	)
	if config.StellarPrivateKey != "" {
		signerKey, err = keypair.ParseFull(config.StellarPrivateKey)
		if err != nil {
			log.Fatalf("cannot parse signer secret key: %v", err)
		}
	}

	contractIDBytes, err := hex.DecodeString(config.StellarBridgeContractID)
	if err != nil {
		log.Fatalf("cannot parse bridge contract id: %v", err)
	}
	if len(contractIDBytes) != 32 {
		log.Fatalf("invalid contract id: %v", config.StellarBridgeContractID)
	}
	var bridgeContractID [32]byte
	copy(bridgeContractID[:], contractIDBytes)

	stellarObserver := stellar.NewObserver(
		bridgeContractID,
		&horizonclient.Client{
			HorizonURL: config.HorizonURL,
			// TODO set proper timeouts
			HTTP: http.DefaultClient,
		},
		&stellarcore.Client{URL: config.CoreURL, HTTP: http.DefaultClient},
	)

	ethRPCClient, err := ethclient.Dial(config.EthereumRPCURL)
	if err != nil {
		log.WithField("err", err).Fatal("could not dial ethereum node")
	}
	ethObserver, err := ethereum.NewObserver(ethRPCClient, config.EthereumBridgeAddress)
	if err != nil {
		log.WithField("err", err).Fatal("could not create ethereum observer")
	}

	converter, err := controllers.NewAssetConverter(config.NetworkPassphrase, config.StellarBridgeAccount, config.AssetMapping)
	if err != nil {
		log.Fatalf("unable to create asset converter: %v", err)
	}

	domainSeparator, err := ethObserver.GetDomainSeparator(a.appCtx)
	if err != nil {
		log.Fatalf("unable to fetch domain separator: %v", err)
	}

	ethSigner, err := ethereum.NewSigner(config.EthereumPrivateKey, domainSeparator)
	if err != nil {
		log.Fatalf("cannot create ethereum signer: %v", err)
	}

	stellarBuilder := &stellar.Builder{
		BridgeAccount:    config.StellarBridgeAccount,
		BridgeContractID: bridgeContractID,
	}

	stellarSigner := &stellar.Signer{
		NetworkPassphrase: config.NetworkPassphrase,
		Signer:            signerKey,
	}

	httpServer, err := httpx.NewServer(httpx.ServerConfig{
		Ctx:                a.appCtx,
		Port:               config.Port,
		AdminPort:          config.AdminPort,
		PrometheusRegistry: a.prometheusRegistry,
		StellarWithdrawalHandler: &controllers.StellarWithdrawalHandler{
			StellarBuilder:         stellarBuilder,
			StellarSigner:          stellarSigner,
			StellarObserver:        stellar.Observer{},
			WithdrawalWindow:       config.WithdrawalWindow,
			Converter:              converter,
			EthereumObserver:       ethObserver,
			EthereumFinalityBuffer: config.EthereumFinalityBuffer,
		},
		EthereumWithdrawalHandler: &controllers.EthereumWithdrawalHandler{
			EthereumObserver:       ethObserver,
			StellarObserver:        stellarObserver,
			EthereumSigner:         ethSigner,
			EthereumFinalityBuffer: config.EthereumFinalityBuffer,
			WithdrawalWindow:       config.WithdrawalWindow,
			Converter:              converter,
		},
		EthereumRefundHandler: &controllers.EthereumRefundHandler{
			EthereumObserver:       ethObserver,
			StellarObserver:        stellarObserver,
			EthereumSigner:         ethSigner,
			EthereumFinalityBuffer: config.EthereumFinalityBuffer,
			WithdrawalWindow:       config.WithdrawalWindow,
		},
		StellarRefundHandler: &controllers.StellarRefundHandler{
			StellarBuilder:         stellarBuilder,
			StellarSigner:          stellarSigner,
			EthereumObserver:       ethObserver,
			StellarObserver:        stellarObserver,
			WithdrawalWindow:       config.WithdrawalWindow,
			EthereumFinalityBuffer: config.EthereumFinalityBuffer,
		},
	})
	if err != nil {
		log.Fatal("unable to create http server", err)
	}
	a.httpServer = httpServer
}

// Run starts all services and block until they are gracefully shut down.
func (a *App) Run() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		a.RunHTTPServer()
		wg.Done()
	}()

	wg.Wait()
	log.Info("Bye")
}

// RunHTTPServer starts http server
func (a *App) RunHTTPServer() {
	err := a.httpServer.Serve()
	if err != nil {
		log.WithField("error", err).Error("error running http server")
	}
}

func (a *App) Close() {
	a.cancelCtx()
}

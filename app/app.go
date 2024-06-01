package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/support/db"
	"github.com/stellar/go/support/log"

	"github.com/stellar/starbridge/backend"
	"github.com/stellar/starbridge/controllers"
	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/httpx"
	"github.com/stellar/starbridge/stellar/signer"
	"github.com/stellar/starbridge/stellar/txbuilder"
	"github.com/stellar/starbridge/stellar/txobserver"
	"github.com/stellar/starbridge/store"
)

type App struct {
	appCtx    context.Context
	cancelCtx context.CancelFunc

	httpServer      *httpx.Server
	worker          *backend.Worker
	session         *db.Session
	stellarObserver *txobserver.Observer

	prometheusRegistry *prometheus.Registry
}

type Config struct {
	Port      uint16 `toml:"port" valid:"-"`
	AdminPort uint16 `toml:"admin_port" valid:"-"`

	PostgresDSN string `toml:"postgres_dsn" valid:"-"`

	HorizonURL           string `toml:"horizon_url" valid:"-"`
	NetworkPassphrase    string `toml:"network_passphrase" valid:"-"`
	StellarBridgeAccount string `toml:"stellar_bridge_account" valid:"stellar_accountid"`
	StellarPrivateKey    string `toml:"stellar_private_key" valid:"stellar_seed"`

	EthereumRPCURL        string `toml:"ethereum_rpc_url" valid:"-"`
	EthereumBridgeAddress string `toml:"ethereum_bridge_address" valid:"-"`
	EthereumPrivateKey    string `toml:"ethereum_private_key" valid:"-"`

	AssetMapping []backend.AssetMappingConfigEntry `toml:"asset_mapping" valid:"-"`

	EthereumFinalityBuffer uint64        `toml:"-" valid:"-"`
	WithdrawalWindow       time.Duration `toml:"-" valid:"-"`
}

func NewApp(config Config) *App {
	app := &App{
		prometheusRegistry: prometheus.NewRegistry(),
	}

	client := &horizonclient.Client{
		HorizonURL: config.HorizonURL,
		// TODO set proper timeouts
		HTTP: http.DefaultClient,
	}

	app.initDB(config)
	app.initGracefulShutdown()
	app.stellarObserver = txobserver.NewObserver(
		config.StellarBridgeAccount,
		client,
		app.NewStore(),
	)
	ethRPCClient, err := ethclient.Dial(config.EthereumRPCURL)
	if err != nil {
		log.WithField("err", err).Fatal("could not dial ethereum node")
	}
	ethObserver, err := ethereum.NewObserver(ethRPCClient, config.EthereumBridgeAddress)
	if err != nil {
		log.WithField("err", err).Fatal("could not create ethereum observer")
	}
	app.initHTTP(config, client, ethObserver)
	app.initWorker(config, client, ethObserver)
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

func (a *App) initDB(config Config) {
	session, err := db.Open("postgres", config.PostgresDSN)
	if err != nil {
		log.Fatalf("cannot open DB: %v", err)
	}

	a.session = session
	err = store.InitSchema(session.DB.DB)
	if err != nil {
		log.Fatalf("cannot init DB: %v", err)
	}
}

func (a *App) initWorker(config Config, client *horizonclient.Client, ethObserver ethereum.Observer) {
	var (
		signerKey *keypair.Full
		err       error
	)
	if config.StellarPrivateKey != "" {
		signerKey, err = keypair.ParseFull(config.StellarPrivateKey)
		if err != nil {
			log.Fatalf("cannot pase signer secret key: %v", err)
		}
	}

	converter, err := backend.NewAssetConverter(config.AssetMapping)
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

	a.worker = &backend.Worker{
		Store:         a.NewStore(),
		StellarClient: client,
		StellarBuilder: &txbuilder.Builder{
			BridgeAccount: config.StellarBridgeAccount,
		},
		StellarSigner: &signer.Signer{
			NetworkPassphrase: config.NetworkPassphrase,
			Signer:            signerKey,
		},
		StellarObserver: a.stellarObserver,
		EthereumSigner:  ethSigner,
		StellarWithdrawalValidator: backend.StellarWithdrawalValidator{
			Session:          a.session.Clone(),
			WithdrawalWindow: config.WithdrawalWindow,
			Converter:        converter,
		},
		StellarRefundValidator: backend.StellarRefundValidator{
			Session:                a.session.Clone(),
			WithdrawalWindow:       config.WithdrawalWindow,
			Observer:               ethObserver,
			EthereumFinalityBuffer: config.EthereumFinalityBuffer,
		},
		EthereumWithdrawalValidator: backend.EthereumWithdrawalValidator{
			Observer:               ethObserver,
			EthereumFinalityBuffer: config.EthereumFinalityBuffer,
			WithdrawalWindow:       config.WithdrawalWindow,
			Converter:              converter,
		},
		EthereumRefundValidator: backend.EthereumRefundValidator{
			Session:          a.session.Clone(),
			WithdrawalWindow: config.WithdrawalWindow,
		},
	}
}

func (a *App) initHTTP(config Config, client *horizonclient.Client, ethObserver ethereum.Observer) {
	converter, err := backend.NewAssetConverter(config.AssetMapping)
	if err != nil {
		log.Fatal("unable to create asset converter", err)
	}

	httpServer, err := httpx.NewServer(httpx.ServerConfig{
		Ctx:                a.appCtx,
		Port:               config.Port,
		AdminPort:          config.AdminPort,
		PrometheusRegistry: a.prometheusRegistry,
		StellarWithdrawalHandler: &controllers.StellarWithdrawalHandler{
			StellarClient: client,
			Observer:      ethObserver,
			Store:         a.NewStore(),
			StellarWithdrawalValidator: backend.StellarWithdrawalValidator{
				Session:          a.session.Clone(),
				WithdrawalWindow: config.WithdrawalWindow,
				Converter:        converter,
			},
			EthereumFinalityBuffer: config.EthereumFinalityBuffer,
		},
		EthereumWithdrawalHandler: &controllers.EthereumWithdrawalHandler{
			Store: a.NewStore(),
			EthereumWithdrawalValidator: backend.EthereumWithdrawalValidator{
				Observer:               ethObserver,
				EthereumFinalityBuffer: config.EthereumFinalityBuffer,
				WithdrawalWindow:       config.WithdrawalWindow,
				Converter:              converter,
			},
		},
		EthereumRefundHandler: &controllers.EthereumRefundHandler{
			Observer: ethObserver,
			Store:    a.NewStore(),
			EthereumRefundValidator: backend.EthereumRefundValidator{
				Session:          a.session.Clone(),
				WithdrawalWindow: config.WithdrawalWindow,
			},
			EthereumFinalityBuffer: config.EthereumFinalityBuffer,
		},
		StellarRefundHandler: &controllers.StellarRefundHandler{
			StellarClient: client,
			Store:         a.NewStore(),
			StellarRefundValidator: backend.StellarRefundValidator{
				Session:                a.session.Clone(),
				WithdrawalWindow:       config.WithdrawalWindow,
				Observer:               ethObserver,
				EthereumFinalityBuffer: config.EthereumFinalityBuffer,
			},
		},
		TestDepositHandler: &controllers.TestDeposit{
			Store: a.NewStore(),
			// This will crash if no asset mappings - probably fine for a demo
			// because it requires at least one mapping.
			Token: config.AssetMapping[0].EthereumToken,
		},
	})
	if err != nil {
		log.Fatal("unable to create http server", err)
	}
	a.httpServer = httpServer
}

// NewStore returns a new instance of store.DB
func (a *App) NewStore() *store.DB {
	return &store.DB{Session: a.session.Clone()}
}

// Run starts all services and block until they are gracefully shut down.
func (a *App) Run() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		a.RunHTTPServer()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a.RunBackendWorker()
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

// RunBackendWorker starts backend worker responsible for building and signing
// transactions
func (a *App) RunBackendWorker() {
	a.worker.Run(a.appCtx)
}

func (a *App) Close() {
	a.cancelCtx()
}

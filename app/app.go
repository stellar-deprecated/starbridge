package app

import (
	"context"
	"github.com/stellar/starbridge/concordium"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	concordiumproto "github.com/Concordium/concordium-go-sdk/grpc-api"

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

type perRPCCredentials string

func (c perRPCCredentials) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{
		"authentication": string(c),
	}, nil
}

func (c perRPCCredentials) RequireTransportSecurity() bool {
	return false
}

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

	EthereumRPCURL              string `toml:"ethereum_rpc_url" valid:"-"`
	EthereumBridgeAddress       string `toml:"ethereum_bridge_address" valid:"-"`
	EthereumBridgeConfigVersion uint32 `toml:"ethereum_bridge_config_version" valid:"-"`
	EthereumPrivateKey          string `toml:"ethereum_private_key" valid:"-"`

	OkxRPCURL              string `toml:"okx_rpc_url" valid:"-"`
	OkxBridgeAddress       string `toml:"okx_bridge_address" valid:"-"`
	OkxBridgeConfigVersion uint32 `toml:"okx_bridge_config_version" valid:"-"`
	OkxPrivateKey          string `toml:"okx_private_key" valid:"-"`

	ConcordiumNodeService         string `toml:"concordium_node_service" valid:"-"`
	ConcordiumGRPCURL             string `toml:"concordium_grpc_url" valid:"-"`
	ConcordiumAuthToken           string `toml:"concordium_auth_token" valid:"-"`
	ConcordiumBridgeAddress       string `toml:"concordium_bridge_address" valid:"-"`
	ConcordiumPrivateKey          string `toml:"concordium_private_key" valid:"-"`
	ConcordiumBridgeConfigVersion uint32 `toml:"concordium_bridge_config_version" valid:"-"`

	AssetMapping []backend.AssetMappingConfigEntry `toml:"asset_mapping" valid:"-"`

	EthereumFinalityBuffer uint64        `toml:"ethereum_finality_buffer" valid:"-"`
	OkxFinalityBuffer      uint64        `toml:"okx_finality_buffer" valid:"-"`
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
	okxRPCClient, err := ethclient.Dial(config.OkxRPCURL)
	if err != nil {
		log.WithField("err", err).Fatal("could not dial ethereum node")
	}
	okxObserver, err := ethereum.NewObserver(okxRPCClient, config.OkxBridgeAddress)
	if err != nil {
		log.WithField("err", err).Fatal("could not create ethereum observer")
	}
	ccdGRPCClientInterface, err := grpc.Dial(config.ConcordiumGRPCURL, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(perRPCCredentials(config.ConcordiumAuthToken)),
	}...)
	if err != nil {
		log.WithField("err", err).Fatal("could not dial concordium node")
	}
	ccdGRPCClient := concordiumproto.NewP2PClient(ccdGRPCClientInterface)
	ccdObserver, err := concordium.NewObserver(ccdGRPCClient, config.ConcordiumBridgeAddress, config.ConcordiumNodeService)
	if err != nil {
		log.WithField("err", err).Fatal("could not create concordium observer")
	}
	if err != nil {
		log.WithField("err", err).Fatal("could not create concordium observer")
	}
	app.initHTTP(config, client, ethObserver, ccdObserver, okxObserver)
	app.initWorker(config, client, ethObserver, ccdObserver, okxObserver)
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

func (a *App) initWorker(config Config, client *horizonclient.Client, ethObserver ethereum.Observer, ccdObserver concordium.Observer, okxObserver ethereum.Observer) {
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
		log.Fatal("unable to create asset converter", err)
	}

	ethSigner, err := ethereum.NewSigner(config.EthereumPrivateKey, config.EthereumBridgeConfigVersion)
	if err != nil {
		log.Fatalf("cannot create ethereum signer: %v", err)
	}

	okxSigner, err := ethereum.NewSigner(config.OkxPrivateKey, config.OkxBridgeConfigVersion)
	if err != nil {
		log.Fatalf("cannot create okx signer: %v", err)
	}

	ccdSigner, err := concordium.NewSigner(config.ConcordiumPrivateKey, config.ConcordiumBridgeConfigVersion, ccdObserver)
	if err != nil {
		log.Fatalf("cannot create concordium signer: %v", err)
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
		StellarObserver:  a.stellarObserver,
		EthereumSigner:   ethSigner,
		OkxSigner:        okxSigner,
		ConcordiumSigner: ccdSigner,
		StellarWithdrawalValidator: backend.StellarWithdrawalValidator{
			Session:          a.session.Clone(),
			WithdrawalWindow: config.WithdrawalWindow,
			Converter:        converter,
			CcdToken:         config.AssetMapping[0].ConcordiumToken,
		},
		ConcordiumWithdrawalValidator: backend.ConcordiumWithdrawalValidator{
			Observer:         ccdObserver,
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
		OkxWithdrawalValidator: backend.OkxWithdrawalValidator{
			Observer:          okxObserver,
			OkxFinalityBuffer: config.OkxFinalityBuffer,
			WithdrawalWindow:  config.WithdrawalWindow,
			Converter:         converter,
		},
	}
}

func (a *App) initHTTP(config Config, client *horizonclient.Client, ethObserver ethereum.Observer, ccdObserver concordium.Observer, okxObserver ethereum.Observer) {
	converter, err := backend.NewAssetConverter(config.AssetMapping)
	if err != nil {
		log.Fatal("unable to create asset converter", err)
	}

	httpServer, err := httpx.NewServer(httpx.ServerConfig{
		Ctx:                a.appCtx,
		Port:               config.Port,
		AdminPort:          config.AdminPort,
		PrometheusRegistry: a.prometheusRegistry,
		EthereumStellarWithdrawalHandler: &controllers.EthereumStellarWithdrawalHandler{
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
		StellarEthereumWithdrawalHandler: &controllers.StellarEthereumWithdrawalHandler{
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
		EthereumDepositHandler: &controllers.EthereumDepositHandler{
			Observer: ethObserver,
			Store:    a.NewStore(),
			StellarWithdrawalValidator: backend.StellarWithdrawalValidator{
				Session:          a.session.Clone(),
				WithdrawalWindow: config.WithdrawalWindow,
				Converter:        converter,
			},
			EthereumFinalityBuffer: config.EthereumFinalityBuffer,
			Token:                  config.AssetMapping[0].EthereumToken,
		},
		ConcordiumDepositHandler: &controllers.ConcordiumDepositHandler{
			Observer: ccdObserver,
			Store:    a.NewStore(),
		},
		ConcordiumStellarWithdrawalHandler: &controllers.ConcordiumStellarWithdrawalHandler{
			StellarClient: client,
			Observer:      ccdObserver,
			Store:         a.NewStore(),
			ConcordiumToStellarWithdrawalValidator: backend.StellarWithdrawalValidator{
				Session:          a.session.Clone(),
				WithdrawalWindow: config.WithdrawalWindow,
				Converter:        converter,
				CcdToken:         config.AssetMapping[0].ConcordiumToken,
			},
		},
		StellarConcordiumWithdrawalHandler: &controllers.StellarConcordiumWithdrawalHandler{
			Store: a.NewStore(),
			ConcordiumWithdrawalValidator: backend.ConcordiumWithdrawalValidator{
				Observer:         ccdObserver,
				WithdrawalWindow: config.WithdrawalWindow,
				Converter:        converter,
			},
		},
		OkxStellarWithdrawalHandler: &controllers.OkxStellarWithdrawalHandler{
			StellarClient: client,
			Observer:      okxObserver,
			Store:         a.NewStore(),
			StellarWithdrawalValidator: backend.StellarWithdrawalValidator{
				Session:          a.session.Clone(),
				WithdrawalWindow: config.WithdrawalWindow,
				Converter:        converter,
			},
			OkxFinalityBuffer: config.OkxFinalityBuffer,
		},
		StellarOkxWithdrawalHandler: &controllers.StellarOkxWithdrawalHandler{
			Store: a.NewStore(),
			OkxWithdrawalValidator: backend.EthereumWithdrawalValidator{
				Observer:               okxObserver,
				EthereumFinalityBuffer: config.OkxFinalityBuffer,
				WithdrawalWindow:       config.WithdrawalWindow,
				Converter:              converter,
			},
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

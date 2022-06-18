package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/support/db"
	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/backend"
	"github.com/stellar/starbridge/httpx"
	"github.com/stellar/starbridge/stellar/signer"
	"github.com/stellar/starbridge/stellar/txbuilder"
	"github.com/stellar/starbridge/stellar/txobserver"
	"github.com/stellar/starbridge/store"
)

type App struct {
	appCtx    context.Context
	cancelCtx context.CancelFunc

	httpServer *httpx.Server
	worker     *backend.Worker
	store      *store.DB

	prometheusRegistry *prometheus.Registry
}

type Config struct {
	Port      uint16 `valid:"-"`
	AdminPort uint16 `toml:"admin_port" valid:"-"`

	PostgresDSN string `toml:"postgres_dsn" valid:"-"`

	HorizonURL        string `toml:"horizon_url" valid:"-"`
	NetworkPassphrase string `toml:"network_passphrase" valid:"-"`

	MainAccountID   string `toml:"main_account_id" valid:"stellar_accountid"`
	SignerSecretKey string `toml:"signer_secret_key" valid:"optional,stellar_seed"`
}

func NewApp(config Config) *App {
	app := &App{
		store:              &store.DB{},
		prometheusRegistry: prometheus.NewRegistry(),
	}

	app.initGracefulShutdown()
	app.initHTTP(config)
	app.initWorker(config)
	app.initStore(config)
	app.initLogger()
	app.initPrometheus()

	return app
}

func (a *App) GetStore() *store.DB {
	return a.store
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

func (a *App) initStore(config Config) {
	session, err := db.Open("postgres", config.PostgresDSN)
	if err != nil {
		log.Fatalf("cannot open DB: %v", err)
	}

	a.store.Session = session
	err = a.store.InitSchema()
	if err != nil {
		log.Fatalf("cannot init DB: %v", err)
	}
}

func (a *App) initWorker(config Config) {
	var (
		signerKey *keypair.Full
		err       error
	)
	if config.SignerSecretKey != "" {
		signerKey, err = keypair.ParseFull(config.SignerSecretKey)
		if err != nil {
			log.Fatalf("cannot pase signer secret key: %v", err)
		}
	}

	a.worker = &backend.Worker{
		Ctx:   a.appCtx,
		Store: a.store,
		StellarBuilder: &txbuilder.Builder{
			HorizonURL:    config.HorizonURL,
			BridgeAccount: config.MainAccountID,
		},
		StellarSigner: &signer.Signer{
			NetworkPassphrase: config.NetworkPassphrase,
			Signer:            signerKey,
		},
		StellarObserver: txobserver.NewObserver(a.appCtx, horizonclient.DefaultTestNetClient, a.store),
	}
}

func (a *App) initHTTP(config Config) {
	httpServer, err := httpx.NewServer(httpx.ServerConfig{
		Ctx:                a.appCtx,
		Port:               config.Port,
		AdminPort:          config.AdminPort,
		PrometheusRegistry: a.prometheusRegistry,
		Store:              a.store,
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
	a.worker.Run()
}

func (a *App) Close() {
	a.cancelCtx()
}

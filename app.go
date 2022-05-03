package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/backend"
	"github.com/stellar/starbridge/httpx"
	"github.com/stellar/starbridge/stellar/signer"
	"github.com/stellar/starbridge/stellar/txbuilder"
	"github.com/stellar/starbridge/stellar/txobserver"
	"github.com/stellar/starbridge/store"
)

type App struct {
	httpServer *httpx.Server
	worker     *backend.Worker
	store      *store.Memory

	stellarObserver *txobserver.Observer

	prometheusRegistry *prometheus.Registry
}

type Config struct {
	Port      uint16
	AdminPort uint16
}

func NewApp(config Config) *App {
	app := &App{
		store:              &store.Memory{},
		prometheusRegistry: prometheus.NewRegistry(),
	}

	app.initHTTP(config)
	app.initWorker()
	app.initStellarTxObserver()
	app.initLogger()
	app.initPrometheus()

	return app
}

func (a *App) initPrometheus() {
	a.httpServer.RegisterMetrics(a.prometheusRegistry)
}

func (a *App) initLogger() {
	log.SetLevel(log.InfoLevel)
}

func (a *App) initWorker() {
	a.worker = &backend.Worker{
		Store: a.store,
		StellarBuilder: &txbuilder.Builder{
			BridgeAccount: "GBMULMDOT22YJ6SFUCADW7OQCQUBE5LMMHN6GXJ4A5P5IBOK56YYUK6M",
		},
		StellarSigner: &signer.Signer{
			NetworkPassphrase: network.TestNetworkPassphrase,
			SecretKey:         "SAV3VE7CMIDIY5GWPZ3WPTMXCD342CGRVKP2SHX4FHAU5D35QW7HNJLS",
		},
	}
}

func (a *App) initStellarTxObserver() {
	a.stellarObserver = &txobserver.Observer{
		Store:  a.store,
		Client: horizonclient.DefaultTestNetClient,
	}
}

func (a *App) initHTTP(config Config) {
	httpServer, err := httpx.NewServer(httpx.ServerConfig{
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
	err := a.worker.Run()
	if err != nil {
		log.WithField("error", err).Error("error running backend worker")
	}
}

// RunStellarTxObserver starts backend worker responsible for observing Stellar
// transactions
func (a *App) RunStellarTxObserver() {
	err := a.stellarObserver.Run()
	if err != nil {
		log.WithField("error", err).Error("error running stellar tx observer")
	}
}

package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/httpx"
)

type App struct {
	httpServer *httpx.Server

	prometheusRegistry *prometheus.Registry
}

type Config struct {
	Port      uint16
	AdminPort uint16
}

func NewApp(config Config) *App {
	app := &App{
		prometheusRegistry: prometheus.NewRegistry(),
	}

	app.initHTTP(config)
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

func (a *App) initHTTP(config Config) {
	httpServer, err := httpx.NewServer(httpx.ServerConfig{
		Port:               config.Port,
		AdminPort:          config.AdminPort,
		PrometheusRegistry: a.prometheusRegistry,
	})
	if err != nil {
		log.Fatal("unable to create http server", err)
	}
	a.httpServer = httpServer
}

// Run starts all services, including http server
func (a *App) Run() {
	err := a.httpServer.Serve()
	if err != nil {
		log.WithField("error", err).Error("error running http server")
	}
}

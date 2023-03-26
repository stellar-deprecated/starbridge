package httpx

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/stellar/go/support/errors"
	stellarhttp "github.com/stellar/go/support/http"
	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/controllers"
	"github.com/stellar/starbridge/html"
)

type ServerMetrics struct {
	RequestDurationSummary *prometheus.SummaryVec
}

type TLSConfig struct {
	CertPath, KeyPath string
}

type ServerConfig struct {
	Ctx                context.Context
	Port               uint16
	AdminPort          uint16
	TLSConfig          *TLSConfig
	PrometheusRegistry *prometheus.Registry

	EthereumStellarWithdrawalHandler   *controllers.EthereumStellarWithdrawalHandler
	StellarEthereumWithdrawalHandler   *controllers.StellarEthereumWithdrawalHandler
	OkxStellarWithdrawalHandler        *controllers.OkxStellarWithdrawalHandler
	StellarOkxWithdrawalHandler        *controllers.StellarOkxWithdrawalHandler
	ConcordiumStellarWithdrawalHandler *controllers.ConcordiumStellarWithdrawalHandler
	StellarConcordiumWithdrawalHandler *controllers.StellarConcordiumWithdrawalHandler

	StellarRefundHandler  *controllers.StellarRefundHandler
	EthereumRefundHandler *controllers.EthereumRefundHandler

	EthereumDepositHandler   *controllers.EthereumDepositHandler
	ConcordiumDepositHandler *controllers.ConcordiumDepositHandler
}

type Server struct {
	Metrics *ServerMetrics

	ctx context.Context

	server      *http.Server
	adminServer *http.Server

	tlsConfig          *TLSConfig
	prometheusRegistry *prometheus.Registry
}

func NewServer(serverConfig ServerConfig) (*Server, error) {
	metrics := &ServerMetrics{
		RequestDurationSummary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: "starbridge", Subsystem: "http", Name: "requests_duration_seconds",
				Help: "HTTP requests durations, sliding window = 10m",
			},
			[]string{"status", "method"},
		),
	}

	server := &Server{
		Metrics: metrics,

		ctx: serverConfig.Ctx,

		prometheusRegistry: serverConfig.PrometheusRegistry,
		tlsConfig:          serverConfig.TLSConfig,

		server: &http.Server{
			Addr:        fmt.Sprintf(":%d", serverConfig.Port),
			ReadTimeout: 5 * time.Second,
		},
	}
	server.initMux(serverConfig)

	if serverConfig.AdminPort != 0 {
		server.adminServer = &http.Server{
			Addr:        fmt.Sprintf(":%d", serverConfig.AdminPort),
			ReadTimeout: 5 * time.Second,
		}
		server.initAdminMux()
	}

	return server, nil
}

func (s *Server) initMux(serverConfig ServerConfig) {
	mux := stellarhttp.NewAPIMux(log.DefaultLogger)

	// Public middlewares
	mux.Use(middleware.StripSlashes)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Date"},
	})
	mux.Use(c.Handler)
	mux.Use(middleware.NoCache)
	mux.Use(prometheusMiddleware(s.Metrics))
	mux.Use(middleware.Timeout(10 * time.Second))

	// Public routes
	mux.Method(http.MethodPost, "/ethereum/deposit", serverConfig.EthereumDepositHandler)
	mux.Method(http.MethodPost, "/concordium/deposit", serverConfig.ConcordiumDepositHandler)
	//mux.Method(http.MethodPost, "/okx/deposit", serverConfig.EthereumDepositHandler)

	mux.Method(http.MethodPost, "/ethereum/withdraw/stellar", serverConfig.EthereumStellarWithdrawalHandler)
	mux.Method(http.MethodPost, "/stellar/withdraw/ethereum", serverConfig.StellarEthereumWithdrawalHandler)
	mux.Method(http.MethodPost, "/concordium/withdraw/stellar", serverConfig.ConcordiumStellarWithdrawalHandler)
	mux.Method(http.MethodPost, "/stellar/withdraw/concordium", serverConfig.StellarConcordiumWithdrawalHandler)
	mux.Method(http.MethodPost, "/okx/withdraw/stellar", serverConfig.OkxStellarWithdrawalHandler)
	mux.Method(http.MethodPost, "/stellar/withdraw/okx", serverConfig.StellarOkxWithdrawalHandler)

	mux.Method(http.MethodPost, "/ethereum/refund", serverConfig.EthereumRefundHandler)
	mux.Method(http.MethodPost, "/stellar/refund", serverConfig.StellarRefundHandler)

	staticServer := http.FileServer(http.FS(html.Files))
	mux.Handle("/*", staticServer)

	s.server.Handler = mux
}

func (s *Server) initAdminMux() {
	adminMux := stellarhttp.NewAPIMux(log.DefaultLogger)

	// Admin middlewares
	adminMux.Use(middleware.NoCache)
	adminMux.Use(prometheusMiddleware(s.Metrics))
	adminMux.Use(middleware.Timeout(10 * time.Second))

	// Admin routes
	adminMux.Get("/metrics", promhttp.HandlerFor(s.prometheusRegistry, promhttp.HandlerOpts{}).ServeHTTP)

	s.adminServer.Handler = adminMux
}

// RegisterMetrics registers the prometheus metrics
func (s *Server) RegisterMetrics(registry *prometheus.Registry) {
	registry.MustRegister(s.Metrics.RequestDurationSummary)
}

func (s *Server) Serve() error {
	go func() {
		<-s.ctx.Done()
		if s.adminServer != nil {
			_ = s.adminServer.Shutdown(context.Background())
		}
		_ = s.server.Shutdown(context.Background())
	}()

	if s.adminServer != nil {
		go func() {
			log.Infof("starting admin server on %s", s.adminServer.Addr)
			if err := s.adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Warn(errors.Wrap(err, "error in internalServer.ListenAndServe()"))
			}
		}()
	}

	log.Infof("starting server on %s", s.server.Addr)
	var err error
	if s.tlsConfig != nil {
		err = s.server.ListenAndServeTLS(s.tlsConfig.CertPath, s.tlsConfig.KeyPath)
	} else {
		err = s.server.ListenAndServe()
	}

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	var wg sync.WaitGroup
	defer wg.Wait()
	if s.adminServer != nil {
		wg.Add(1)
		go func() {
			err := s.adminServer.Shutdown(ctx)
			if err != nil {
				log.Warn(errors.Wrap(err, "error in adminServer.Shutdown()"))
			}
			wg.Done()
		}()
	}
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

package httpx

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
)

type ServerMetrics struct {
	RequestDurationSummary *prometheus.SummaryVec
}

type TLSConfig struct {
	CertPath, KeyPath string
}

type ServerConfig struct {
	Port               uint16
	AdminPort          uint16
	TLSConfig          *TLSConfig
	PrometheusRegistry *prometheus.Registry
}

type Server struct {
	Metrics *ServerMetrics

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
			[]string{"status", "route", "method"},
		),
	}

	server := &Server{
		Metrics: metrics,

		prometheusRegistry: serverConfig.PrometheusRegistry,
		tlsConfig:          serverConfig.TLSConfig,

		server: &http.Server{
			Addr:        fmt.Sprintf(":%d", serverConfig.Port),
			ReadTimeout: 5 * time.Second,
		},
	}
	server.initMux()

	if serverConfig.AdminPort != 0 {
		server.adminServer = &http.Server{
			Addr:        fmt.Sprintf(":%d", serverConfig.AdminPort),
			ReadTimeout: 5 * time.Second,
		}
		server.initAdminMux()
	}

	return server, nil
}

func (s *Server) initMux() {
	mux := chi.NewMux()

	// Public middlewares
	mux.Use(middleware.StripSlashes)
	mux.Use(middleware.NoCache)
	mux.Use(middleware.RequestID)
	mux.Use(loggerMiddleware(s.Metrics))
	mux.Use(middleware.Timeout(10 * time.Second))
	mux.Use(middleware.Recoverer)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Date"},
	})
	mux.Use(c.Handler)

	// Public routes
	mux.Method(http.MethodGet, "/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "world.")
	}))

	s.server.Handler = mux
}

func (s *Server) initAdminMux() {
	adminMux := chi.NewMux()

	// Admin middlewares
	adminMux.Use(middleware.NoCache)
	adminMux.Use(loggerMiddleware(s.Metrics))
	adminMux.Use(middleware.Timeout(10 * time.Second))
	adminMux.Use(middleware.Recoverer)

	// Admin routes
	adminMux.Get("/metrics", promhttp.HandlerFor(s.prometheusRegistry, promhttp.HandlerOpts{}).ServeHTTP)

	s.adminServer.Handler = adminMux
}

// RegisterMetrics registers the prometheus metrics
func (s *Server) RegisterMetrics(registry *prometheus.Registry) {
	registry.MustRegister(s.Metrics.RequestDurationSummary)
}

func (s *Server) Serve() error {
	if s.adminServer != nil {
		go func() {
			log.Infof("Starting admin server on %s", s.adminServer.Addr)
			if err := s.adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Warn(errors.Wrap(err, "error in internalServer.ListenAndServe()"))
			}
		}()
	}

	log.Infof("Starting server on %s", s.server.Addr)
	var err error
	if s.tlsConfig != nil {
		err = s.server.ListenAndServeTLS(s.tlsConfig.CertPath, s.tlsConfig.KeyPath)
	} else {
		err = s.server.ListenAndServe()
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

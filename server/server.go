// Package server provides functionality to easily set up a subscriber of pubsub events.
//
// The server holds all the clients it needs. The clients should be set up in the Create method.
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"template-subscriber-go/client/database"
	"template-subscriber-go/client/pubsub"
	"template-subscriber-go/config"
	"template-subscriber-go/monitoring/metrics"
	"template-subscriber-go/monitoring/trace"
	"template-subscriber-go/server/internal/event"
	"template-subscriber-go/server/internal/handler"

	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Server holds an HTTP server, config and all the clients.
type Server struct {
	Config          *config.Config
	HTTP            *http.Server
	DB              *database.Client
	PubSub          *pubsub.Client
	TracerProvider  *tracesdk.TracerProvider
	MetricsProvider *metricsdk.MeterProvider
}

// Create sets up a server with necessary all clients.
// Returns an error if an error occurs.
func (s *Server) Create(ctx context.Context, config *config.Config) error {
	var dbClient database.Client
	if err := dbClient.Init(ctx, config); err != nil {
		return fmt.Errorf("database client: %w", err)
	}

	var psClient pubsub.Client
	if err := psClient.Init(ctx, config); err != nil {
		return fmt.Errorf("pubsub client: %w", err)
	}

	s.DB = &dbClient
	s.PubSub = &psClient
	s.Config = config
	s.HTTP = &http.Server{
		Addr: fmt.Sprintf(":%s", s.Config.Port),
	}

	return nil
}

// Serve starts subscribing for messages.
// It also makes sure that the server gracefully shuts down on exit.
// Returns an error if an error occurs.
func (s *Server) Serve(ctx context.Context, errc chan<- error) {
	s.addTracingAndMetrics(errc)

	go s.serveHTTP(errc)
	go s.subscribeAndListen(ctx, errc)

	log.Info("Ready")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Info("Shutdown signal received")

	s.shutdown(ctx)
}

func (s *Server) serveHTTP(errc chan<- error) {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/_healthz", handler.Healthz)

	if err := s.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		errc <- err
	}
}

func (s *Server) addTracingAndMetrics(errc chan<- error) {
	var err error
	s.TracerProvider, err = trace.TracerProvider(s.Config)
	if err != nil {
		errc <- fmt.Errorf("init global tracer: %w", err)
	}

	s.MetricsProvider, err = metrics.MetricsProvider(s.Config)
	if err != nil {
		errc <- fmt.Errorf("init metrics: %w", err)
	}

}

func (s *Server) subscribeAndListen(ctx context.Context, errc chan<- error) {
	for _, e := range event.GetPubSubEvents(s.DB) {
		go func(e event.PubSubEvent) {
			e.SubscribeAndListen(ctx, s.PubSub, errc)
		}(e)
	}
	for _, e := range event.GetAppEvents() {
		go func(e event.AppEvent) {
			e.SubscribeAndListen(ctx)
		}(e)
	}

}

func (s *Server) shutdown(ctx context.Context) {
	if err := s.TracerProvider.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}

	if err := s.MetricsProvider.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}

	if err := s.DB.Close(); err != nil {
		log.Error(err.Error())
	}

	if err := s.HTTP.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}

}

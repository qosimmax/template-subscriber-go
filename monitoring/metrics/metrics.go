// Package metrics sets up and handles our prometheus collectors.
package metrics

import (
	"context"
	"template-subscriber-go/config"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

var (
	messagesReceived api.Int64Counter
	errorsOccurred   api.Int64Counter
	timeToProcess    api.Float64Histogram
)

// MetricsProvider tells prometheus to set up collectors.
func MetricsProvider(cfg *config.Config) (*metric.MeterProvider, error) {
	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(cfg.ServiceName)

	messagesReceived, _ = meter.Int64Counter("messages_received",
		api.WithDescription("Number of messages received from PubSub."),
		api.WithUnit("{call}"),
	)

	errorsOccurred, _ = meter.Int64Counter("errors_occurred",
		api.WithDescription("Number of errors occurred when processing PubSub messages."),
		api.WithUnit("{call}"),
	)

	timeToProcess, _ = meter.Float64Histogram("task_duration",
		api.WithDescription("Amount of time spent processing."),
		api.WithUnit("s"),
	)

	otel.SetMeterProvider(provider)

	return provider, nil
}

// ReceivedMessage records number of messages of each type received.
func ReceivedMessage(ctx context.Context, msgType string, t int64) {
	opt := api.WithAttributes(
		attribute.Key("message_type").String(msgType),
	)
	messagesReceived.Add(ctx, t, opt)
}

// OccurredError records number of errors occurred while processing messages
// of each type.
func OccurredError(ctx context.Context, msgType string) {
	opt := api.WithAttributes(
		attribute.Key("not_processed_message_type").String(msgType),
	)
	errorsOccurred.Add(ctx, 1, opt)
}

// ObserveTimeToProcess records amount of time spent processing messages.
func ObserveTimeToProcess(ctx context.Context, t float64) {
	timeToProcess.Record(ctx, t)
}

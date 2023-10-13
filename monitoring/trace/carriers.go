// Package trace contains functions for tracing requests.
package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// InjectIntoCarrier returns a textMapCarrier, basically a map[string]string,
// which can be used to transmit a span context to another service with ExtractFromCarrier.
func InjectIntoCarrier(ctx context.Context, carrier propagation.TextMapCarrier) {
	prop := otel.GetTextMapPropagator()
	prop.Inject(ctx, carrier)

}

// ExtractFromCarrier returns a span with context passed  by the carrier.
// ctx should not already have span in it.
func ExtractFromCarrier(ctx context.Context, carrier propagation.TextMapCarrier, spanName string) context.Context {
	prop := otel.GetTextMapPropagator()
	return prop.Extract(ctx, carrier)
}

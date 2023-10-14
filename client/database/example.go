package database

import (
	"context"
	"template-subscriber-go/example"
	"time"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("database")

func (c *Client) RecordExampleData(ctx context.Context, exampleData example.Data) error {
	ctx, span := tracer.Start(ctx, "RecordExampleData")
	defer span.End()
	_, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return nil
}

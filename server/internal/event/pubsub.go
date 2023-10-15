package event

import (
	"context"
	"errors"
	"fmt"
	"template-subscriber-go/client/pubsub"
	"template-subscriber-go/example"
	"template-subscriber-go/monitoring/metrics"
	"template-subscriber-go/monitoring/trace"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

// PubSubEvents contains a slice of PubSubEvent.
type PubSubEvents []PubSubEvent

// PubSubEvent contains the data for a PubSub event type.
type PubSubEvent struct {
	Name             string
	Queue            string
	SubscriptionName string
	Handler          Handler
	Subscription     nats.JetStreamContext
}

// SubscribeAndListen subscribes to a PubSubEvent.
func (e *PubSubEvent) SubscribeAndListen(ctx context.Context, c *pubsub.Client, errc chan<- error) {
	e.Subscription = c
	go e.receive(ctx, errc)
}

func (e *PubSubEvent) receive(ctx context.Context, errc chan<- error) {
	var tracer = otel.Tracer(e.Name)

	handler := func(ctx context.Context, msg *nats.Msg) {
		ctx = trace.ExtractFromCarrier(ctx, propagation.HeaderCarrier(msg.Header), e.Name)
		ctx, span := tracer.Start(ctx, e.Name)
		defer span.End()

		metrics.ReceivedMessage(ctx, e.Name, 1)
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			metrics.ObserveTimeToProcess(ctx, duration.Seconds())
		}()

		var errNonRecoverable example.ErrNonRecoverable
		var errExpected example.ErrExpected

		err := e.Handler.Handle(ctx, msg.Data)
		if err != nil {
			// If the error is not an expected error, log and record the error
			if !errors.As(err, &errExpected) {
				log.Error(err.Error())
				span.SetStatus(codes.Error, "handle event failed")
				span.RecordError(err)
				metrics.OccurredError(ctx, e.Name)
			}

			// If the error is not a non-recoverable error, it means it is
			// recoverable, so return before acking
			if !errors.As(err, &errNonRecoverable) {
				return
			}
		}

		_ = msg.Ack()
	}

	sub, err := e.Subscription.PullSubscribe(e.SubscriptionName, e.Queue, nats.DeliverAll())
	if err != nil {
		errc <- fmt.Errorf("subscription receive(%s): %w", e.SubscriptionName, err)
	}

	//@TODO get worker numbers from env
	queue := make(chan *nats.Msg, 1)
	worker := func(ctx context.Context, queue chan *nats.Msg) {
		for msg := range queue {
			handler(ctx, msg)
		}
	}

	for i := 0; i < 10; i++ {
		go worker(ctx, queue)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msgs, _ := sub.Fetch(10, nats.MaxWait(50*time.Millisecond))
		for _, msg := range msgs {
			queue <- msg
		}
	}

}

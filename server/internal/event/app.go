package event

import (
	"context"
	"errors"
	"template-subscriber-go/example"
	"time"

	log "github.com/sirupsen/logrus"
)

// AppEvents contains a slice of AppEvent.
type AppEvents []AppEvent

// AppEvent contains the data for an in-app event type.
type AppEvent struct {
	Name    string
	Rate    time.Duration
	Handler Handler
}

// SubscribeAndListen subscribes to an AppEvent.
func (e *AppEvent) SubscribeAndListen(ctx context.Context) {
	for t := range time.Tick(e.Rate) {
		go func(t time.Time) {

			var errExpected example.ErrExpected
			err := e.Handler.Handle(ctx, nil)
			if err != nil && !errors.As(err, &errExpected) {
				log.Error(t, err.Error())
			}
		}(t)
	}
}

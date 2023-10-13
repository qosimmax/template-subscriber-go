// Package event handles configuration and setup for receiving events.
//
// Events to subscribe to should be defined in GetPubSubEvents
package event

import (
	"context"
	"template-subscriber-go/client/database"
	"template-subscriber-go/server/internal/handler"
)

// Handler is an interface that all event handles must implement.
type Handler interface {
	Handle(ctx context.Context, data []byte) error
}

// GetPubSubEvents describes all the pubsub events to listen to.
func GetPubSubEvents(db *database.Client) PubSubEvents {
	// Define your google pubsub events here
	psEvents := PubSubEvents{
		PubSubEvent{
			Name:             "Example",
			Queue:            "example",
			SubscriptionName: "example",
			Handler:          handler.Example{},
		},
	}

	return psEvents
}

// GetAppEvents describes all the app events to listen to.
func GetAppEvents() AppEvents {
	appEvents := AppEvents{}

	return appEvents
}

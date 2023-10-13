package pubsub

import (
	"context"

	"github.com/nats-io/nats.go"

	"template-subscriber-go/config"
)

// Client holds the PubSub client.
type Client struct {
	nats.JetStreamContext
}

// Init sets up a new pubsub client.
func (c *Client) Init(ctx context.Context, config *config.Config) error {
	nc, err := nats.Connect(config.NatsURL)
	if err != nil {
		return err
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(10000))
	if err != nil {
		return err
	}

	if err := c.createStreams(js); err != nil {
		return err
	}

	c.JetStreamContext = js
	return nil
}

func (c *Client) createStreams(js nats.JetStreamContext) error {

	type streamConf struct {
		Name     string
		Subjects []string
	}

	streams := []streamConf{
		{
			Name:     "example",
			Subjects: []string{"example"},
		},
	}

	for _, stream := range streams {
		_, err := js.StreamInfo(stream.Name)
		if err != nil {
			if err == nats.ErrStreamNotFound {
				_, err = js.AddStream(&nats.StreamConfig{
					Name:     stream.Name,
					Subjects: stream.Subjects,
				})

				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

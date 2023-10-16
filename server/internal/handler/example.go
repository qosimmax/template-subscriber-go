// Package handler contains handlers for events.
package handler

import (
	"context"
	"fmt"
	"template-subscriber-go/example"
	"template-subscriber-go/example/pb/fakeapi"

	"google.golang.org/protobuf/proto"
)

// Example is an example event.
type Example struct {
	DB example.DataRecorder
}

// Handle is the handler for the example event.
func (e Example) Handle(ctx context.Context, data []byte) error {
	fakeData := &fakeapi.FakeData{}
	err := proto.Unmarshal(data, fakeData)
	if err != nil {
		return example.ErrNonRecoverable{
			Err: fmt.Errorf("failed to unmarshal example data: %w", err),
		}
	}

	exampleData := example.Data{
		IsFake: fakeData.IsFake,
		Date:   fakeData.GetDate().AsTime(),
	}

	// Do stuff here
	err = e.DB.RecordExampleData(ctx, exampleData)
	if err != nil {
		return err
	}

	return nil
}

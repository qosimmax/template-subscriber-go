// Package handler contains handlers for events.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"template-subscriber-go/example"
)

// Example is an example event.
type Example struct{}

// Handle is the handler for the example event.
func (e Example) Handle(ctx context.Context, data []byte) error {
	var exampleData example.Data
	err := json.Unmarshal(data, &exampleData)
	if err != nil {
		return example.ErrNonRecoverable{
			Err: fmt.Errorf("failed to unmarshal example data: %w", err),
		}
	}
	// Do stuff here
	return nil
}

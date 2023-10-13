package example

import (
	"context"
	"time"
)

// Data contains example data indicating if the data is fake
// or not.
type Data struct {
	IsFake bool      `json:"isFake"`
	Date   time.Time `json:"date"`
}

// DataRecorder is an interface for recording example data.
type DataRecorder interface {
	RecordExampleData(ctx context.Context, exampleData Data) error
}

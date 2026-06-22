package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Convert Time to Timestamp
func ToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

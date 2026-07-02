package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Convert Timestamp to Time
func ToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

// Convert Time to Timestamp
func ToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

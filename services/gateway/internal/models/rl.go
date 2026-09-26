package models

import "time"

type (
	LimitByFixedWindowRL struct {
		IPAddress string
		Namespace string
		Window    time.Duration
		Limit     uint16
	}

	LimitBySlidingWindowRL struct {
		AuthID    uint64
		Namespace string
		Now       time.Time
		Window    time.Duration
		Limit     uint16
	}
)

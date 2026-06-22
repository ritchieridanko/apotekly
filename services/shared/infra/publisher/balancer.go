package publisher

import "github.com/segmentio/kafka-go"

func balancer(name string) kafka.Balancer {
	switch name {
	case "murmur2":
		return &kafka.Murmur2Balancer{Consistent: true}
	default:
		return &kafka.RoundRobin{}
	}
}

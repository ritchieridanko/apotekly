package utils

import "github.com/google/uuid"

// Create a new random UUID
func GenerateUUID() uuid.UUID {
	return uuid.New()
}

// Create a new random UUID v7
func GenerateUUIDv7() (uuid.UUID, error) {
	return uuid.NewV7()
}

// Convert value to UUID
// NOTE: Panic if value cannot be converted
func ToUUID(value string) uuid.UUID {
	return uuid.MustParse(value)
}

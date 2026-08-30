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

// Create a new random UUID v7
// NOTE: Panic if UUID fails to create
func MustGenerateUUIDv7() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

// Convert value to UUID
// NOTE: Panic if value cannot be converted
func ToUUID(value string) uuid.UUID {
	return uuid.MustParse(value)
}

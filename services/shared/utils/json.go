package utils

import "encoding/json"

// Convert value to JSON Raw Message
func ToJSONRawMessage(value any) (json.RawMessage, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

// Convert JSON Raw Message back to original value
func FromJSONRawMessage[T any](value json.RawMessage) (*T, error) {
	var v T
	if err := json.Unmarshal(value, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

package utils

import "encoding/json"

// Convert JSON Raw Message back to original value
func FromJSONRawMessage[T any](value json.RawMessage) (*T, error) {
	var v T
	if err := json.Unmarshal(value, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// Convert JSON Raw Message to []byte
func ToByte(rm *json.RawMessage) []byte {
	if rm == nil {
		return nil
	}
	return *rm
}

// Convert []byte to JSON Raw Message
func ToJSON(b []byte) *json.RawMessage {
	if len(b) == 0 {
		return nil
	}
	res := json.RawMessage(b)
	return &res
}

// Convert value to JSON Raw Message
func ToJSONRawMessage(value any) (json.RawMessage, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

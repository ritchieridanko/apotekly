package utils

import (
	"fmt"
	"strconv"
)

// Convert value to uint64
func ToUint64(value any) (uint64, error) {
	s, ok := value.(string)
	if !ok {
		return 0, fmt.Errorf("unable to convert to uint64: %v (type: %T)", value, value)
	}
	return strconv.ParseUint(s, 10, 64)
}

package utils

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titlecaser = cases.Title(language.English)

// Convert string to all lowercase
// NOTE: Return nil if s is nil
func ToLowerPtr(s *string) *string {
	if s == nil {
		return nil
	}
	res := strings.ToLower(*s)
	return &res
}

// Convert string to titlecase
func ToTitlecase(s string) string {
	values := strings.Fields(s)
	if len(values) == 0 {
		return ""
	}
	return titlecaser.String(strings.Join(values, " "))
}

// Convert string to titlecase
// NOTE: Return nil if s is nil
func ToTitlecasePtr(s *string) *string {
	if s == nil {
		return nil
	}
	res := ToTitlecase(*s)
	return &res
}

// Convert value to uint64
func ToUint64(value any) (uint64, error) {
	s, ok := value.(string)
	if !ok {
		return 0, fmt.Errorf("unable to convert to uint64: %v (type: %T)", value, value)
	}
	return strconv.ParseUint(s, 10, 64)
}

// Strip string of leading and trailing whitespaces
// NOTE: Return nil if s is nil
func TrimSpacePtr(s *string) *string {
	if s == nil {
		return nil
	}
	res := strings.TrimSpace(*s)
	return &res
}

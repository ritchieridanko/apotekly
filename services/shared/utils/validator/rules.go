package validator

import "regexp"

var (
	rgxEmail        *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rgxLowercase    *regexp.Regexp = regexp.MustCompile(`[a-z]`)
	rgxNumber       *regexp.Regexp = regexp.MustCompile(`[0-9]`)
	rgxPhone        *regexp.Regexp = regexp.MustCompile(`^0\d{7,15}$`)
	rgxSpecialChars *regexp.Regexp = regexp.MustCompile(`[!@#$%^&*()_+\-={};:'"/\\|,.<>?]`)
	rgxUppercase    *regexp.Regexp = regexp.MustCompile(`[A-Z]`)
)

const (
	bcryptMaxBytes     int = 72
	nameMaxLength      int = 100
	nameMinLength      int = 3
	passwordMaxLength  int = 50
	passwordMinLength  int = 8
	userAgentMaxLength int = 512
)

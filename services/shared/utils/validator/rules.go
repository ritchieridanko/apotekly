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
	bcryptMaxBytes       int = 72
	labelMaxLength       int = 50
	labelMinLength       int = 3
	nameMaxLength        int = 100
	nameMinLength        int = 3
	notesMaxLength       int = 100
	passwordMaxLength    int = 50
	passwordMinLength    int = 8
	postalCodeMaxLength  int = 6
	postalCodeMinLength  int = 4
	streetMaxLength      int = 250
	streetMinLength      int = 3
	subdivisionMaxLength int = 250
	userAgentMaxLength   int = 512

	maxLatitude  float64 = 90
	minLatitude  float64 = -90
	maxLongitude float64 = 180
	minLongitude float64 = -180
)

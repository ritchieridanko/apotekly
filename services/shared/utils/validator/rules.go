package validator

import "regexp"

var (
	rgxEmail            *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rgxLowercase        *regexp.Regexp = regexp.MustCompile(`[a-z]`)
	rgxNumber           *regexp.Regexp = regexp.MustCompile(`[0-9]`)
	rgxPhone            *regexp.Regexp = regexp.MustCompile(`^0\d{7,15}$`)
	rgxPrice            *regexp.Regexp = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	rgxSpecialChars     *regexp.Regexp = regexp.MustCompile(`[!@#$%^&*()_+\-={};:'"/\\|,.<>?]`)
	rgxStorageURL       *regexp.Regexp = regexp.MustCompile(`^https://res\.cloudinary\.com/dta3lzmww/image/upload/(?:v\d+/)?(.+)$`)
	rgxTimeRange        *regexp.Regexp = regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d-(?:[01]\d|2[0-3]):[0-5]\d$`)
	rgxUppercase        *regexp.Regexp = regexp.MustCompile(`[A-Z]`)
	rgxZeroDecimalPrice *regexp.Regexp = regexp.MustCompile(`^\d+$`)
)

const (
	bcryptMaxBytes          int = 72
	descMaxLength           int = 1000
	dosageStrengthMaxLength int = 20
	labelMaxLength          int = 50
	labelMinLength          int = 3
	nameMaxLength           int = 255
	nameMinLength           int = 3
	notesMaxLength          int = 100
	passwordMaxLength       int = 50
	passwordMinLength       int = 8
	postalCodeMaxLength     int = 6
	postalCodeMinLength     int = 4
	regIdentifierMaxLength  int = 50
	regIdentifierMinLength  int = 2
	searchMaxLength         int = 255
	searchMinLength         int = 3
	subdivisionMaxLength    int = 255
	userAgentMaxLength      int = 512

	maxRadius uint32 = 100_000
	minRadius uint32 = 100

	maxLatitude  float64 = 90
	minLatitude  float64 = -90
	maxLongitude float64 = 180
	minLongitude float64 = -180
)

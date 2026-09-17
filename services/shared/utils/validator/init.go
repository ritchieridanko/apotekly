package validator

import (
	"encoding/json"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ritchieridanko/apotekly/services/shared/utils"
)

type Validator struct{}

func Init() *Validator {
	return &Validator{}
}

func (v *Validator) AddrLabel(value string) (bool, string) {
	length := utf8.RuneCountInString(value)
	if length < labelMinLength {
		return false, "Label must be at least " + strconv.Itoa(labelMinLength) + " characters"
	}
	if length > labelMaxLength {
		return false, "Label must not exceed " + strconv.Itoa(labelMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) AddrNotes(value string) (bool, string) {
	if utf8.RuneCountInString(value) > notesMaxLength {
		return false, "Notes must not exceed " + strconv.Itoa(notesMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) AddrSubdivision(value string) (bool, string) {
	if utf8.RuneCountInString(value) > subdivisionMaxLength {
		return false, "Subdivision name must not exceed " + strconv.Itoa(subdivisionMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) Birthdate(value time.Time) (bool, string) {
	if value.After(time.Now().UTC()) {
		return false, "Birthdate is invalid: " + value.Format("2 Jan 2006")
	}
	return true, ""
}

func (v *Validator) Country(value string, name string) (bool, string) {
	if _, exists := countries[value]; !exists {
		return false, name + " is invalid: " + utils.ToTitlecase(value)
	}
	return true, ""
}

func (v *Validator) Currency(value string) (bool, string) {
	if _, exists := currencies[value]; !exists {
		return false, "Currency is invalid: " + strings.ToUpper(value)
	}
	return true, ""
}

func (v *Validator) Description(value string) (bool, string) {
	if utf8.RuneCountInString(value) > descMaxLength {
		return false, "Description must not exceed " + strconv.Itoa(descMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) DosageForm(value string) (bool, string) {
	if _, exists := dosageForms[value]; !exists {
		return false, "Dosage form is invalid: " + utils.ToTitlecase(value)
	}
	return true, ""
}

func (v *Validator) DosageStrength(value string) (bool, string) {
	if utf8.RuneCountInString(value) > dosageStrengthMaxLength {
		return false, "Strength must not exceed " + strconv.Itoa(dosageStrengthMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) Email(value string) (bool, string) {
	if !rgxEmail.MatchString(value) {
		return false, "Email is invalid: " + value
	}
	return true, ""
}

func (v *Validator) IPAddress(value string) (bool, string) {
	if ip := net.ParseIP(value); ip == nil {
		return false, "IP address is invalid: " + value
	}
	return true, ""
}

func (v *Validator) Latitude(value float64) (bool, string) {
	if value < minLatitude || value > maxLatitude {
		return false, "Latitude is invalid: " + strconv.FormatFloat(value, 'f', 6, 64)
	}
	return true, ""
}

func (v *Validator) Longitude(value float64) (bool, string) {
	if value < minLongitude || value > maxLongitude {
		return false, "Longitude is invalid: " + strconv.FormatFloat(value, 'f', 6, 64)
	}
	return true, ""
}

func (v *Validator) Measurement(value float64, name string) (bool, string) {
	if value < 0 {
		return false, name + " is invalid: " + strconv.FormatFloat(value, 'f', 2, 64)
	}
	return true, ""
}

func (v *Validator) Name(value, name string) (bool, string) {
	length := utf8.RuneCountInString(value)
	if length < nameMinLength {
		return false, name + " must be at least " + strconv.Itoa(nameMinLength) + " characters"
	}
	if length > nameMaxLength {
		return false, name + " must not exceed " + strconv.Itoa(nameMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) OnlineHours(value json.RawMessage) (bool, string) {
	data, err := utils.FromJSONRawMessage[map[string][]string](value)
	if err != nil {
		return false, "Online hours are invalid"
	}
	for day, hourRanges := range *data {
		if _, exists := days[day]; !exists {
			return false, "Day in online hours is invalid: " + utils.ToTitlecase(day)
		}
		for _, hourRange := range hourRanges {
			why := "Range in online hours is invalid: " + hourRange
			if !rgxTimeRange.MatchString(hourRange) {
				return false, why
			}

			parts := strings.SplitN(hourRange, "-", 2)
			start, err := time.Parse("15:04", parts[0])
			if err != nil {
				return false, why
			}

			end, err := time.Parse("15:04", parts[1])
			if err != nil {
				return false, why
			}
			if !start.Before(end) {
				return false, why
			}
		}
	}
	return true, ""
}

func (v *Validator) PackSize(value int) (bool, string) {
	if value <= 0 {
		return false, "Pack size is invalid: " + strconv.Itoa(value)
	}
	return true, ""
}

func (v *Validator) PackUnit(value string) (bool, string) {
	if _, exists := packUnits[value]; !exists {
		return false, "Pack unit is invalid: " + utils.ToTitlecase(value)
	}
	return true, ""
}

func (v *Validator) Password(value string) (bool, string) {
	length := utf8.RuneCountInString(value)
	if length < passwordMinLength {
		return false, "Password must be at least " + strconv.Itoa(passwordMinLength) + " characters"
	}
	if length > passwordMaxLength {
		return false, "Password must not exceed " + strconv.Itoa(passwordMaxLength) + " characters"
	}
	if !rgxLowercase.MatchString(value) {
		return false, "Password must include at least one lowercase letter"
	}
	if !rgxUppercase.MatchString(value) {
		return false, "Password must include at least one uppercase letter"
	}
	if !rgxNumber.MatchString(value) {
		return false, "Password must include at least one number"
	}
	if !rgxSpecialChars.MatchString(value) {
		return false, "Password must include at least one special character: " + rgxSpecialChars.String()
	}
	if len([]byte(value)) > bcryptMaxBytes {
		return false, "Password too long"
	}
	return true, ""
}

func (v *Validator) Phone(value string) (bool, string) {
	if !rgxPhone.MatchString(value) {
		return false, "Phone is invalid: " + value
	}
	return true, ""
}

func (v *Validator) PostalCode(value string) (bool, string) {
	length := utf8.RuneCountInString(value)
	if length < postalCodeMinLength {
		return false, "Postal code must be at least " + strconv.Itoa(postalCodeMinLength) + " characters"
	}
	if length > postalCodeMaxLength {
		return false, "Postal code must not exceed " + strconv.Itoa(postalCodeMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) Price(value, currency string) (bool, string) {
	if _, exists := utils.ZeroDecimalCurrencies[currency]; exists {
		if !rgxZeroDecimalPrice.MatchString(value) {
			return false, "Price is invalid: " + value
		}
		return true, ""
	}
	if !rgxPrice.MatchString(value) {
		return false, "Price is invalid: " + value
	}
	return true, ""
}

func (v *Validator) Quantity(value int) (bool, string) {
	if value < 0 {
		return false, "Quantity is invalid: " + strconv.Itoa(value)
	}
	return true, ""
}

func (v *Validator) RegIdentifier(value string) (bool, string) {
	length := utf8.RuneCountInString(value)
	if length < regIdentifierMinLength {
		return false, "Regulatory identifier must be at least " + strconv.Itoa(regIdentifierMinLength) + " characters"
	}
	if length > regIdentifierMaxLength {
		return false, "Regulatory identifier must not exceed " + strconv.Itoa(regIdentifierMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) Sex(value string) (bool, string) {
	if value != "male" && value != "female" {
		return false, "Sex is invalid: " + utils.ToTitlecase(value)
	}
	return true, ""
}

func (v *Validator) URL(value string) (bool, string) {
	why := "URL is invalid: " + value
	u, err := url.ParseRequestURI(value)
	if err != nil {
		return false, why
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false, why
	}
	if u.Hostname() == "" {
		return false, why
	}
	return true, ""
}

func (v *Validator) UserAgent(value string) (bool, string) {
	if utf8.RuneCountInString(value) > userAgentMaxLength {
		return false, "User agent must not exceed " + strconv.Itoa(userAgentMaxLength) + " characters"
	}
	return true, ""
}

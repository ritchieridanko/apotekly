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

func (v *Validator) AddrRecipient(value string) (bool, string) {
	length := utf8.RuneCountInString(value)
	if length < nameMinLength {
		return false, "Recipient name must be at least " + strconv.Itoa(nameMinLength) + " characters"
	}
	if length > nameMaxLength {
		return false, "Recipient name must not exceed " + strconv.Itoa(nameMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) AddrStreet(value string) (bool, string) {
	length := utf8.RuneCountInString(value)
	if length < streetMinLength {
		return false, "Street must be at least " + strconv.Itoa(streetMinLength) + " characters"
	}
	if length > streetMaxLength {
		return false, "Street must not exceed " + strconv.Itoa(streetMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) AddrSubdivision(value string) (bool, string) {
	if utf8.RuneCountInString(value) > subdivisionMaxLength {
		return false, "Subdivision must not exceed " + strconv.Itoa(subdivisionMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) Birthdate(value time.Time) (bool, string) {
	if value.After(time.Now().UTC()) {
		return false, "Birthdate is invalid: " + value.Format("2 Jan 2006")
	}
	return true, ""
}

func (v *Validator) Country(value string) (bool, string) {
	if _, ok := countries[value]; !ok {
		return false, "Country is invalid: " + value
	}
	return true, ""
}

func (v *Validator) Description(value string) (bool, string) {
	if utf8.RuneCountInString(value) > descMaxLength {
		return false, "Description must not exceed " + strconv.Itoa(descMaxLength) + " characters"
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
		return false, "IP Address is invalid: " + value
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

func (v *Validator) Name(value string) (bool, string) {
	length := utf8.RuneCountInString(value)
	if length < nameMinLength {
		return false, "Name must be at least " + strconv.Itoa(nameMinLength) + " characters"
	}
	if length > nameMaxLength {
		return false, "Name must not exceed " + strconv.Itoa(nameMaxLength) + " characters"
	}
	return true, ""
}

func (v *Validator) OnlineHours(value json.RawMessage) (bool, string) {
	data, err := utils.FromJSONRawMessage[map[string][]string](value)
	if err != nil {
		return false, "Online Hours is invalid"
	}
	for day, hours := range *data {
		if _, ok := days[day]; !ok {
			return false, "Day in Online Hours is invalid: " + day
		}
		for _, hour := range hours {
			why := "Range in Online Hours is invalid: " + hour
			if !rgxTimeRange.MatchString(hour) {
				return false, why
			}

			parts := strings.SplitN(hour, "-", 2)
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

func (v *Validator) Sex(value string) (bool, string) {
	if value != "male" && value != "female" {
		return false, "Sex is invalid: " + value
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
		return false, "User Agent must not exceed " + strconv.Itoa(userAgentMaxLength) + " characters"
	}
	return true, ""
}

package validator

import (
	"net"
	"strconv"
	"time"
	"unicode/utf8"
)

type Validator struct{}

func Init() *Validator {
	return &Validator{}
}

func (v *Validator) Birthdate(value time.Time) (bool, string) {
	if value.After(time.Now().UTC()) {
		return false, "Birthdate is invalid: " + value.Format("2 Jan 2006")
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

func (v *Validator) Name(value string) (bool, string) {
	length := len(value)
	if length < nameMinLength {
		return false, "Name must be at least " + strconv.Itoa(nameMinLength) + " characters"
	}
	if length > nameMaxLength {
		return false, "Name must not exceed " + strconv.Itoa(nameMaxLength) + " characters"
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

func (v *Validator) Sex(value string) (bool, string) {
	if value != "male" && value != "female" {
		return false, "Sex is invalid: " + value
	}
	return true, ""
}

func (v *Validator) UserAgent(value string) (bool, string) {
	if len(value) > userAgentMaxLength {
		return false, "User Agent must not exceed " + strconv.Itoa(userAgentMaxLength) + " characters"
	}
	return true, ""
}

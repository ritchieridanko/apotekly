package validator

import (
	"net"
	"strconv"
)

type Validator struct{}

func Init() *Validator {
	return &Validator{}
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

func (v *Validator) Password(value string) (bool, string) {
	length := len(value)
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
		return false, "Password must include at least one special character: " + specialChars
	}
	return true, ""
}

func (v *Validator) UserAgent(value string) (bool, string) {
	if len(value) > userAgentMaxLength {
		return false, "User Agent must not exceed " + strconv.Itoa(userAgentMaxLength) + " characters"
	}
	return true, ""
}

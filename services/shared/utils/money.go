package utils

import (
	"fmt"
	"strconv"
	"strings"
)

var ZeroDecimalCurrencies map[string]struct{} = map[string]struct{}{
	"idr": {},
	"jpy": {},
	"khr": {},
	"krw": {},
	"lak": {},
	"mmk": {},
	"vnd": {},
}

func ParseMoney(money, currency string) (int64, error) {
	parts := strings.Split(money, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid money format: %s", money)
	}

	value, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}
	if _, exists := ZeroDecimalCurrencies[currency]; exists {
		return value, nil
	}

	cents := int64(0)
	if len(parts) == 2 {
		fraction := parts[1]
		if len(fraction) > 2 {
			return 0, fmt.Errorf("invalid money format: %s", money)
		}
		for len(fraction) < 2 {
			fraction += "0"
		}

		cents, err = strconv.ParseInt(fraction, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	if value < 0 {
		return (value * 100) - cents, nil
	}
	return (value * 100) + cents, nil
}

func StringifyMoney(amount int64, currency string) string {
	if _, exists := ZeroDecimalCurrencies[currency]; exists {
		return strconv.FormatInt(amount, 10)
	}
	return fmt.Sprintf("%d.%02d", amount/100, amount%100)
}

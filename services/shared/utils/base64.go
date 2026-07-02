package utils

import "encoding/base64"

// Encode value to MIME Base64
func ToMIMEBase64(value string) string {
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(value)) + "?="
}

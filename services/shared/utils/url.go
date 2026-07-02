package utils

import "net/url"

// Create a new tokenized URL
func GenerateTokenizedURL(baseURL, path, token string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	u.Path = path
	q := u.Query()
	q.Set("token", token)

	u.RawQuery = q.Encode()
	return u.String(), nil
}

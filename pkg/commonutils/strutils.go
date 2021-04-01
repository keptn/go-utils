package commonutils

import "net/url"

// IsValidURL checks whether the given string is a valid URL or not
func IsValidURL(strURL string) bool {
	_, err := url.ParseRequestURI(strURL)
	if err != nil {
		return false
	}
	u, err := url.Parse(strURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

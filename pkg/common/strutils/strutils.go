package strutils

import "errors"

// Sringp returns a string pointer of the given string
func Stringp(str string) *string {
	return &str
}

func AllSet(vals ...string) error {
	for _, val := range vals {
		if val == "" {
			return errors.New("empty value")
		}
	}
	return nil
}

package strutils

// Sringp returns a string pointer of the given string
func Stringp(str string) *string {
	return &str
}

// AllSet checks whether all provided string values are non-empty
func AllSet(vals ...string) bool {
	for _, val := range vals {
		if val == "" {
			return false
		}
	}
	return true
}

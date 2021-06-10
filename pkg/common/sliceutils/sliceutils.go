package sliceutils

// ContainsStr checks if a string str is present in a slice
func ContainsStr(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

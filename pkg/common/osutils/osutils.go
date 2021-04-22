package osutils

import "os"

// GetOSEnv retrieves the value of the environment variable named by the key.
// If the environment variable is not present the default value will be returned
func GetOSEnvOrDefault(key, defaultVal string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return defaultVal
}

// GetOSEnv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
func GetOSEnv(key string) string {
	return os.Getenv(key)
}

// GetAndCompareOSEnv retrieves the value of the environment variable named by the key
// and compares it to the value of compareStr. If the environment variable is not present
// it returns false
func GetAndCompareOSEnv(key, compareStr string) bool {
	v := GetOSEnv(key)
	if v == "" {
		return false
	}
	return v == compareStr
}

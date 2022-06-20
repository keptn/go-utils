package osutils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GetOSEnvOrDefault retrieves the value of the environment variable named by the key.
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

// ExecuteCommand exectues the command using the args
func ExecuteCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("Error executing command %s %s: %s\n%s", command, strings.Join(args, " "), err.Error(), string(out))
	}
	return string(out), nil
}

// ExecuteCommandInDirectory executes the command using the args within the specified directory
func ExecuteCommandInDirectory(command string, args []string, directory string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = directory
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("Error executing command %s %s: %s\n%s", command, strings.Join(args, " "), err.Error(), string(out))
	}
	return string(out), nil
}

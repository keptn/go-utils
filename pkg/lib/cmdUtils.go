package keptn

import (
	"fmt"
	"os/exec"
	"strings"
)

// ExecuteCommand exectues the command using the args
func ExecuteCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error executing command %s %s: %s\n%s", command, strings.Join(args, " "), err.Error(), string(out))
	}
	return string(out), nil
}

// ExecuteCommandInDirectory executes the command using the args within the specified directory
func ExecuteCommandInDirectory(command string, args []string, directory string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = directory
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error executing command %s %s: %s", command, strings.Join(args, " "), err.Error())
	}
	return string(out), nil
}

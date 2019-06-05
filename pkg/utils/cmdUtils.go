package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// ExecuteCommand exectues the command using the args
func ExecuteCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error executing command %s %s: %s", command, strings.Join(args, " "), err.Error())
	}
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error executing command %s %s: %s", command, strings.Join(args, " "), err.Error())
	}
	return string(out), nil
}

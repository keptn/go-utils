package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecuteCommand exectues the command using the args
func ExecuteCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error executing command %s %s: %s", command, strings.Join(args, " "), err.Error())
	}
	return nil
}

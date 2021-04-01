package fileutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ReadFile reads a file with the given file name and returns its content
// as a slice of bytes. This function also supports file paths which are pointing to the
// user's home directory, i.e. starting with ~/
func ReadFile(fileName string) ([]byte, error) {
	fileName = ExpandTilde(fileName)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("Cannot find file %s", fileName)
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ReadFileAsStr does the same thing as ReadFile but ruturns the content
// of the file as a string
func ReadFileAsStr(fileName string) (string, error) {
	content, err := ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), err
}

// FileExists checks if the file with the gieven name exists.
// This function also supports file paths which are pointing to the
// user's home directory, i.e. starting with ~/
func FileExists(fileName string) bool {
	info, err := os.Stat(ExpandTilde(fileName))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ExpandTilde expands the tilde symbol (~) to the user's home directory
// This function is also useful for expanding file paths like "~/a/b/c"
func ExpandTilde(fileName string) string {
	if fileName == "~" {
		return UserHomeDir()
	} else if strings.HasPrefix(fileName, "~/") {
		return filepath.Join(UserHomeDir(), fileName[2:])
	}
	return fileName
}

// UserHomeDir returns the user's home directory
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

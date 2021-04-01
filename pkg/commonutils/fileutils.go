package commonutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

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

func ReadFileAsStr(fileName string) (string, error) {
	content, err := ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), err
}

func FileExists(filename string) bool {
	info, err := os.Stat(ExpandTilde(filename))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ExpandTilde(fileName string) string {
	if fileName == "~" {
		return UserHomeDir()
	} else if strings.HasPrefix(fileName, "~/") {
		return filepath.Join(UserHomeDir(), fileName[2:])
	}
	return fileName
}

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

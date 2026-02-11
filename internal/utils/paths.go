package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/ChristianLapinig/aem-local-cli/constants"
)

func PathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func PathExistsWithError(path string) error {
	if !PathExists(path) {
		msg := fmt.Sprintf("%s %s\n", constants.PathDoesNotExist, path)
		return errors.New(msg)
	}
	return nil
}

func GetHomePath() string {
	if testPath := os.Getenv("AEMLOCAL_TEST_HOME"); testPath != "" {
		return testPath
	}
	home, _ := os.UserHomeDir()
	return home
}

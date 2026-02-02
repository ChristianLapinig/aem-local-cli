package utils

import (
	"os"
)

func ErrorAndCleanup(path string, err error) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return err
}

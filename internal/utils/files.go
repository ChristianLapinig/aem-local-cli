package utils

import (
	"io"
	"os"
	"path/filepath"
)

func LoadMarkerFile() ([]byte, error) {
	home := GetHomePath()
	markerFile := filepath.Join(home, ".aemlocal_path")
	if err := PathExistsWithError(markerFile); err != nil {
		return []byte{}, err
	}
	path, err := os.ReadFile(markerFile)
	if err != nil {
		return []byte{}, err
	}
	return path, nil
}

func CopyFile(src, dest string) error {
	if err := PathExistsWithError(src); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	if err := destFile.Sync(); err != nil {
		return err
	}

	return nil
}

package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/constants"
	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
)

func TestCopyFile(t *testing.T) {
	tmp := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmp, "src"), 0755); err != nil {
		t.Fatalf("Error creating src folder: %v", err)
	}

	if err := os.Mkdir(filepath.Join(tmp, "dest"), 0755); err != nil {
		t.Fatalf("Error creating dest folder: %v", err)
	}

	src := filepath.Join(tmp, "src", "test.txt")
	dest := filepath.Join(tmp, "dest", "test.txt")
	srcFile, err := os.Create(src)
	if err != nil {
		t.Fatalf("Error creating %s: %v", src, err)
	}
	defer srcFile.Close()

	if err := utils.CopyFile(src, dest); err != nil {
		t.Fatalf("Error copying file from %s to %s: %v", src, dest, err)
	}

	if !utils.PathExists(dest) {
		t.Errorf("FAILED: expected %s to exist.", dest)
	}
}

func TestCopyFile_SrcPath_Doesnt_Exist(t *testing.T) {
	tmp := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmp, "dest"), 0755); err != nil {
		t.Fatalf("Error creating dest folder: %v", err)
	}

	src := filepath.Join(tmp, "src", "test.txt")
	dest := filepath.Join(tmp, "dest", "test.txt")
	err := utils.CopyFile(src, dest)

	if err != nil && !strings.Contains(err.Error(), constants.PathDoesNotExist) {
		t.Errorf("FAILED: got error %s, expected error to contain %s", err.Error(), constants.PathDoesNotExist)
	}
}

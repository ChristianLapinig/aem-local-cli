package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/cmd"
	"github.com/spf13/cobra"
)

type data struct {
	Path string
	File *os.File
}

// Setup test files
func SetupFile(t testing.TB, name, path string) *data {
	dest := filepath.Join(path, name)
	f, err := os.Create(dest)
	if err != nil {
		t.Fatalf("error creating file %s", dest)
	}
	return &data{
		Path: dest,
		File: f,
	}
}

func SetupTempDir(t testing.TB) string {
	tmp := t.TempDir()
	t.Setenv("AEMLOCAL_TEST_HOME", tmp)
	return tmp
}

func SetupWithRootCmd(t testing.TB) (*cobra.Command, string) {
	tmp := SetupTempDir(t)
	rootCmd := cmd.NewRootCmd()
	return rootCmd, tmp
}

func SetupWithInitCmd(t testing.TB) (*cobra.Command, string) {
	rootCmd, tmp := SetupWithRootCmd(t)
	rootCmd.AddCommand(cmd.NewInitCmd())
	return rootCmd, tmp
}

func SetupForSubcommands(t testing.TB) (*cobra.Command, string) {
	rootCmd, tmp := SetupWithRootCmd(t)
	rootCmd.AddCommand(cmd.NewInitCmd())
	rootCmd.SetArgs([]string{"init", "-p", tmp})
	return rootCmd, tmp
}

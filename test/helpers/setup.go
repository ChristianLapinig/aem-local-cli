package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/cmd"
	"github.com/spf13/cobra"
)

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
	rootCmd, tmp := SetupWithInitCmd(t)
	envsPath := filepath.Join(tmp, "envs")
	if err := os.Mkdir(envsPath, 0o755); err != nil {
		t.Fatalf("Error creating folder %s: %v", envsPath, err)
	}
	rootCmd.SetArgs([]string{"init", "-p", tmp, "-e", envsPath})
	return rootCmd, tmp
}

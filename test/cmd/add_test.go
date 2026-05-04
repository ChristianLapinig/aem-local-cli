package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/cmd"
	"github.com/ChristianLapinig/aem-local-cli/constants"
	"github.com/ChristianLapinig/aem-local-cli/models/config"
	"github.com/ChristianLapinig/aem-local-cli/test/helpers"
)

func TestAddCommand_Success(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	envPath := filepath.Join(tmp, "my-env")
	if err := os.Mkdir(envPath, 0o755); err != nil {
		t.Fatalf("error creating environment directory: %v", err)
	}

	addCmd := cmd.NewAddCommand()
	addCmd.SetArgs([]string{"my-env", envPath})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("error executing add command: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}
	if len(cfg.Environments) == 0 {
		t.Fatal("FAILED: expected environment to be added to config")
	}

	env := cfg.Environments[0]
	if env.Name != "my-env" {
		t.Errorf("FAILED: expected environment name to be 'my-env', got %s", env.Name)
	}
	if env.Path != envPath {
		t.Errorf("FAILED: expected environment path to be %s, got %s", envPath, env.Path)
	}
}

func TestAddCommand_Path_Doesnt_Exist(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	addCmd := cmd.NewAddCommand()
	addCmd.SetArgs([]string{"my-env", filepath.Join(tmp, "nonexistent")})

	err := addCmd.Execute()
	if err == nil {
		t.Fatal("FAILED: expected error when path does not exist")
	}
	if !strings.Contains(err.Error(), constants.PathDoesNotExist) {
		t.Errorf("FAILED: expected '%s' error, got: %v", constants.PathDoesNotExist, err)
	}
}

func TestAddCommand_Not_Enough_Args(t *testing.T) {
	rootCmd, _ := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	addCmd := cmd.NewAddCommand()
	addCmd.SetArgs([]string{})

	err := addCmd.Execute()
	if err == nil {
		t.Fatal("FAILED: expected error when no args provided")
	}
	if !strings.Contains(err.Error(), "accepts 2 arg") {
		t.Errorf("FAILED: expected not enough args error, got: %v", err)
	}
}

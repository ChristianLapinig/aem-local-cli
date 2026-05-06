package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/cmd"
	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
	"github.com/ChristianLapinig/aem-local-cli/models/config"
	"github.com/ChristianLapinig/aem-local-cli/models/environment"
	"github.com/ChristianLapinig/aem-local-cli/test/helpers"
)

func addEnvToConfig(t *testing.T, name, path string) {
	t.Helper()
	configPath, err := config.GetConfigPath()
	if err != nil {
		t.Fatalf("error getting config path: %v", err)
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}
	cfg.Environments = append(cfg.Environments, environment.Environment{Name: name, Path: path})
	if err := config.UpdateConfig(configPath, cfg); err != nil {
		t.Fatalf("error updating config: %v", err)
	}
}

func TestDeleteCommand_NoEnvironments(t *testing.T) {
	rootCmd, _ := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	deleteCmd := cmd.NewDeleteCmd()
	deleteCmd.SetArgs([]string{})
	if err := deleteCmd.Execute(); err != nil {
		t.Fatalf("FAILED: expected no error when no environments configured, got: %v", err)
	}
}

func TestDeleteCommand_EnvironmentNotFound(t *testing.T) {
	rootCmd, _ := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	addEnvToConfig(t, "other-env", "/some/path")

	deleteCmd := cmd.NewDeleteCmd()
	deleteCmd.SetArgs([]string{"--name", "nonexistent"})
	err := deleteCmd.Execute()
	if err == nil {
		t.Fatal("FAILED: expected error when environment not found")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("FAILED: expected error to mention environment name, got: %v", err)
	}
}

func TestDeleteCommand_Abort(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	envPath := filepath.Join(tmp, "my-env")
	if err := os.Mkdir(envPath, 0o755); err != nil {
		t.Fatalf("error creating environment directory: %v", err)
	}
	addEnvToConfig(t, "my-env", envPath)

	deleteCmd := cmd.NewDeleteCmd()
	deleteCmd.SetArgs([]string{"--name", "my-env"})
	deleteCmd.SetIn(strings.NewReader("n\n"))
	if err := deleteCmd.Execute(); err != nil {
		t.Fatalf("FAILED: expected no error on abort, got: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}
	if len(cfg.Environments) != 1 {
		t.Errorf("FAILED: expected environment to remain in config after abort, got %d environments", len(cfg.Environments))
	}
}

func TestDeleteCommand_Success_ConfigOnly(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	envPath := filepath.Join(tmp, "my-env")
	if err := os.Mkdir(envPath, 0o755); err != nil {
		t.Fatalf("error creating environment directory: %v", err)
	}
	addEnvToConfig(t, "my-env", envPath)

	// "y" to confirm deletion, "n" to skip folder removal
	deleteCmd := cmd.NewDeleteCmd()
	deleteCmd.SetArgs([]string{"--name", "my-env"})
	deleteCmd.SetIn(strings.NewReader("y\nn\n"))
	if err := deleteCmd.Execute(); err != nil {
		t.Fatalf("FAILED: expected no error, got: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}
	if len(cfg.Environments) != 0 {
		t.Errorf("FAILED: expected environment to be removed from config, got %d environments", len(cfg.Environments))
	}
	if !utils.PathExists(envPath) {
		t.Errorf("FAILED: expected environment folder to still exist at %s", envPath)
	}
}

func TestDeleteCommand_Success_WithPurgeFlag(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	envPath := filepath.Join(tmp, "my-env")
	if err := os.Mkdir(envPath, 0o755); err != nil {
		t.Fatalf("error creating environment directory: %v", err)
	}
	addEnvToConfig(t, "my-env", envPath)

	// "y" to confirm deletion; --purge flag skips the folder prompt
	deleteCmd := cmd.NewDeleteCmd()
	deleteCmd.SetArgs([]string{"--name", "my-env", "--purge"})
	deleteCmd.SetIn(strings.NewReader("y\n"))
	if err := deleteCmd.Execute(); err != nil {
		t.Fatalf("FAILED: expected no error, got: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}
	if len(cfg.Environments) != 0 {
		t.Errorf("FAILED: expected environment to be removed from config, got %d environments", len(cfg.Environments))
	}
	if utils.PathExists(envPath) {
		t.Errorf("FAILED: expected environment folder to be deleted at %s", envPath)
	}
}

func TestDeleteCommand_Purge_EmptyPath(t *testing.T) {
	rootCmd, _ := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	addEnvToConfig(t, "my-env", "")

	deleteCmd := cmd.NewDeleteCmd()
	deleteCmd.SetArgs([]string{"--name", "my-env", "--purge"})
	deleteCmd.SetIn(strings.NewReader("y\n"))
	err := deleteCmd.Execute()
	if err == nil {
		t.Fatal("FAILED: expected error when purging environment with no path configured")
	}
	if !strings.Contains(err.Error(), "my-env") {
		t.Errorf("FAILED: expected error to mention environment name, got: %v", err)
	}
}

func TestDeleteCommand_Success_PurgeViaPrompt(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	envPath := filepath.Join(tmp, "my-env")
	if err := os.Mkdir(envPath, 0o755); err != nil {
		t.Fatalf("error creating environment directory: %v", err)
	}
	addEnvToConfig(t, "my-env", envPath)

	// "y" to confirm deletion, "y" to also delete folder
	deleteCmd := cmd.NewDeleteCmd()
	deleteCmd.SetArgs([]string{"--name", "my-env"})
	deleteCmd.SetIn(strings.NewReader("y\ny\n"))
	if err := deleteCmd.Execute(); err != nil {
		t.Fatalf("FAILED: expected no error, got: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}
	if len(cfg.Environments) != 0 {
		t.Errorf("FAILED: expected environment to be removed from config, got %d environments", len(cfg.Environments))
	}
	if utils.PathExists(envPath) {
		t.Errorf("FAILED: expected environment folder to be deleted at %s", envPath)
	}
}

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/constants"
	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
	"github.com/ChristianLapinig/aem-local-cli/test/helpers"
)

func TestInitCommand_Default_Options(t *testing.T) {
	rootCmd, _ := helpers.SetupWithInitCmd(t)
	rootCmd.SetArgs([]string{"init"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing inti command: %v", err)
	}

	home := utils.GetHomePath()
	configPath := filepath.Join(home, constants.AemLocalFolder)
	if !utils.PathExists(configPath) {
		t.Errorf("FAILED: expected .aemlocal folder to be under %s", home)
	}
	if !utils.PathExists(filepath.Join(home, constants.MarkerFile)) {
		t.Errorf("FAILED: expected .aemlocal_path folder to be under %s", home)
	}
	if !utils.PathExists(filepath.Join(configPath, "config.json")) {
		t.Errorf("FAILED: expected config.json to be under %s", configPath)
	}
}

func TestInitCommand_With_Path_Flag(t *testing.T) {
	rootCmd, tmp := helpers.SetupWithInitCmd(t)
	rootCmd.SetArgs([]string{"init", "-p", tmp})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Error executing init command: %v", err)
	}

	configPath := filepath.Join(tmp, constants.AemLocalFolder)
	if !utils.PathExists(configPath) {
		t.Errorf("FAILED: expected folder %s to exist", configPath)
	}

	if !utils.PathExists(filepath.Join(configPath, "temp")) {
		t.Errorf("FAILED: expected folder %s/temp to exist", configPath)
	}
	if !utils.PathExists(filepath.Join(configPath, "config.json")) {
		t.Errorf("FAILED: expected config.json to be under %s", configPath)
	}

	helpers.Teardown(t)
}

func TestInitCommand_Existing_Setup_Overwrite_Yes(t *testing.T) {
	rootCmd, tmp := helpers.SetupWithInitCmd(t)

	// Pre-create .aemlocal with a sentinel file to prove it gets replaced
	existingConfig := filepath.Join(tmp, constants.AemLocalFolder)
	if err := os.Mkdir(existingConfig, 0o755); err != nil {
		t.Fatalf("failed to create existing .aemlocal: %v", err)
	}
	sentinelPath := filepath.Join(existingConfig, "old_file.txt")
	if err := os.WriteFile(sentinelPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("failed to create sentinel file: %v", err)
	}

	rootCmd.SetIn(strings.NewReader("y\n"))
	rootCmd.SetArgs([]string{"init", "-p", tmp})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing init command: %v", err)
	}

	if utils.PathExists(sentinelPath) {
		t.Error("FAILED: expected old .aemlocal to be replaced, but sentinel file still exists")
	}
	if !utils.PathExists(filepath.Join(existingConfig, "config.json")) {
		t.Errorf("FAILED: expected config.json to exist after overwrite at %s", existingConfig)
	}

	helpers.Teardown(t)
}

func TestInitCommand_Existing_Setup_Overwrite_No(t *testing.T) {
	rootCmd, tmp := helpers.SetupWithInitCmd(t)

	existingConfig := filepath.Join(tmp, constants.AemLocalFolder)
	if err := os.Mkdir(existingConfig, 0o755); err != nil {
		t.Fatalf("failed to create existing .aemlocal: %v", err)
	}
	sentinelPath := filepath.Join(existingConfig, "old_file.txt")
	if err := os.WriteFile(sentinelPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("failed to create sentinel file: %v", err)
	}

	rootCmd.SetIn(strings.NewReader("n\n"))
	rootCmd.SetArgs([]string{"init", "-p", tmp})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing init command: %v", err)
	}

	if !utils.PathExists(sentinelPath) {
		t.Error("FAILED: expected existing .aemlocal to be preserved when user declines overwrite")
	}

	helpers.Teardown(t)
}

func TestInitCommand_With_Path_Flag_Non_Existent(t *testing.T) {
	rootCmd, _ := helpers.SetupWithInitCmd(t)
	rootCmd.SetArgs([]string{"init", "-p", "./temp"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("FAILED: Expected error to be thrown for non-existent path.")
	}

	if !strings.Contains(err.Error(), "Path does not exist") {
		t.Error("FAILED: Expected path non-existent error to be thrown")
	}

	helpers.Teardown(t)
}

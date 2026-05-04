package cmd

import (
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

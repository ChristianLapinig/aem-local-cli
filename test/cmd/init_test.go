package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/cmd"
	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
	"github.com/ChristianLapinig/aem-local-cli/models"
)

func TestInitCommand_With_Path_Flag(t *testing.T) {
	tmp := t.TempDir()
	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(cmd.NewInitCmd())
	rootCmd.SetArgs([]string{"init", "-p", tmp})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Error executing init command: %v", err)
	}

	configPath := filepath.Join(tmp, ".aemlocal")
	if !utils.PathExists(configPath) {
		t.Errorf("FAILED: expected folder %s to exist", configPath)
	}

	if !utils.PathExists(filepath.Join(configPath, "temp")) {
		t.Errorf("FAILED: expected folder %s/temp to exist", configPath)
	}

	if !utils.PathExists(filepath.Join(configPath, "config.json")) {
		t.Errorf("FAILED: expected %s/config.json to exist", configPath)
	}
}

func TestInitCommand_With_EnvsPath_Flag(t *testing.T) {
	tmp := t.TempDir()
	envsPath := filepath.Join(tmp, "envs")
	// Represents where local AEM environments are stored
	if err := os.Mkdir(envsPath, 0755); err != nil {
		t.Fatalf("Error creating folder %s: %v", envsPath, err)
	}
	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(cmd.NewInitCmd())
	rootCmd.SetArgs([]string{"init", "-p", tmp, "-e", envsPath})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Error executing init command: %v", err)
	}

	jsonPath := filepath.Join(tmp, ".aemlocal", "config.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("Error opening config.json file: %v", err)
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Error unmarshalling json data: %v", err)
	}
	if config.EnvsPath != envsPath {
		t.Errorf("FAILED: expected envsPath to be %s. Got %s", envsPath, config.EnvsPath)
	}
}

func TestInitCommand_With_Path_Flag_Non_Existent(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(cmd.NewInitCmd())
	rootCmd.SetArgs([]string{"init", "-p", "./temp"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("FAILED: Expected error to be thrown for non-existent path.")
	}

	if !strings.Contains(err.Error(), "Path does not exist") {
		t.Error("FAILED: Expected path non-existent error to be thrown")
	}
}

func TestInitCommand_With_EnvsPath_Flag_Non_Existent(t *testing.T) {
	tmp := t.TempDir()
	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(cmd.NewInitCmd())
	rootCmd.SetArgs([]string{"init", "-p", tmp, "-e", "does-not-exist"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("FAILED: Expected error to be thrown for non-existent path.")
	}

	if !strings.Contains(err.Error(), "Environments path does not exist") {
		t.Error("FAILED: Expected environments path non-existent error to be thrown")
	}
}

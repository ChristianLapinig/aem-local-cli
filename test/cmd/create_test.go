package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/cmd"
	"github.com/ChristianLapinig/aem-local-cli/constants"
	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
	"github.com/ChristianLapinig/aem-local-cli/models/config"
	"github.com/ChristianLapinig/aem-local-cli/test/helpers"
)

func authorAndPublishExist(t testing.TB, path string) {
	authorPath := filepath.Join(path, constants.Author)
	publishPath := filepath.Join(path, constants.Publish)
	if !utils.PathExists(authorPath) {
		t.Errorf("FAILED: expected folder to exist at %s", authorPath)
	}
	if !utils.PathExists(publishPath) {
		t.Errorf("FAILED: expected folder to exist at %s", publishPath)
	}
}

func TestCreateCommand_Default_Options(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	licenseProps := helpers.SetupFile(t, constants.LicenseProperties, tmp)
	quickstartJar := helpers.SetupFile(t, "cq-quickstart.jar", tmp)
	defer licenseProps.File.Close()
	defer quickstartJar.File.Close()

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{licenseProps.Path, quickstartJar.Path, "-p", tmp, "-n", "my-env"})
	if err := createCmd.Execute(); err != nil {
		t.Fatalf("error executing command: %v", err)
	}

	envDir := filepath.Join(tmp, "my-env")
	authorAndPublishExist(t, envDir)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config file: %v", err)
	}
	if len(cfg.Environments) == 0 {
		t.Error("FAILED: expected new environment to have been added to config.json")
	}
	if cfg.Environments[0].Path != envDir {
		t.Errorf("FAILED: expected environment path to be %s, got %s", envDir, cfg.Environments[0].Path)
	}
}

func TestCreateCommand_With_Options(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	licenseProps := helpers.SetupFile(t, constants.LicenseProperties, tmp)
	quickstartJar := helpers.SetupFile(t, "cq-quickstart.jar", tmp)
	defer licenseProps.File.Close()
	defer quickstartJar.File.Close()

	name := "test"
	baseDir := filepath.Join(tmp, "cloud-service")
	if err := os.Mkdir(baseDir, 0o755); err != nil {
		t.Fatalf("error creating base directory: %v", err)
	}
	authorPort := "8080"
	publishPort := "8081"

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{
		licenseProps.Path,
		quickstartJar.Path,
		"-p", baseDir,
		"-n", name,
		"--author-port", authorPort,
		"--publish-port", publishPort,
	})

	if err := createCmd.Execute(); err != nil {
		t.Fatalf("error executing command: %v", err)
	}

	envPath := filepath.Join(baseDir, name)
	authorAndPublishExist(t, envPath)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config file: %v", err)
	}
	if len(cfg.Environments) == 0 {
		t.Error("FAILED: expected new environment to have been added to config.json")
	}

	environment := cfg.Environments[0]
	if environment.Name != name {
		t.Errorf("FAILED: expected environment name to be %s, got %s", name, environment.Name)
	}
	if environment.Path != envPath {
		t.Errorf("FAILED: expected environment path to be %s, got %s", envPath, environment.Path)
	}
}

func TestCreateCommand_Not_Enough_Args(t *testing.T) {
	rootCmd, _ := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{})

	err := createCmd.Execute()
	if err == nil {
		t.Errorf("FAILED: expected error to be thrown")
	}
	if err != nil && !strings.Contains(err.Error(), "accepts 2 arg") {
		t.Errorf("FAILED: expected not enough args error to be thrown: %v", err)
	}
}

func TestCreateCommand_LicenseProps_Doesnt_Exist(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	quickstartJar := helpers.SetupFile(t, "cq-quickstart.jar", tmp)
	defer quickstartJar.File.Close()

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{
		filepath.Join(tmp, constants.LicenseProperties),
		quickstartJar.Path,
	})
	err := createCmd.Execute()
	if err == nil {
		t.Error("FAILED: expected error to be thrown")
	}

	if err != nil && !strings.Contains(err.Error(), constants.PathDoesNotExist) {
		t.Errorf("FAILED: expected %s error to be thrown.", constants.PathDoesNotExist)
	}
}

func TestCreateCommand_Quickstart_JAR_Doesnt_Exist(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	licenseProps := helpers.SetupFile(t, constants.LicenseProperties, tmp)
	defer licenseProps.File.Close()

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{
		licenseProps.Path,
		filepath.Join(tmp, "cq-quickstart.jar"),
	})

	err := createCmd.Execute()
	if err == nil {
		t.Error("FAILED: expected error to be thrown")
	}

	if err != nil && !strings.Contains(err.Error(), constants.PathDoesNotExist) {
		t.Errorf("FAILED: expected %s error to be thrown.", constants.PathDoesNotExist)
	}
}

func TestCreateCommand_NamedDest_AlreadyExists(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	licenseProps := helpers.SetupFile(t, constants.LicenseProperties, tmp)
	quickstartJar := helpers.SetupFile(t, "cq-quickstart.jar", tmp)
	defer licenseProps.File.Close()
	defer quickstartJar.File.Close()

	if err := os.Mkdir(filepath.Join(tmp, "existing-env"), 0o755); err != nil {
		t.Fatalf("failed to pre-create destination: %v", err)
	}

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{licenseProps.Path, quickstartJar.Path, "-p", tmp, "-n", "existing-env"})
	err := createCmd.Execute()
	if err == nil {
		t.Fatal("FAILED: expected error when named destination already exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("FAILED: expected 'already exists' error, got: %v", err)
	}
}

func TestCreateCommand_DuplicateName(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	addEnvToConfig(t, "my-env", filepath.Join(tmp, "my-env"))

	licenseProps := helpers.SetupFile(t, constants.LicenseProperties, tmp)
	quickstartJar := helpers.SetupFile(t, "cq-quickstart.jar", tmp)
	defer licenseProps.File.Close()
	defer quickstartJar.File.Close()

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{licenseProps.Path, quickstartJar.Path, "-p", tmp, "-n", "my-env"})
	err := createCmd.Execute()
	if err == nil {
		t.Fatal("FAILED: expected error when environment name already exists in config")
	}
	if !strings.Contains(err.Error(), "my-env") {
		t.Errorf("FAILED: expected error to mention environment name, got: %v", err)
	}
}

func TestCreateCommand_Invalid_Port_Flag_Value(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	licenseProps := helpers.SetupFile(t, constants.LicenseProperties, tmp)
	quickstartJar := helpers.SetupFile(t, "cq-quickstart.jar", tmp)
	defer licenseProps.File.Close()
	defer quickstartJar.File.Close()

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{
		licenseProps.Path,
		quickstartJar.Path,
		"-n",
		"test",
		"--author-port",
		"abcd",
	})

	if err := createCmd.Execute(); err == nil {
		t.Errorf("FAILED: expected error to be thrown.")
	}
}

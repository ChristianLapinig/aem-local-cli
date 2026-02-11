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

type data struct {
	path string
	file *os.File
}

// Setup test files
func setupFile(t testing.TB, name, path string) *data {
	dest := filepath.Join(path, name)
	f, err := os.Create(dest)
	if err != nil {
		t.Fatalf("error creating file %s", dest)
	}
	return &data{
		path: dest,
		file: f,
	}
}

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

	licenseProps := setupFile(t, constants.LicenseProperties, tmp)
	quickstartJar := setupFile(t, "cq-quickstart.jar", tmp)
	defer licenseProps.file.Close()
	defer quickstartJar.file.Close()

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{licenseProps.path, quickstartJar.path})
	if err := createCmd.Execute(); err != nil {
		t.Fatalf("error executing command: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config file: %v", err)
	}

	authorAndPublishExist(t, filepath.Join(cfg.EnvsPath, "aem"))

	cfg, err = config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config file: %v", err)
	}
	if len(cfg.Environments) == 0 {
		t.Error("FAILED: expected new environment to have been added to config.json")
	}
}

func TestCreateCommand_With_Options(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	licenseProps := setupFile(t, constants.LicenseProperties, tmp)
	quickstartJar := setupFile(t, "cq-quickstart.jar", tmp)
	defer licenseProps.file.Close()
	defer quickstartJar.file.Close()

	name := "test"
	pathFlag := "cloud-service"
	authorPort := "8080"
	publishPort := "8081"
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config file: %v", err)
	}

	if err := os.Mkdir(filepath.Join(cfg.EnvsPath, pathFlag), 0o755); err != nil {
		t.Fatalf("error creating folder: %v", err)
	}

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{
		licenseProps.path,
		quickstartJar.path,
		"-p",
		pathFlag,
		"-n",
		name,
		"--author-port",
		authorPort,
		"--publish-port",
		publishPort,
	})

	if err := createCmd.Execute(); err != nil {
		t.Fatalf("error executing command: %v", err)
	}

	envPath := filepath.Join(cfg.EnvsPath, pathFlag, name)
	authorAndPublishExist(t, envPath)

	cfg, err = config.LoadConfig()
	if err != nil {
		t.Fatalf("error loading config file: %v", err)
	}
	if len(cfg.Environments) == 0 {
		t.Error("FAILED: expected new environment to have been added to config.json")
	}

	environment := cfg.Environments[0]
	if environment.Name != name {
		t.Errorf("FAILED: expected environment name to be %s, got %s", environment.Name, name)
	}

	if environment.Path != envPath {
		t.Errorf("FAILED: expected environment path to be %s, got %s", environment.Path, envPath)
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

	quickstartJar := setupFile(t, "cq-quickstart.jar", tmp)
	defer quickstartJar.file.Close()

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{
		filepath.Join(tmp, constants.LicenseProperties),
		quickstartJar.path,
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

	licenseProps := setupFile(t, constants.LicenseProperties, tmp)
	defer licenseProps.file.Close()

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{
		licenseProps.path,
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

func TestCreateCommand_Invalid_Port_Flag_Value(t *testing.T) {
	rootCmd, tmp := helpers.SetupForSubcommands(t)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("error executing root command: %v", err)
	}

	licenseProps := setupFile(t, constants.LicenseProperties, tmp)
	quickstartJar := setupFile(t, "cq-quickstart.jar", tmp)
	defer licenseProps.file.Close()
	defer quickstartJar.file.Close()

	createCmd := cmd.NewCreateCommand()
	createCmd.SetArgs([]string{
		licenseProps.path,
		quickstartJar.path,
		"-n",
		"test",
		"--author-port",
		"abcd",
	})

	if err := createCmd.Execute(); err == nil {
		t.Errorf("FAILED: expected error to be thrown.")
	}
}

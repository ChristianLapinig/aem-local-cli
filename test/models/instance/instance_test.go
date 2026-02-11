package instance

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/constants"
	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
	"github.com/ChristianLapinig/aem-local-cli/models/instance"
	"github.com/ChristianLapinig/aem-local-cli/models/paths"
	"github.com/ChristianLapinig/aem-local-cli/test/helpers"
)

func TestInstance_Create(t *testing.T) {
	tmp := helpers.SetupTempDir(t)
	licensePropsPath := filepath.Join(tmp, "license.properties")
	licenseProps, err := os.Create(licensePropsPath)
	if err != nil {
		t.Fatalf("error creating file %s: %v", licensePropsPath, err)
	}
	defer licenseProps.Close()

	quickstartPath := filepath.Join(tmp, "cq-quickstart.jar")
	quickstartJar, err := os.Create(quickstartPath)
	if err != nil {
		t.Fatalf("error creating file %s: %v", quickstartPath, err)
	}
	defer quickstartJar.Close()

	paths := &paths.Paths{
		Name:              filepath.Join(tmp, "test"),
		LicenseProperties: licensePropsPath,
		QuickstartJAR:     quickstartPath,
	}
	instance := &instance.Instance{
		Name: constants.Author,
		Port: constants.DefaultAuthorPort,
	}

	if err := instance.Create(paths); err != nil {
		t.Fatalf("error generating instance: %v", err)
	}

	aemPath := filepath.Join(paths.Name, instance.Name)
	if !utils.PathExists(aemPath) {
		t.Errorf("FAILED: expected folder at %s to be exist", aemPath)
	}

	if !utils.PathExists(filepath.Join(aemPath, constants.LicenseProperties)) {
		t.Errorf("FAILED: expected license.properties at %s to be exist", aemPath)
	}

	aemJar := fmt.Sprintf("aem-%s-p%d.jar", instance.Name, instance.Port)
	if !utils.PathExists(filepath.Join(aemPath, aemJar)) {
		t.Errorf("FAILED: expected %s at %s to be exist", aemJar, aemPath)
	}

	helpers.Teardown(t)
}

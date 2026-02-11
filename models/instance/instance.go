package instance

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ChristianLapinig/aem-local-cli/constants"
	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
	"github.com/ChristianLapinig/aem-local-cli/models/paths"
)

type Instance struct {
	Name string
	Port int
}

func NewInstance(name string, port int) *Instance {
	return &Instance{
		Name: name,
		Port: port,
	}
}

func (i *Instance) Create(paths *paths.Paths) error {
	if !utils.PathExists(paths.Name) {
		err := os.Mkdir(paths.Name, 0o755)
		if err != nil {
			return err
		}
	}

	instancePath := filepath.Join(paths.Name, i.Name)
	if err := os.Mkdir(instancePath, 0o755); err != nil {
		return err
	}

	licenseProperties := filepath.Join(instancePath, constants.LicenseProperties)
	if err := utils.CopyFile(paths.LicenseProperties, licenseProperties); err != nil {
		return err
	}

	quickstartJar := fmt.Sprintf("%s/aem-%s-p%d.jar", instancePath, i.Name, i.Port)
	if err := utils.CopyFile(paths.QuickstartJAR, quickstartJar); err != nil {
		return err
	}

	return nil
}

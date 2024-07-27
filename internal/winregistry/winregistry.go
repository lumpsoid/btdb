package winregistry

import (
	"btdb/internal/utility"
	"fmt"
	"os/exec"
	"path/filepath"
)

func Export(
	mountpoint string,
	registryPath string,
	prefix string,
	registryQuery string,
	outputPath string,
) error {
	// reged -x /mnt/Windows/System32/config/SYSTEM PREFIX "ControlSet001\Services\BTHPORT" ./out.reg
	registryPathFull := filepath.Join(mountpoint, registryPath)
	regExist := utility.FileExists(registryPathFull)
	if !regExist {
		return fmt.Errorf("registry file not found: %s", registryPathFull)
	}

	outputExist := utility.FileExists(outputPath)
	if outputExist {
		return fmt.Errorf("output file found: %s", outputPath)
	}

	cmd := exec.Command("reged", "-x", registryPathFull, prefix, registryQuery, outputPath)

	// Execute the command
	if err := cmd.Run(); err != nil {
		return err // Program not found
	}
	return nil
}

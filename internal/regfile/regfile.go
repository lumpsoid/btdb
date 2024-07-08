package regfile

import (
	"btdb/internal/device"
)

type Regfile struct {
	// Path to the registry file
	Filepath string
	// Windows mountpoint
	Mountpoint string
	// Registry path
	RegPath string
	// Prefix to the registry path
	Prefix string
	// Bluetooth devices
	Devices map[string]device.Bluetooth
	// Parsed registry data
	Data map[string]interface{}
}

package regfile

type Regfile struct {
	// Path to the registry file
	Filepath string
	// Windows mountpoint
	Mountpoint string
	// Registry path
	RegPath string
	// Prefix to the registry path
	Prefix string
	// Parsed registry data
	Data map[string]interface{}
}

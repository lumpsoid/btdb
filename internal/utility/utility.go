package utility

import (
	"os"
	"os/exec"
)

func CheckProgramAvailability(program string) error {
	// Run "which" command to check program availability
	cmd := exec.Command("which", program)

	// Execute the command
	if err := cmd.Run(); err != nil {
		return err // Program not found
	}

	return nil // Program found
}

func FileExists(filePath string) bool {
	// Use os.Stat to get file info
	_, err := os.Stat(filePath)

	// Check if there was no error (file exists)
	if err == nil {
		return true
	}

	// Check if the error indicates that the file doesn't exist
	if os.IsNotExist(err) {
		return false
	}

	// For any other error, assume the file exists (to be safe)
	return true
}

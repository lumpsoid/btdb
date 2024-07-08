package utility

import "os/exec"

func CheckProgramAvailability(program string) error {
	// Run "which" command to check program availability
	cmd := exec.Command("which", program)

	// Execute the command
	if err := cmd.Run(); err != nil {
		return err // Program not found
	}

	return nil // Program found
}

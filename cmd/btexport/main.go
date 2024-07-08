package main

import (
	"fmt"
	"os/exec"
)

func main() {
	remoteMachine := "RemoteMachineName"
	command := "PsExec64.exe"
	args := []string{"\\" + remoteMachine, "-accepteula", "-s", "-i", "reg", "query", "\"HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\BTHPORT\\Parameters\\Keys\""}

	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running PsExec:", err)
		return
	}

	fmt.Println("Registry Key Values:")
	fmt.Println(string(output))
}

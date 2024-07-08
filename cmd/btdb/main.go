package main

import (
	"btdb/internal/regfile"
	"btdb/internal/utility"
	"fmt"
	"log"
)

func main() {
	err := utility.CheckProgramAvailability("reged")
	if err != nil {
		log.Fatalln("`reged` is not available. Probably `chntpw` includes it.")
	}

	rf, err := regfile.Parse(
		"/home/qq/out.reg",
		`TESTPREFIX\ControlSet001\Services\BTHPORT\Parameters\Keys`,
	)
	if err != nil {
		log.Fatal(err)
	}
	for Interface, devices := range rf {
		fmt.Printf("Interface: %s, Value: %v\n", Interface, devices)
		for d, v := range devices.(map[string]interface{}) {
			fmt.Printf("Device: %s, Value: %v\n", d, v)
		}
	}
	fmt.Println(rf["b48c9d1b07c6"].(map[string]interface{})["CentralIRK"])

}

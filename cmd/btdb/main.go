package main

import (
	"btdb/internal/parser"
	"fmt"
	"log"
)

func main() {
	data, err := parser.ParseRegFile(
		"/home/qq/out.reg",
		`TESTPREFIX\ControlSet001\Services\BTHPORT\Parameters\Keys`,
	)
	if err != nil {
		log.Fatal(err)
	}
	for Interface, devices := range data {
		fmt.Printf("Interface: %s, Value: %v\n", Interface, devices)
		for d, v := range devices.(map[string]interface{}) {
			fmt.Printf("Device: %s, Value: %v\n", d, v)
		}
	}
	fmt.Println(data["b48c9d1b07c6"].(map[string]interface{})["CentralIRK"])

}

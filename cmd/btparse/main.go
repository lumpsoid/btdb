package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// parser.ListDevicesMac()

	// Open the registry hive file for reading
	file, err := os.Open("C:\\systemhive")
	if err != nil {
		fmt.Println("Error opening registry hive file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read from the file
	scanner := bufio.NewScanner(file)

	// Variables to hold key and value information
	var currentKey string
	var currentValues []string

	// Process each line of the file
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		// Skip empty lines or lines that don't start with a valid registry key
		if strings.TrimSpace(line) == "" || !strings.HasPrefix(line, "[") {
			continue
		}

		// Extract the registry key
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentKey = strings.TrimPrefix(line, "[")
			currentKey = strings.TrimSuffix(currentKey, "]")
			continue
		}

		// Extract registry values and their data
		if strings.HasPrefix(line, "@=") {
			// Handle default value
			valueData := strings.TrimPrefix(line, "@=")
			currentValues = append(currentValues, valueData)
		} else if strings.Contains(line, "=") {
			// Handle named values
			parts := strings.SplitN(line, "=", 2)
			valueName := strings.TrimSpace(parts[0])
			valueData := strings.TrimSpace(parts[1])

			// Example: Print or process valueName and valueData
			fmt.Printf("Key: %s, Name: %s, Data: %s\n", currentKey, valueName, valueData)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading registry hive file:", err)
		return
	}
}

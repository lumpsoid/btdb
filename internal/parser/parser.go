package parser

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func reverseSlice(slice interface{}) error {
	// Obtain the value and kind of the provided slice
	v := reflect.ValueOf(slice)
	k := v.Kind()

	// Check if the provided argument is a slice
	if k != reflect.Slice {
		return fmt.Errorf("provided argument is not a slice")
	}

	// Get the length of the slice
	length := v.Len()

	// Swap elements from the beginning and end of the slice
	for i := 0; i < length/2; i++ {
		j := length - i - 1
		// Swap elements at index i and j
		temp := v.Index(i).Interface()
		v.Index(i).Set(v.Index(j))
		v.Index(j).Set(reflect.ValueOf(temp))
	}
	return nil
}

func decodeDword(value string) (int64, error) {
	hexString := strings.TrimPrefix(value, "dword:")
	intValue, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func decodeHex(value string) ([]byte, error) {
	// Split the hex string by commas
	hexString := strings.Split(value, ":")[1]
	hexParts := strings.Split(hexString, ",")

	// Initialize a byte slice to store the converted bytes
	byteSlice := make([]byte, len(hexParts))

	// Convert each hex string to byte
	for i, hexPart := range hexParts {
		// Parse hex string to byte
		hexBytes, err := hex.DecodeString(hexPart)
		if err != nil {
			return nil, err
		}

		// Ensure we have exactly one byte per hex string part
		if len(hexBytes) != 1 {
			return nil, fmt.Errorf("unexpected length after decoding hex string: %v", hexBytes)
		}

		// Store the byte in the byte slice
		byteSlice[i] = hexBytes[0]
	}
	return byteSlice, nil
}

func parseValue(value string) (interface{}, error) {
	switch {
	case strings.HasPrefix(value, "dword:"):
		intValue, err := decodeDword(value)
		if err != nil {
			return nil, err
		}
		return intValue, nil
	case strings.HasPrefix(value, "hex"):
		intValue, err := decodeHex(value)
		if err != nil {
			return nil, err
		}
		return intValue, nil
	default:
		// return "", fmt.Errorf("unsupported value: %s", value)
		return value, nil
	}
}

func stripQuotes(value string) string {
	if len(value) > 1 && value[0] == '"' && value[len(value)-1] == '"' {
		return value[1 : len(value)-1]
	}
	return value
}

func parseRegLine(line string) (lineType string, key string, value interface{}, err error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return "empty", "", "", nil // skip empty lines or comments
	}

	// Extract section headers
	if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
		key = line[1 : len(line)-1]
		return "header", key, "", nil
	}

	// Extract key-value pairs
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "error", "", "", fmt.Errorf("invalid line format: %s", line)
	}
	key = strings.TrimSpace(parts[0])
	valuePre := strings.TrimSpace(parts[1])

	// Strip surrounding quotes from value
	key = stripQuotes(key)
	valuePre = stripQuotes(valuePre)

	parsedValue, err := parseValue(valuePre)
	if err != nil {
		return "error", "", "", err
	}

	return "pair", key, parsedValue, nil
}

func getRegistryVersion(scanner *bufio.Scanner) (string, error) {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Windows Registry Editor Version") {
			lineSplited := strings.Split(line, " ")
			return lineSplited[len(lineSplited)-1], nil
		}
	}
	return "", fmt.Errorf("registry version not found")
}

func checkRegistryVersion(ver string) (bool, error) {
	if ver != "5.00" {
		return false, fmt.Errorf("unsupported registry version: %s", ver)
	}
	return true, nil
}

func insertKey(
	lineType string,
	regData map[string]interface{},
	key string,
	value interface{},
) map[string]interface{} {
	currentMap := regData
	if lineType == "header" {
		keySlice := strings.Split(key, "\\")
		for _, key := range keySlice {
			_, ok := currentMap[key]
			if ok {
				currentMap = currentMap[key].(map[string]interface{})
			} else {
				if key != "" {
					currentMap[key] = make(map[string]interface{})
					currentMap = currentMap[key].(map[string]interface{})
				}
			}
		}
		return currentMap
	}
	if key != "" {
		currentMap[key] = value
	}
	return currentMap
}

func initDict(key string) map[string]interface{} {
	dictReg := make(map[string]interface{})
	currentDict := dictReg
	keySlice := strings.Split(key, "\\")
	for _, key := range keySlice {
		if _, ok := dictReg[key]; !ok {
			currentDict[key] = make(map[string]interface{})
			currentDict = currentDict[key].(map[string]interface{})
		}
	}
	return dictReg
}

func traverseNestedMap(regData map[string]interface{}, header string) (map[string]interface{}, bool) {
	keys := strings.Split(header, "\\")
	current := regData

	for _, key := range keys {
		value, ok := current[key]
		if !ok {
			return nil, false // Key not found
		}
		if nested, isMap := value.(map[string]interface{}); isMap {
			current = nested
		}
	}

	return current, true
}

func ParseRegFile(path string, header string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ver, err := getRegistryVersion(scanner)
	if err != nil {
		return nil, err
	}
	ok, err := checkRegistryVersion(ver)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, err
	}

	regData := make(map[string]interface{})
	currentDict := regData

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Trim(line, " ") == "" {
			continue
		}

		lineType, key, value, err := parseRegLine(line)
		if err != nil {
			return nil, err
		}
		if lineType == "header" {
			currentDict = insertKey(lineType, regData, key, value)
		} else {
			currentDict = insertKey(lineType, currentDict, key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	regDataReturn, found := traverseNestedMap(regData, header)
	if !found {
		return nil, fmt.Errorf("path not found in traversing: %s", header)
	}

	return regDataReturn, nil
}

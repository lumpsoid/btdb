package main

import (
	"btdb/internal/device"
	"btdb/internal/regfile"
	"btdb/internal/utility"
	"btdb/internal/winregistry"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	err := utility.CheckProgramAvailability("reged")
	if err != nil {
		log.Fatalln("`reged` is not available. Probably `chntpw` package includes it.")
	}

	// Define flags
	var mountpoint string
	var outputPathYaml string

	// Bind flags to variables
	flag.StringVar(&mountpoint, "mnt", "", "mountpoint for windows")
	flag.StringVar(&outputPathYaml, "output", "", "path to the final yaml file with devices information")

	flag.Parse()

  if mountpoint == "" {
		log.Fatal("-mnt parameter is empty. mount windows partition and retry")
  }

	prefix := "GOREGPARSE"
	registryQueryExport := `ControlSet001\Services\BTHPORT\Parameters\Keys`
	registryQuery := fmt.Sprint(prefix, `\ControlSet001\Services\BTHPORT\Parameters\Keys`)
	registryPath := "Windows/System32/config/SYSTEM"
	outputFilename := prefix + ".reg"
	outputPathReg := filepath.Join(os.TempDir(), outputFilename)

	os.Remove(outputPathReg)

	err = winregistry.Export(mountpoint, registryPath, prefix, registryQueryExport, outputPathReg)
	if err != nil {
		log.Fatal(err)
	}

	rf, err := regfile.Parse(outputPathReg, registryQuery)
	if err != nil {
		log.Fatal(err)
	}

	//err = os.Remove(outputPathReg)
	//if err != nil {
	//	log.Fatal(err)
	//}

	for btInterface, devices := range rf {
		fmt.Printf("Interface: %s\n", btInterface)
		for d, v := range devices.(map[string]interface{}) {
			switch vType := v.(type) {
			case map[string]interface{}:
				dN := device.New(d,
					vType["LTK"].([]byte),
					vType["IRK"].([]byte),
					vType["ERand"].([]uint8),
					vType["EDIV"].(int64),
				)
				fmt.Printf("\n%s", dN.String())
			case []uint8:
				// Handle the case where devices is []uint8
				fmt.Printf(
					"\nDevice: %s\n  Key: %v\n",
					d,
					strings.ToUpper(hex.EncodeToString(vType)),
				)
				// You might want to process the []uint8 here as needed
			default:
				// Handle unexpected types
				fmt.Printf("Unexpected type: %T\n", vType)
			}
		}
	}
}

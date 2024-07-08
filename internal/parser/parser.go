package parser

import (
	"btdb/internal/device"
	"fmt"
	"log"
	"reflect"

	"golang.org/x/sys/windows/registry"
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

func ParseTest() {
	rPath := `SYSTEM\CurrentControlSet\Services\BTHPORT\Parameters\Keys`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, rPath, registry.READ)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	names, err := k.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Names: ", names)
}

func ListDevicesMac() {
	rPath := `SYSTEM\CurrentControlSet\Services\BTHPORT\Parameters\Keys`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, rPath, registry.READ)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	btInterfaces, err := k.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range btInterfaces {
		rDevicePath := rPath + `\` + i
		kD, err := registry.OpenKey(registry.LOCAL_MACHINE, rDevicePath, registry.READ)
		if err != nil {
			log.Fatal(err)
		}
		defer kD.Close()

		btMacs, err := kD.ReadValueNames(-1)
		if err != nil {
			log.Fatal(err)
		}
		btSubMacs, err := kD.ReadSubKeyNames(-1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("BT Devices: ", btMacs, btSubMacs)
	}
}

func ParseDevices() {
	rPath := `SYSTEM\CurrentControlSet\Services\BTHPORT\Parameters\Keys`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, rPath, registry.READ)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	btInterfaces, err := k.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range btInterfaces {
		fmt.Println("Device: ", i)
		rDevicePath := rPath + `\` + i
		kD, err := registry.OpenKey(registry.LOCAL_MACHINE, rDevicePath, registry.READ)
		if err != nil {
			log.Fatal(err)
		}
		defer kD.Close()

		btMacs, err := kD.ReadSubKeyNames(-1)
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range btMacs {
			fmt.Println("Mac: ", m)
			rMacPath := rDevicePath + `\` + m
			kM, err := registry.OpenKey(registry.LOCAL_MACHINE, rMacPath, registry.READ)
			if err != nil {
				log.Fatal(err)
			}
			defer kM.Close()

			macValueNames, err := kM.ReadValueNames(-1)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Macs value names: %q\n", macValueNames)

			ltk, _, err := kM.GetBinaryValue("LTK")
			if err != nil {
				log.Fatal(err)
			}
			irk, _, err := kM.GetBinaryValue("IRK")
			if err != nil {
				log.Fatal(err)
			}
			err = reverseSlice(irk)
			if err != nil {
				log.Fatal(err)
			}
			eRand, _, err := kM.GetIntegerValue("ERand")
			if err != nil {
				log.Fatal(err)
			}
			ediv, _, err := kM.GetIntegerValue("EDIV")
			if err != nil {
				log.Fatal(err)
			}

			btDevice := device.New(m, ltk, irk, eRand, ediv)
			fmt.Printf("Device info: %s\n", btDevice)
		}
	}
}

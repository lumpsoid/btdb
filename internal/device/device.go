package device

import (
	"encoding/hex"
	"fmt"
)

type Device struct {
	Mac   string
	Ltk   string
	Irk   string
	Erand uint64
	Ediv  uint64
}

func New(
	mac string,
	ltk []byte,
	irk []byte,
	erand uint64,
	ediv uint64,
) Device {
	return Device{
		Mac:   mac,
		Ltk:   hex.EncodeToString(ltk),
		Irk:   hex.EncodeToString(irk),
		Erand: erand,
		Ediv:  ediv,
	}

}

func (d Device) String() string {
	return fmt.Sprintf(
		"Device: %s\nLTK: %s, IRK: %s, ERand: %d, EDIV: %d\n",
		d.Mac,
		d.Ltk,
		d.Irk,
		d.Erand,
		d.Ediv,
	)
}

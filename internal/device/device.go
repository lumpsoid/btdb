package device

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func reverseSlice[T any](s []T) []T {
	n := len(s)
	reversed := make([]T, n)
	copy(reversed, s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	return reversed
}

func uintSliceToByteSlice(uints []uint8) []byte {
	bytes := make([]byte, len(uints))
	for i, u := range uints {
		if u > 255 {
			fmt.Printf("Warning: uint value %d at index %d is out of byte range and will be truncated.\n", u, i)
		}
		bytes[i] = byte(u)
	}
	return bytes
}

type Bluetooth struct {
	Mac   string
	Ltk   string
	Irk   string
	Erand uint64
	Ediv  int64
}

func New(
	mac string,
	ltk []byte,
	irk []byte,
	erand []uint8,
	ediv int64,
) Bluetooth {
	return Bluetooth{
		Mac: mac,
		Ltk: hex.EncodeToString(ltk),
		Irk: hex.EncodeToString(irk),
		Erand: binary.LittleEndian.Uint64(erand),
		Ediv: ediv,
	}

}

func (d Bluetooth) String() string {
	return fmt.Sprintf(
		"Device: %s\n  LTK: %s\n  IRK: %s\n  ERand: %d\n  EDIV: %d\n",
		d.Mac,
		d.Ltk,
		d.Irk,
		d.Erand,
		d.Ediv,
	)
}

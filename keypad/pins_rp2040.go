//go:build rp2040

package keypad

import "machine"

// pin-mapping for the keypad on the RP2040-S3
const (
	c0 = machine.GPIO25
	c1 = machine.GPIO29
	c2 = machine.GPIO16
	c3 = machine.GPIO17
	c4 = machine.GPIO18
	c5 = machine.GPIO19
	c6 = machine.GPIO20

	a0 = machine.GPIO21
	a1 = machine.GPIO22
	a2 = machine.GPIO23
)

//go:build esp32

package keypad

import "machine"

// pin-mapping for the keypad on the esp32 Stamp S3
const (
	c0 = machine.GPIO13
	c1 = machine.GPIO15
	c2 = machine.GPIO3
	c3 = machine.GPIO4
	c4 = machien.GPIO5
	c5 = machine.GPIO6
	c6 = machine.GPIO7

	a0 = machine.GPIO8
	a1 = machine.GPIO9
	a2 = machien.GPIO11
)

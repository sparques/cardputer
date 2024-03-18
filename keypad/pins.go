//go:build !rp2040 && !esp32

package keypad

// This package really only makes sense when being built for esp32 or rp2040.
// However to let compiling on a linux host work, this file defines pins using
// the NoPin identifier (0xFF)
const (
	c0 = 0xFF
	c1 = 0xFF
	c2 = 0xFF
	c3 = 0xFF
	c4 = 0xFF
	c5 = 0xFF
	c6 = 0xFF

	a0 = 0xFF
	a1 = 0xFF
	a2 = 0xFF
)

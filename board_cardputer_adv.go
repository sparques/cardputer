//go:build esp32 && cardputer_adv

package cardputer // github.com/sparques/cardputer

import "machine"

// Board identifies the Cardputer-Adv build.
const Board = "cardputer-adv"

const (
	// Cardputer-Adv keyboard controller pins.
	KeypadIRQ = machine.GPIO11
	KeypadSDA = machine.GPIO8
	KeypadSCL = machine.GPIO9
)

const (
	// Shared I2C bus for the keypad controller, audio codec, and IMU.
	I2CSharedSDA = machine.GPIO8
	I2CSharedSCL = machine.GPIO9
)

const (
	// Cardputer-Adv audio pins routed through the ES8311 codec.
	I2SClock    = machine.GPIO43
	SpeakerBK   = machine.GPIO41
	SpeakerData = machine.GPIO42
	MicData     = machine.GPIO46
)

const (
	// Cardputer-Adv expansion connector pins.
	ExtReset = machine.GPIO3
	ExtInt   = machine.GPIO4
	ExtBusy  = machine.GPIO6
	ExtCS    = machine.GPIO5
	ExtSCK   = machine.GPIO40
	ExtMOSI  = machine.GPIO14
	ExtMISO  = machine.GPIO39
	ExtSDA   = machine.GPIO8
	ExtSCL   = machine.GPIO9
	ExtRX    = machine.GPIO13
	ExtTX    = machine.GPIO15
)

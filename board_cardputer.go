//go:build esp32 && !cardputer_adv

package cardputer // github.com/sparques/cardputer

import "machine"

const Board = "cardputer"

const (
	// Keypad pins
	// KeypadCX are the sense lines for the keypad.
	KeypadC0 = machine.GPIO13
	KeypadC1 = machine.GPIO15
	KeypadC2 = machine.GPIO3
	KeypadC3 = machine.GPIO4
	KeypadC4 = machine.GPIO5
	KeypadC5 = machine.GPIO6
	KeypadC6 = machine.GPIO7
	// KeypadAX are the address lines for the keypad.
	KeypadA0 = machine.GPIO8
	KeypadA1 = machine.GPIO9
	KeypadA2 = machine.GPIO11
)

const (
	// I2S Speaker and Mic
	// I2SClock is shared between speaker and microphone.
	I2SClock    = machine.GPIO43
	SpeakerBK   = machine.GPIO41
	SpeakerData = machine.GPIO42
	MicData     = machine.GPIO46
)

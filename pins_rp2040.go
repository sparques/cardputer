//go:build rp2040

package cardputer // github.com/sparques/cardputer

import "machine"

const (
	// Keypad pins
	// KeypadCX are the sense lines for the keypad
	KeypadC0 = machine.GPIO25
	KeypadC1 = machine.GPIO29
	KeypadC2 = machine.GPIO16
	KeypadC3 = machine.GPIO17
	KeypadC4 = machine.GPIO18
	KeypadC5 = machine.GPIO19
	KeypadC6 = machine.GPIO20
	// KeypadAX are the address lines for the keypad
	KeypadA0 = machine.GPIO21
	KeypadA1 = machine.GPIO22
	KeypadA2 = machine.GPIO23
)

const (
	// Grove Port
	// Grove i2c pins
	GroveSCL = machine.GPIO13
	GroveSDA = machine.GPIO12
	// Grove UART pins
	GroveRX = machine.GPIO13
	GroveTx = machine.GPIO12
)

const (
	// Battery Voltage
	BatterySense = machine.ADC0
)

const (
	// SD Card
	SDMOSI = machine.GPIO27
	SDMISO = machine.GPIO8
	SDSCK  = machine.GPIO10
	SDCS   = machine.GPIO24
)

const (
	// I2S Speaker and Mic
	// IS2Clock is shared between Speaker and Mic
	I2SClock    = machine.GPIO15
	SpeakerBK   = machine.GPIO1
	SpeakerData = machine.GPIO11
	MicData     = machine.GPIO14
)

// Infrared LED
const (
	IRPin = machine.GPIO2
)

const (
	// LCD / ST7789 pins
	LCDBacklight = machine.GPIO3
	LCDReset     = machine.GPIO4
	LCDRS        = machine.GPIO5
	LCDMOSI      = machine.GPIO7
	LCDSCK       = machine.GPIO6
	LCDCS        = machine.GPIO9
)

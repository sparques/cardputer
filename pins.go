//go:build esp32 || esp32s3

package cardputer

import "machine"

const (
	// GroveSCL and GroveSDA are the Cardputer family HY2.0-4P "custom" port pins.
	// The same pins also double as the port's UART signals.
	GroveSCL = machine.GPIO2
	GroveSDA = machine.GPIO1
	GroveRX  = machine.GPIO2
	GroveTx  = machine.GPIO1
)

const (
	// BatterySense is the ADC input connected to the battery divider.
	BatterySense = machine.GPIO10
)

const (
	// Built-in microSD slot pins.
	SDMOSI = machine.GPIO14
	SDMISO = machine.GPIO39
	SDSCK  = machine.GPIO40
	SDCS   = machine.GPIO12
)

const (
	// IRPin drives the built-in infrared LED.
	IRPin = machine.GPIO44
)

const (
	// Built-in ST7789 LCD pins.
	LCDBacklight = machine.GPIO38
	LCDReset     = machine.GPIO33
	LCDRS        = machine.GPIO34
	LCDMOSI      = machine.GPIO35
	LCDSCK       = machine.GPIO36
	LCDCS        = machine.GPIO37
)

const (
	// Shared audio data pins present on both Cardputer variants.
	I2SClock    = machine.GPIO43
	SpeakerBK   = machine.GPIO41
	SpeakerData = machine.GPIO42
	MicData     = machine.GPIO46
)

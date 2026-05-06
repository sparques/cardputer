//go:build esp32

package cardputer

import "machine"

const (
	// Grove Port
	// On Cardputer devices this is the HY2.0-4P "custom" port: yellow is GPIO2, white is GPIO1.
	GroveSCL = machine.GPIO2
	GroveSDA = machine.GPIO1
	GroveRX  = machine.GPIO2
	GroveTx  = machine.GPIO1
)

const (
	// Battery Voltage
	BatterySense = machine.GPIO10
)

const (
	// SD Card
	SDMOSI = machine.GPIO14
	SDMISO = machine.GPIO39
	SDSCK  = machine.GPIO40
	SDCS   = machine.GPIO12
)

const (
	// Infrared LED
	IRPin = machine.GPIO44
)

const (
	// LCD / ST7789 pins
	LCDBacklight = machine.GPIO38
	LCDReset     = machine.GPIO33
	LCDRS        = machine.GPIO34
	LCDMOSI      = machine.GPIO35
	LCDSCK       = machine.GPIO36
	LCDCS        = machine.GPIO37
)

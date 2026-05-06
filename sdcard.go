package cardputer

import (
	"machine"

	"tinygo.org/x/drivers/sdcard"
)

// SDCard exposes the built-in microSD slot as a TinyGo sdcard.Device.
var SDCard = &sdc{
	Device: sdcard.New(machine.SPI1, SDSCK, SDMOSI, SDMISO, SDCS),
}

type sdc struct {
	sdcard.Device
}

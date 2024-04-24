package cardputer

import (
	"machine"

	"tinygo.org/x/drivers/sdcard"
)

var SDCard = &sdc{
	Device: sdcard.New(machine.SPI1, SDSCK, SDMOSI, SDMISO, SDCS),
}

type sdc struct {
	sdcard.Device
}

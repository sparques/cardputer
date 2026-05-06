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

// Init configures the built-in microSD slot and probes the inserted card.
func (s *sdc) Init() error {
	return s.Device.Configure()
}

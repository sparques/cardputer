//go:build esp32 && cardputer_adv

package cardputer

import (
	"time"

	"github.com/sparques/cardputer/internal/adv"
	"machine"
)

const (
	es8311Addr = 0x18
)

const (
	es8311VolumeMute = 0x00
	es8311Volume0dB  = 0xBF
)

// ES8311 register map.
const (
	es8311RegReset      = 0x00
	es8311RegClock1     = 0x01
	es8311RegClock2     = 0x02
	es8311RegClock3     = 0x03
	es8311RegClock4     = 0x04
	es8311RegClock5     = 0x05
	es8311RegClock6     = 0x06
	es8311RegClock7     = 0x07
	es8311RegClock8     = 0x08
	es8311RegClock9     = 0x09
	es8311RegClock10    = 0x0A
	es8311RegSystem1    = 0x0B
	es8311RegSystem2    = 0x0C
	es8311RegSystem3    = 0x0D
	es8311RegSystem4    = 0x0E
	es8311RegSystem5    = 0x0F
	es8311RegSystem6    = 0x10
	es8311RegSystem7    = 0x11
	es8311RegSystem8    = 0x12
	es8311RegSystem9    = 0x13
	es8311RegSystem10   = 0x14
	es8311RegSystem11   = 0x15
	es8311RegADC1       = 0x16
	es8311RegADCVolume  = 0x17
	es8311RegADCALC1    = 0x18
	es8311RegADCALC2    = 0x19
	es8311RegADCALC3    = 0x1A
	es8311RegDAC1       = 0x31
	es8311RegDACVolume  = 0x32
	es8311RegDACOffset  = 0x33
	es8311RegDACDRC1    = 0x34
	es8311RegDACDRC2    = 0x35
	es8311RegDAC6       = 0x37
	es8311RegDAC7       = 0x38
	es8311RegGPIO       = 0x44
	es8311RegGP1        = 0x45
	es8311RegGP2        = 0x46
	es8311RegChipStatus = 0x4F
)

const (
	es8311ResetShutdown = 0x1F
	es8311ResetRelease  = 0x80
)

const (
	es8311DACMuteMask = 0x60
	es8311ALCEnable   = 0x80
	es8311DRCEnable   = 0x80
)

type ES8311Config struct {
	DACVolume uint8
	ADCVolume uint8
	ADCGain   uint8
	Muted     bool
}

func DefaultES8311Config() ES8311Config {
	return ES8311Config{
		DACVolume: es8311Volume0dB,
		ADCVolume: es8311Volume0dB,
		ADCGain:   4, // datasheet default: +24dB
		Muted:     true,
	}
}

// ES8311 provides low-level register access for the Cardputer-Adv audio codec.
type ES8311 struct {
	bus     *machine.I2C
	address uint16
}

func newES8311(bus *machine.I2C, address uint16) *ES8311 {
	return &ES8311{
		bus:     bus,
		address: address,
	}
}

func openES8311() (*ES8311, error) {
	bus, err := adv.SharedI2C()
	if err != nil {
		return nil, err
	}
	return newES8311(bus, es8311Addr), nil
}

func (c *ES8311) ReadRegister(reg uint8) (uint8, error) {
	buf := []byte{reg}
	out := []byte{0}
	if err := c.bus.Tx(c.address, buf, out); err != nil {
		return 0, err
	}
	return out[0], nil
}

func (c *ES8311) WriteRegister(reg, value uint8) error {
	return c.bus.Tx(c.address, []byte{reg, value}, nil)
}

func (c *ES8311) UpdateRegisterBits(reg, mask, value uint8) error {
	cur, err := c.ReadRegister(reg)
	if err != nil {
		return err
	}
	next := (cur &^ mask) | (value & mask)
	if next == cur {
		return nil
	}
	return c.WriteRegister(reg, next)
}

func (c *ES8311) Reset() error {
	if err := c.WriteRegister(es8311RegReset, es8311ResetShutdown); err != nil {
		return err
	}
	time.Sleep(2 * time.Millisecond)
	if err := c.WriteRegister(es8311RegReset, 0x00); err != nil {
		return err
	}
	time.Sleep(2 * time.Millisecond)
	return c.WriteRegister(es8311RegReset, es8311ResetRelease)
}

// ConfigureDefaults applies a conservative software-only codec configuration.
// This intentionally avoids clock-tree and I2S format programming until the
// transport layer is implemented.
func (c *ES8311) ConfigureDefaults(cfg ES8311Config) error {
	if err := c.DisableALC(); err != nil {
		return err
	}
	if err := c.DisableDRC(); err != nil {
		return err
	}
	if err := c.SetADCGain(cfg.ADCGain); err != nil {
		return err
	}
	if err := c.SetADCVolume(cfg.ADCVolume); err != nil {
		return err
	}
	if err := c.SetDACVolume(cfg.DACVolume); err != nil {
		return err
	}
	return c.SetMuted(cfg.Muted)
}

func (c *ES8311) SetMuted(muted bool) error {
	var value uint8
	if muted {
		value = es8311DACMuteMask
	}
	return c.UpdateRegisterBits(es8311RegDAC1, es8311DACMuteMask, value)
}

func (c *ES8311) DisableALC() error {
	return c.UpdateRegisterBits(es8311RegADCALC1, es8311ALCEnable, 0)
}

func (c *ES8311) DisableDRC() error {
	return c.UpdateRegisterBits(es8311RegDACDRC1, es8311DRCEnable, 0)
}

func (c *ES8311) SetADCGain(scale uint8) error {
	if scale > 7 {
		scale = 7
	}
	return c.UpdateRegisterBits(es8311RegADC1, 0x07, scale)
}

func (c *ES8311) SetDACVolume(volume uint8) error {
	return c.WriteRegister(es8311RegDACVolume, volume)
}

func (c *ES8311) SetADCVolume(volume uint8) error {
	return c.WriteRegister(es8311RegADCVolume, volume)
}

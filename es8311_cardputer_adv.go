//go:build esp32 && cardputer_adv

package cardputer

import (
	"errors"
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
	es8311RegSystem14   = 0x1C
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

var errES8311UnsupportedClock = errors.New("unsupported ES8311 clock configuration")

type es8311ClockCoefficient struct {
	mclk    uint32
	rate    uint32
	preDiv  uint8
	preMult uint8
	adcDiv  uint8
	dacDiv  uint8
	fsMode  uint8
	lrckH   uint8
	lrckL   uint8
	bclkDiv uint8
	adcOSR  uint8
	dacOSR  uint8
}

var es8311ClockCoefficients = [...]es8311ClockCoefficient{
	{4096000, 16000, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0xff, 0x04, 0x10, 0x20},
	{5644800, 22050, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0xff, 0x04, 0x10, 0x10},
	{6144000, 24000, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0xff, 0x04, 0x10, 0x10},
	{8192000, 32000, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0xff, 0x04, 0x10, 0x10},
	{11289600, 44100, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0xff, 0x04, 0x10, 0x10},
	{12288000, 48000, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0xff, 0x04, 0x10, 0x10},
}

// ES8311Config describes the shared Adv audio codec configuration.
type ES8311Config struct {
	// SampleRate is the target PCM sample rate in hertz.
	SampleRate uint32
	// BitsPerSample is the target PCM word size.
	BitsPerSample ES8311Resolution
	// UseMCLK requests a master clock driven at 256*SampleRate.
	UseMCLK bool
	// UseMicrophone enables the codec microphone input path.
	UseMicrophone bool
	// DACVolume is the raw DAC volume register value.
	DACVolume uint8
	// ADCVolume is the raw ADC volume register value.
	ADCVolume uint8
	// ADCGain is the ADC gain field value.
	ADCGain uint8
	// Muted requests the DAC mute bit.
	Muted bool
}

// DefaultES8311Config returns the package's conservative Adv audio defaults.
// The default profile is muted, uses 16-bit mono PCM at 16kHz, and enables the microphone path.
func DefaultES8311Config() ES8311Config {
	return ES8311Config{
		SampleRate:    16000,
		BitsPerSample: ES8311Resolution16,
		UseMCLK:       true,
		UseMicrophone: true,
		DACVolume:     es8311Volume0dB,
		ADCVolume:     0xC8,
		ADCGain:       4, // datasheet default: +24dB
		Muted:         true,
	}
}

func normalizeES8311Config(cfg ES8311Config) ES8311Config {
	defaults := DefaultES8311Config()
	if cfg.SampleRate == 0 {
		cfg.SampleRate = defaults.SampleRate
	}
	if cfg.BitsPerSample == 0 {
		cfg.BitsPerSample = defaults.BitsPerSample
	}
	if cfg.ADCGain > 7 {
		cfg.ADCGain = 7
	}
	return cfg
}

// ES8311 provides low-level register access for the Cardputer-Adv audio codec.
type ES8311 struct {
	bus     *machine.I2C
	address uint16
}

var (
	sharedES8311         *ES8311
	sharedES8311Defaults bool
	sharedES8311Config   = DefaultES8311Config()
)

func newES8311(bus *machine.I2C, address uint16) *ES8311 {
	return &ES8311{
		bus:     bus,
		address: address,
	}
}

func openES8311() (*ES8311, error) {
	if sharedES8311 != nil {
		return sharedES8311, nil
	}
	bus, err := adv.SharedI2C()
	if err != nil {
		return nil, err
	}
	sharedES8311 = newES8311(bus, es8311Addr)
	return sharedES8311, nil
}

// ReadRegister reads a single ES8311 register.
func (c *ES8311) ReadRegister(reg uint8) (uint8, error) {
	buf := []byte{reg}
	out := []byte{0}
	if err := c.bus.Tx(c.address, buf, out); err != nil {
		return 0, err
	}
	return out[0], nil
}

// WriteRegister writes a single ES8311 register.
func (c *ES8311) WriteRegister(reg, value uint8) error {
	return c.bus.Tx(c.address, []byte{reg, value}, nil)
}

// UpdateRegisterBits updates only the masked bits in an ES8311 register.
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

// Reset performs a software reset sequence for the codec.
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

// ConfigureDefaults applies a conservative codec configuration for future I2S
// transport use. This assumes a standard I2S link with MCLK = 256*Fs.
func (c *ES8311) ConfigureDefaults(cfg ES8311Config) error {
	if err := c.ConfigureClock(cfg); err != nil {
		return err
	}
	if err := c.ConfigureFormat(cfg); err != nil {
		return err
	}
	if err := c.ConfigureMicrophone(cfg); err != nil {
		return err
	}
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
	if err := c.PowerUpPlaybackPath(); err != nil {
		return err
	}
	return c.SetMuted(cfg.Muted)
}

func initES8311Defaults(cfg ES8311Config) (*ES8311, error) {
	cfg = normalizeES8311Config(cfg)
	codec, err := openES8311()
	if err != nil {
		return nil, err
	}
	if sharedES8311Defaults {
		return codec, nil
	}
	if err := codec.Reset(); err != nil {
		return nil, err
	}
	if err := codec.ConfigureDefaults(cfg); err != nil {
		return nil, err
	}
	sharedES8311Defaults = true
	sharedES8311Config = cfg
	return codec, nil
}

func configureSharedES8311(update func(*ES8311Config)) (*ES8311, error) {
	cfg := sharedES8311Config
	update(&cfg)
	cfg = normalizeES8311Config(cfg)

	codec, err := initES8311Defaults(sharedES8311Config)
	if err != nil {
		return nil, err
	}
	if err := codec.ConfigureDefaults(cfg); err != nil {
		return nil, err
	}
	if err := configureSharedAudioTransport(cfg); err != nil {
		return nil, err
	}
	sharedES8311Config = cfg
	return codec, nil
}

func audioTransportConfigFromES8311(cfg ES8311Config) AudioTransportConfig {
	return AudioTransportConfig{
		SampleRate:    cfg.SampleRate,
		BitsPerSample: cfg.BitsPerSample,
		Channels:      1,
		UseMCLK:       cfg.UseMCLK,
	}
}

// AudioConfig returns the currently cached shared Adv audio configuration.
func AudioConfig() ES8311Config {
	return sharedES8311Config
}

// ConfigureAudio applies a complete shared Adv audio configuration.
func ConfigureAudio(cfg ES8311Config) error {
	_, err := configureSharedES8311(func(current *ES8311Config) {
		*current = cfg
	})
	return err
}

// SetMuted updates the codec DAC mute state.
func (c *ES8311) SetMuted(muted bool) error {
	var value uint8
	if muted {
		value = es8311DACMuteMask
	}
	return c.UpdateRegisterBits(es8311RegDAC1, es8311DACMuteMask, value)
}

// ConfigureClock programs the codec clock tree for the supplied PCM settings.
// Only MCLK=256*sample-rate profiles present in the built-in coefficient table are supported.
func (c *ES8311) ConfigureClock(cfg ES8311Config) error {
	if !cfg.UseMCLK {
		return errES8311UnsupportedClock
	}

	coeff, ok := lookupES8311ClockCoefficient(cfg.SampleRate, cfg.SampleRate*256)
	if !ok {
		return errES8311UnsupportedClock
	}

	if err := c.WriteRegister(es8311RegClock1, 0x3F); err != nil {
		return err
	}

	reg02, err := c.ReadRegister(es8311RegClock2)
	if err != nil {
		return err
	}
	reg02 &= 0x07
	reg02 |= (coeff.preDiv - 1) << 5
	reg02 |= coeff.preMult << 3
	if err := c.WriteRegister(es8311RegClock2, reg02); err != nil {
		return err
	}

	if err := c.WriteRegister(es8311RegClock3, (coeff.fsMode<<6)|coeff.adcOSR); err != nil {
		return err
	}
	if err := c.WriteRegister(es8311RegClock4, coeff.dacOSR); err != nil {
		return err
	}
	if err := c.WriteRegister(es8311RegClock5, ((coeff.adcDiv-1)<<4)|(coeff.dacDiv-1)); err != nil {
		return err
	}

	reg06, err := c.ReadRegister(es8311RegClock6)
	if err != nil {
		return err
	}
	reg06 &= 0xE0
	if coeff.bclkDiv < 19 {
		reg06 |= coeff.bclkDiv - 1
	} else {
		reg06 |= coeff.bclkDiv
	}
	if err := c.WriteRegister(es8311RegClock6, reg06); err != nil {
		return err
	}
	if err := c.WriteRegister(es8311RegClock7, coeff.lrckH); err != nil {
		return err
	}
	return c.WriteRegister(es8311RegClock8, coeff.lrckL)
}

// ConfigureFormat programs the codec serial audio word size.
func (c *ES8311) ConfigureFormat(cfg ES8311Config) error {
	reg00, err := c.ReadRegister(es8311RegReset)
	if err != nil {
		return err
	}
	reg00 &^= 0x40
	if err := c.WriteRegister(es8311RegReset, reg00); err != nil {
		return err
	}

	resolution := es8311ResolutionBits(cfg.BitsPerSample)
	if err := c.WriteRegister(es8311RegClock9, resolution); err != nil {
		return err
	}
	return c.WriteRegister(es8311RegClock10, resolution)
}

// ConfigureMicrophone programs the codec microphone input path and gain settings.
func (c *ES8311) ConfigureMicrophone(cfg ES8311Config) error {
	reg14 := uint8(0x1A)
	if cfg.UseMicrophone {
		reg14 |= 1 << 6
	}
	if err := c.WriteRegister(es8311RegSystem10, reg14); err != nil {
		return err
	}
	if err := c.SetADCGain(cfg.ADCGain); err != nil {
		return err
	}
	return c.SetADCVolume(cfg.ADCVolume)
}

// PowerUpPlaybackPath enables the codec playback path registers used by the Adv speaker output.
func (c *ES8311) PowerUpPlaybackPath() error {
	if err := c.WriteRegister(es8311RegSystem3, 0x01); err != nil {
		return err
	}
	if err := c.WriteRegister(es8311RegSystem4, 0x02); err != nil {
		return err
	}
	if err := c.WriteRegister(es8311RegSystem8, 0x00); err != nil {
		return err
	}
	if err := c.WriteRegister(es8311RegSystem9, 0x10); err != nil {
		return err
	}
	if err := c.WriteRegister(es8311RegDAC6, 0x08); err != nil {
		return err
	}
	return c.WriteRegister(es8311RegSystem14, 0x6A)
}

// DisableALC disables the codec automatic level control block.
func (c *ES8311) DisableALC() error {
	return c.UpdateRegisterBits(es8311RegADCALC1, es8311ALCEnable, 0)
}

// DisableDRC disables the codec dynamic range control block.
func (c *ES8311) DisableDRC() error {
	return c.UpdateRegisterBits(es8311RegDACDRC1, es8311DRCEnable, 0)
}

// SetADCGain sets the codec ADC gain field, clamped to the device's 3-bit range.
func (c *ES8311) SetADCGain(scale uint8) error {
	if scale > 7 {
		scale = 7
	}
	return c.UpdateRegisterBits(es8311RegADC1, 0x07, scale)
}

// SetDACVolume writes the raw codec DAC volume register value.
func (c *ES8311) SetDACVolume(volume uint8) error {
	return c.WriteRegister(es8311RegDACVolume, volume)
}

// SetADCVolume writes the raw codec ADC volume register value.
func (c *ES8311) SetADCVolume(volume uint8) error {
	return c.WriteRegister(es8311RegADCVolume, volume)
}

func es8311ResolutionBits(resolution ES8311Resolution) uint8 {
	switch resolution {
	case ES8311Resolution16:
		return 3 << 2
	case ES8311Resolution18:
		return 2 << 2
	case ES8311Resolution20:
		return 1 << 2
	case ES8311Resolution24:
		return 0 << 2
	case ES8311Resolution32:
		return 4 << 2
	default:
		return 0
	}
}

func lookupES8311ClockCoefficient(rate, mclk uint32) (es8311ClockCoefficient, bool) {
	for _, coeff := range es8311ClockCoefficients {
		if coeff.rate == rate && coeff.mclk == mclk {
			return coeff, true
		}
	}
	return es8311ClockCoefficient{}, false
}

func percentToES8311Volume(percent uint8) uint8 {
	if percent >= 100 {
		return es8311Volume0dB
	}
	if percent == 0 {
		return es8311VolumeMute
	}
	return uint8((uint16(percent) * uint16(es8311Volume0dB)) / 100)
}

//go:build esp32 && cardputer_adv

package cardputer

// Speaker exposes the Cardputer-Adv speaker path.
// Codec configuration is supported; PCM streaming still depends on a future ESP32 I2S backend.
var Speaker = &speaker{}

type speaker struct {
	codec     *ES8311
	transport audioTransport
}

func (spk *speaker) Init() error {
	codec, err := initES8311Defaults(DefaultES8311Config())
	if err != nil {
		return err
	}
	if err := configureSharedAudioTransport(sharedES8311Config); err != nil {
		return err
	}
	transport, err := openAudioTransport()
	if err != nil {
		return err
	}
	spk.codec = codec
	spk.transport = transport
	return nil
}

func (spk *speaker) SetMuted(muted bool) error {
	codec, err := configureSharedES8311(func(cfg *ES8311Config) {
		cfg.Muted = muted
	})
	if err != nil {
		return err
	}
	spk.codec = codec
	return spk.attachTransport()
}

func (spk *speaker) SetVolume(volume uint8) error {
	codec, err := configureSharedES8311(func(cfg *ES8311Config) {
		cfg.DACVolume = percentToES8311Volume(volume)
	})
	if err != nil {
		return err
	}
	spk.codec = codec
	return spk.attachTransport()
}

func (spk *speaker) SetSampleRate(rate uint32) error {
	codec, err := configureSharedES8311(func(cfg *ES8311Config) {
		cfg.SampleRate = rate
	})
	if err != nil {
		return err
	}
	spk.codec = codec
	return spk.attachTransport()
}

func (spk *speaker) SetBitsPerSample(bits ES8311Resolution) error {
	codec, err := configureSharedES8311(func(cfg *ES8311Config) {
		cfg.BitsPerSample = bits
	})
	if err != nil {
		return err
	}
	spk.codec = codec
	return spk.attachTransport()
}

func (spk *speaker) Write(samples []int16) (int, error) {
	if err := spk.Init(); err != nil {
		return 0, err
	}
	return spk.transport.Write(samples)
}

func (spk *speaker) attachTransport() error {
	transport, err := openAudioTransport()
	if err != nil {
		return err
	}
	spk.transport = transport
	return nil
}

// Beep is currently a no-op until I2S streaming support is added for ESP32-S3.
func (*speaker) Beep() {}

// AckBeep is currently a no-op until I2S streaming support is added for ESP32-S3.
func (*speaker) AckBeep() {}

// NakBeep is currently a no-op until I2S streaming support is added for ESP32-S3.
func (*speaker) NakBeep() {}

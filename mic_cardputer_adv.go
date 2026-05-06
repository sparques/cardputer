//go:build esp32 && cardputer_adv

package cardputer

// Microphone exposes the Cardputer-Adv microphone path.
// Codec configuration is supported; PCM capture still depends on a future ESP32 I2S backend.
var Microphone = &microphone{}

type microphone struct {
	codec     *ES8311
	transport audioTransport
}

func (mic *microphone) Init() error {
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
	mic.codec = codec
	mic.transport = transport
	return nil
}

func (mic *microphone) SetGain(gain uint8) error {
	codec, err := configureSharedES8311(func(cfg *ES8311Config) {
		cfg.ADCGain = gain
	})
	if err != nil {
		return err
	}
	mic.codec = codec
	return mic.attachTransport()
}

func (mic *microphone) SetVolume(volume uint8) error {
	codec, err := configureSharedES8311(func(cfg *ES8311Config) {
		cfg.ADCVolume = percentToES8311Volume(volume)
	})
	if err != nil {
		return err
	}
	mic.codec = codec
	return mic.attachTransport()
}

func (mic *microphone) SetSampleRate(rate uint32) error {
	codec, err := configureSharedES8311(func(cfg *ES8311Config) {
		cfg.SampleRate = rate
	})
	if err != nil {
		return err
	}
	mic.codec = codec
	return mic.attachTransport()
}

func (mic *microphone) SetBitsPerSample(bits ES8311Resolution) error {
	codec, err := configureSharedES8311(func(cfg *ES8311Config) {
		cfg.BitsPerSample = bits
	})
	if err != nil {
		return err
	}
	mic.codec = codec
	return mic.attachTransport()
}

func (mic *microphone) UseDigital(enable bool) error {
	codec, err := configureSharedES8311(func(cfg *ES8311Config) {
		cfg.UseMicrophone = enable
	})
	if err != nil {
		return err
	}
	mic.codec = codec
	return mic.attachTransport()
}

func (mic *microphone) Read(samples []int16) (int, error) {
	if err := mic.Init(); err != nil {
		return 0, err
	}
	return mic.transport.Read(samples)
}

func (mic *microphone) attachTransport() error {
	transport, err := openAudioTransport()
	if err != nil {
		return err
	}
	mic.transport = transport
	return nil
}

//go:build esp32 && cardputer_adv

package cardputer

var Speaker = &speaker{}

type speaker struct {
	codec *ES8311
}

func (spk *speaker) Init() error {
	codec, err := openES8311()
	if err != nil {
		return err
	}
	if err := codec.Reset(); err != nil {
		return err
	}
	if err := codec.ConfigureDefaults(DefaultES8311Config()); err != nil {
		return err
	}
	spk.codec = codec
	return nil
}

// Beep is currently a no-op until I2S streaming support is added for ESP32-S3.
func (*speaker) Beep() {}

// AckBeep is currently a no-op until I2S streaming support is added for ESP32-S3.
func (*speaker) AckBeep() {}

// NakBeep is currently a no-op until I2S streaming support is added for ESP32-S3.
func (*speaker) NakBeep() {}

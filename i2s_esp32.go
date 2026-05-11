//go:build (esp32 || esp32s3) && !cardputer_adv

package cardputer

import "errors"

// Speaker exposes speaker control for the current board.
// On the original Cardputer build, configuration and streaming are not yet implemented.
var Speaker = &speaker{}

type speaker struct{}

func (*speaker) Init() error {
	return errors.New("speaker support is not implemented for the ESP32-S3 cardputer target")
}

func (*speaker) SetMuted(bool) error {
	return errors.New("speaker support is not implemented for the ESP32-S3 cardputer target")
}

func (*speaker) SetVolume(uint8) error {
	return errors.New("speaker support is not implemented for the ESP32-S3 cardputer target")
}

func (*speaker) SetSampleRate(uint32) error {
	return errors.New("speaker support is not implemented for the ESP32-S3 cardputer target")
}

func (*speaker) SetBitsPerSample(ES8311Resolution) error {
	return errors.New("speaker support is not implemented for the ESP32-S3 cardputer target")
}

func (*speaker) Write([]int16) (int, error) {
	return 0, errAudioStreamUnavailable
}

// Beep is currently a no-op until audio support is implemented for ESP32-S3.
func (*speaker) Beep() {}

// AckBeep is currently a no-op until audio support is implemented for ESP32-S3.
func (*speaker) AckBeep() {}

// NakBeep is currently a no-op until audio support is implemented for ESP32-S3.
func (*speaker) NakBeep() {}

//go:build (esp32 || esp32s3) && !cardputer_adv

package cardputer

import "errors"

// Microphone exposes microphone control for the current board.
// On the original Cardputer build, configuration and capture are not yet implemented.
var Microphone = &microphone{}

type microphone struct{}

func (*microphone) Init() error {
	return errors.New("microphone support is not implemented for the ESP32-S3 cardputer target")
}

func (*microphone) SetGain(uint8) error {
	return errors.New("microphone support is not implemented for the ESP32-S3 cardputer target")
}

func (*microphone) SetVolume(uint8) error {
	return errors.New("microphone support is not implemented for the ESP32-S3 cardputer target")
}

func (*microphone) SetSampleRate(uint32) error {
	return errors.New("microphone support is not implemented for the ESP32-S3 cardputer target")
}

func (*microphone) SetBitsPerSample(ES8311Resolution) error {
	return errors.New("microphone support is not implemented for the ESP32-S3 cardputer target")
}

func (*microphone) UseDigital(bool) error {
	return errors.New("microphone support is not implemented for the ESP32-S3 cardputer target")
}

func (*microphone) Read([]int16) (int, error) {
	return 0, errAudioStreamUnavailable
}

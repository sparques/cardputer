//go:build esp32 && !cardputer_adv

package cardputer

import "errors"

var Speaker = &speaker{}

type speaker struct{}

func (*speaker) Init() error {
	return errors.New("speaker support is not implemented for the ESP32-S3 cardputer target")
}

// Beep is currently a no-op until audio support is implemented for ESP32-S3.
func (*speaker) Beep() {}

// AckBeep is currently a no-op until audio support is implemented for ESP32-S3.
func (*speaker) AckBeep() {}

// NakBeep is currently a no-op until audio support is implemented for ESP32-S3.
func (*speaker) NakBeep() {}

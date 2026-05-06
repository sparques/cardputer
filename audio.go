package cardputer

import "errors"

var errAudioStreamUnavailable = errors.New("audio streaming is not implemented for this target")

// ES8311Resolution describes the sample width programmed into the Adv audio codec.
type ES8311Resolution uint8

const (
	// Supported ES8311 PCM sample widths.
	ES8311Resolution16 ES8311Resolution = 16
	ES8311Resolution18 ES8311Resolution = 18
	ES8311Resolution20 ES8311Resolution = 20
	ES8311Resolution24 ES8311Resolution = 24
	ES8311Resolution32 ES8311Resolution = 32
)

// AudioTransportConfig describes the PCM link settings expected by the
// board-specific audio transport layer.
type AudioTransportConfig struct {
	// SampleRate is the PCM sample rate in hertz.
	SampleRate uint32
	// BitsPerSample is the PCM word size.
	BitsPerSample ES8311Resolution
	// Channels is the channel count expected by the transport.
	Channels uint8
	// UseMCLK reports whether the transport is expected to drive a master clock.
	UseMCLK bool
}

type audioTransport interface {
	Configure(AudioTransportConfig) error
	Write([]int16) (int, error)
	Read([]int16) (int, error)
}

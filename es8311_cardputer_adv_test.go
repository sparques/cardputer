//go:build esp32 && cardputer_adv

package cardputer

import "testing"

func TestNormalizeES8311Config(t *testing.T) {
	cfg := normalizeES8311Config(ES8311Config{
		ADCGain: 99,
	})

	if cfg.SampleRate != DefaultES8311Config().SampleRate {
		t.Fatalf("normalizeES8311Config() sample rate = %d, want %d", cfg.SampleRate, DefaultES8311Config().SampleRate)
	}
	if cfg.BitsPerSample != DefaultES8311Config().BitsPerSample {
		t.Fatalf("normalizeES8311Config() bits = %d, want %d", cfg.BitsPerSample, DefaultES8311Config().BitsPerSample)
	}
	if cfg.ADCGain != 7 {
		t.Fatalf("normalizeES8311Config() adc gain = %d, want 7", cfg.ADCGain)
	}
}

func TestAudioTransportConfigFromES8311(t *testing.T) {
	cfg := ES8311Config{
		SampleRate:    44100,
		BitsPerSample: ES8311Resolution24,
		UseMCLK:       true,
	}

	transport := audioTransportConfigFromES8311(cfg)
	if transport.SampleRate != cfg.SampleRate {
		t.Fatalf("transport sample rate = %d, want %d", transport.SampleRate, cfg.SampleRate)
	}
	if transport.BitsPerSample != cfg.BitsPerSample {
		t.Fatalf("transport bits = %d, want %d", transport.BitsPerSample, cfg.BitsPerSample)
	}
	if transport.Channels != 1 {
		t.Fatalf("transport channels = %d, want 1", transport.Channels)
	}
	if transport.UseMCLK != cfg.UseMCLK {
		t.Fatalf("transport UseMCLK = %v, want %v", transport.UseMCLK, cfg.UseMCLK)
	}
}

func TestLookupES8311ClockCoefficient(t *testing.T) {
	coeff, ok := lookupES8311ClockCoefficient(44100, 11289600)
	if !ok {
		t.Fatal("lookupES8311ClockCoefficient() returned !ok for 44.1kHz")
	}
	if coeff.rate != 44100 || coeff.mclk != 11289600 {
		t.Fatalf("lookupES8311ClockCoefficient() = %+v, want rate=44100 mclk=11289600", coeff)
	}

	if _, ok := lookupES8311ClockCoefficient(12345, 999999); ok {
		t.Fatal("lookupES8311ClockCoefficient() unexpectedly found unsupported rate")
	}
}

func TestES8311ResolutionBits(t *testing.T) {
	tests := []struct {
		resolution ES8311Resolution
		want       uint8
	}{
		{ES8311Resolution16, 3 << 2},
		{ES8311Resolution18, 2 << 2},
		{ES8311Resolution20, 1 << 2},
		{ES8311Resolution24, 0 << 2},
		{ES8311Resolution32, 4 << 2},
	}

	for _, tc := range tests {
		if got := es8311ResolutionBits(tc.resolution); got != tc.want {
			t.Fatalf("es8311ResolutionBits(%d) = %d, want %d", tc.resolution, got, tc.want)
		}
	}
}

func TestPercentToES8311Volume(t *testing.T) {
	if got := percentToES8311Volume(0); got != es8311VolumeMute {
		t.Fatalf("percentToES8311Volume(0) = %d, want %d", got, es8311VolumeMute)
	}
	if got := percentToES8311Volume(100); got != es8311Volume0dB {
		t.Fatalf("percentToES8311Volume(100) = %d, want %d", got, es8311Volume0dB)
	}
	if got := percentToES8311Volume(200); got != es8311Volume0dB {
		t.Fatalf("percentToES8311Volume(200) = %d, want %d", got, es8311Volume0dB)
	}
}

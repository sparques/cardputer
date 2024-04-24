//go:build rp2040

package cardputer

import (
	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
)

/*
	The cardputer uses an an I2S speaker module, the NS4186.
	It supports 8kHz~96kHz sample rate.
*/

var Speaker = &speaker{}

type speaker struct {
	*piolib.I2S
}

func (spk *speaker) Init() error {
	// piolib only supports speakers right now, so only config that
	var err error
	sm, _ := pio.PIO0.ClaimStateMachine()

	spk.I2S, err = piolib.NewI2S(sm, SpeakerData, I2SClock)
	return err
}

// Beep emits a short, neutrally toned beep. 1/4 second, 440Hz.
func (spk *speaker) Beep() {
	// 8000 Hz sample rate means each sample is 125 us.
	// if we want 1/4 of second, we need 2000 samples
	// check if looping like this causes stuttering; might have to buffer
	// the sample
	for i := 0; i < 2000/len(Sin440FullWave); i++ {
		spk.I2S.WriteMono(Sin440FullWave)
	}
}

// Beep emits a two-toned positive beep
func (spk *speaker) AckBeep() {
}

// Beep emits a two-toned negative beep
func (spk *speaker) NakBeep() {
}

// Sin440FullWave is one full cycle of a 440Hz sine wave, equivalent to 2.25 ms (18 samples * 125us)
var Sin440FullWave = []uint16{0x8000, 0xab5b, 0xd196, 0xee2c, 0xfdbb, 0xfe6c, 0xf02a, 0xd4a5, 0xaf1e, 0x8405, 0x5872, 0x318c, 0x13ed, 0x315, 0x102, 0xdf3, 0x2861, 0x4d2a}

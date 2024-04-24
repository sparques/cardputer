// keypad provides a device driver to scan over the addressing lines connected
// to the 74HC138, tracks what buttons are currently pressed and also can fire
// events for button releases and callbacks. It will also, by default,
// translate keypresses into ANSI terminal codes/characters.
package keypad // import "github.com/sparques/cardputer/keypad"

import (
	"io"
	"machine"
	"time"
)

var (
	// The fatest typist in the word can type (on a full sized, QWERTY Keyboard) just under
	// 1000 characters per minute. That's about 16.67 Hz or a period of 60 ms.
	// So a DefaultScanPeriod of 20 ms should be plenty responsive, right?
	// TODO: Empirical testing.
	// 20ms scan period means two key presses can be registered in 60ms; the first 20ms
	// are needed to detect the press, the next 20 to detect the release, and the final
	// 20 to detect the second press.
	DefaultScanPeriod time.Duration = 20 * time.Millisecond
)

var (
	// DefaultAddressLines are the pins that are connected to the 74HC138
	DefaultAddressLines = [3]machine.Pin{a1, a1, a2}
	// DefaultSenseLines are the GPIO input pins connected to the keypad.
	DefaultSenseLines = [7]machine.Pin{c0, c1, c2, c3, c4, c5, c6}
)

type Device struct {
	// addressLines the pins used to set the address on the 74HC138
	// The indicies should match, addressLines[0] = A0 on the 74HC138 and is equivalent to G8 on the M5 StampC3
	addressLines [3]machine.Pin
	// senseLines are the GPIO input pins connected to the keypad
	// In order, these should be equivalent to G13, G15, G3, G4, G5, G6, G7.
	senseLines [7]machine.Pin
	// state tracks what buttons are currently pressed/released
	state int64
	// buf is a working buffer for what buttons are currently pressed/released
	buf int64
	// scanPeriod is how often to scan over the addressable lines of the keypad
	scanPeriod time.Duration
	// Receiver is an io.Writer interface that will have keypad presses written to
	// as bytes when the EventPressCallback is set to (*Device).WriteByteCallback.
	// Character repeats (holding down button for multiple characters) is not supported.
	// Not every combination of key presses results in a character
	Receiver io.Writer
	// EventPressCallback is called when a button press is detected
	EventPressCallback func(int64)
	// EventReleaseCallback is called a button release is detected
	EventReleaseCallback func(int64)
	// stop is used to signal the goroutine handling scanning the keypad to return.
	stop chan struct{}
}

// New returns a new *Device. New configures the pins as needed.
// A call to (*Device).Start() is need to start the goroutine that
// scans for key presses. (*Device).Stop() can be called to stop this
// goroutine. This is a wrapper to NewWithPins(), called using
// DefaultAddressLines and DefaultSenseLines.
// By default, key presses will be converted into character bytes and
// written to (*Device).Receiver if it is non-nil.
//
// Example:
//
//	kp := keypad.New()
//	kp.Receiver = terminal // write characters to terminal object
//	kp.Start() // start goroutine scanning keypad
func New() (d *Device) {
	return NewWithPins(DefaultAddressLines, DefaultSenseLines)
}

// NewWithPins lets you specify which pins to use.
// addrLines specify the set of pins connected to the 74HC138.
// senseLines are the input pins for detecting key presses.
func NewWithPins(addrLines [3]machine.Pin, senseLines [7]machine.Pin) *Device {
	// configure addrLines for output
	for i := range addrLines {
		addrLines[i].Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	// configure senseLines for input
	for i := range senseLines {
		senseLines[i].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}
	d := &Device{
		addressLines: addrLines,
		senseLines:   senseLines,
		scanPeriod:   20 * time.Microsecond,
		Receiver:     io.Discard,
	}

	// by default, use our own callback that writes bytes to (*Device).Receiver
	d.EventPressCallback = d.WriteByteCallback

	return d
}

// Start enables interrupts on the sense lines and starts the background process of swapping between
// address lines
func (d *Device) Start() {
	scanSenseLines := func() {
		for i := range d.senseLines {
			if d.senseLines[i].Get() {
				d.buf |= 1 << i
			}
		}
	}

	d.stop = make(chan struct{})

	go func() {
		// yes, there's a cleaner way to program this, but this is running all the time
		// in a tight loop, we want to keep it efficient, so no for loops or a function
		// that converts an integer into 3 calls to addressLines[x].High/Low().
		// The for loop in scanSenseLines SHOULD get unrolled by the compiler.

		ticker := time.NewTicker(d.scanPeriod)
		for {
			select {
			case <-d.stop:
				ticker.Stop()
				d.stop = nil
				return
			case <-ticker.C:
				// read all the senseLines
				// The pattern used is to minimize number of calls used to change addressLines
				// Yes it makes the scan codes kind of random and out of order, but the pattern
				// implemented by the hardware doesn't follow a nice in-order, row-by-row pattern anyway

				d.buf = 0
				// 110
				d.addressLines[0].High()
				d.addressLines[1].High()
				d.addressLines[2].Low()

				scanSenseLines()

				d.buf <<= 7
				// 010
				d.addressLines[0].Low()

				scanSenseLines()

				d.buf <<= 7
				// 011
				d.addressLines[2].High()

				scanSenseLines()

				d.buf <<= 7
				// 111
				d.addressLines[0].High()

				scanSenseLines()

				d.buf <<= 7
				// 101
				d.addressLines[1].Low()

				scanSenseLines()

				d.buf <<= 7
				// 100
				d.addressLines[2].Low()

				scanSenseLines()

				d.buf <<= 7
				// 000
				d.addressLines[0].Low()

				scanSenseLines()

				d.buf <<= 7
				// 001
				d.addressLines[2].High()

				scanSenseLines()

				// no difference? don't need to do anything
				if d.buf == d.state {
					continue
				}

				// a button was pressed or released
				r := released(d.state, d.buf)
				p := pressed(d.state, d.buf)
				d.state = d.buf

				if r != 0 && d.EventReleaseCallback != nil {
					d.EventReleaseCallback(r)
				}

				if p != 0 && d.EventPressCallback != nil {
					d.EventPressCallback(p)
				}
			}
		}
	}()
}

// Stop() stops the goroutine scanning for keypad input if it's running.
func (d *Device) Stop() {
	if d.stop != nil {
		d.stop <- struct{}{}
	}
}

// WriteByteCallback is a callback that generates bytes like an ANSI-terminal.
// For example, Ctrl-H generates a \b (byte 0x08), and enter generates
// a \n (byte 0x0A).
// If you set (*Device).EventPressCallback to WriteByteCallback, Key presses will write bytes to (*Device).Receiver.
// By default Receiver is set to io.Discard.
func (d *Device) WriteByteCallback(int64) {
	// we already have access to d.state, so no need to even consider the argument passed or track what was previously

	// Alt is just whatever was pressed prefixed with ESC (0x1B). So mask off Alt here.
	b, ok := ScancodeToBytes[d.state&^BtnAlt]
	// didn't get an actual byte? which is possible if some alt/ctrl/opt key combo is used
	if !ok {
		return
	}

	// If alt is currently pressed, prefix our write with 0x1B.
	if (d.state & BtnAlt) == BtnAlt {
		b = append([]byte{0x1b}, b...)
	}
	d.Receiver.Write(b)

}

// pressed, given bytes a and b, returns the ones present in b that weren't present in a.
// Zeroes that stay zero and ones that stay one are ignored. That is, capture transitions from 0->1
func pressed(a, b int64) int64 {
	return b - (a & b)
}

// released returns the ones that are present in a and nolonger present in b.
// that is for a->b it returns which bits went from 1->0.
func released(a, b int64) int64 {
	return a - (a & b)
}

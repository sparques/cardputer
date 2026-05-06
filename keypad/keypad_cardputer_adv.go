//go:build cardputer_adv

package keypad // import "github.com/sparques/cardputer/keypad"

import (
	"io"
	"machine"
	"time"

	"github.com/sparques/cardputer/internal/adv"
)

const (
	// DefaultScanPeriod is how often the Adv keypad interrupt line is polled.
	DefaultScanPeriod = 20 * time.Millisecond

	keypadIRQ = machine.GPIO11
	keypadSDA = machine.GPIO8
	keypadSCL = machine.GPIO9
)

// Device reads keyboard events from the Cardputer-Adv TCA8418 controller and tracks button state.
type Device struct {
	// state tracks the currently pressed button bitmask.
	state int64
	// scanPeriod controls how often the interrupt line is sampled.
	scanPeriod time.Duration
	// Receiver receives translated byte output when WriteByteCallback is used.
	Receiver io.Writer
	// EventPressCallback is called when one or more buttons become pressed.
	EventPressCallback func(int64)
	// EventReleaseCallback is called when one or more buttons become released.
	EventReleaseCallback func(int64)

	stop    chan struct{}
	irq     machine.Pin
	bus     *machine.I2C
	ctrl    *tca8418
	initErr error
}

// New constructs a Device using the Cardputer-Adv shared I2C bus and keypad IRQ pin.
func New() *Device {
	d := &Device{
		scanPeriod: DefaultScanPeriod,
		Receiver:   io.Discard,
		irq:        keypadIRQ,
	}
	d.EventPressCallback = d.WriteByteCallback

	d.irq.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	d.bus, d.initErr = adv.SharedI2C()
	if d.initErr != nil {
		return d
	}

	d.ctrl = newTCA8418(d.bus, tca8418Addr)
	if err := d.ctrl.configureMatrix(7, 8); err != nil {
		d.initErr = err
		return d
	}
	if err := d.ctrl.flush(); err != nil {
		d.initErr = err
		return d
	}
	d.initErr = d.ctrl.enableKeyInterrupts()
	return d
}

// Start begins polling the TCA8418 interrupt line and draining queued key events.
func (d *Device) Start() {
	if d.stop != nil || d.initErr != nil || d.ctrl == nil {
		return
	}

	d.stop = make(chan struct{})
	go func() {
		ticker := time.NewTicker(d.scanPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-d.stop:
				d.stop = nil
				return
			case <-ticker.C:
				if d.irq.Get() {
					continue
				}
				d.drainEvents()
			}
		}
	}()
}

// Stop stops the background keypad polling loop if it is running.
func (d *Device) Stop() {
	if d.stop != nil {
		d.stop <- struct{}{}
	}
}

// WriteByteCallback translates the current button state into bytes using ScancodeToBytes
// and writes them to Receiver.
func (d *Device) WriteByteCallback(int64) {
	b, ok := ScancodeToBytes[d.state&^BtnAlt]
	if !ok {
		return
	}
	if (d.state & BtnAlt) == BtnAlt {
		b = append([]byte{0x1b}, b...)
	}
	d.Receiver.Write(b)
}

func (d *Device) drainEvents() {
	for {
		n, err := d.ctrl.eventCount()
		if err != nil || n == 0 {
			break
		}

		for i := uint8(0); i < n; i++ {
			event, err := d.ctrl.getEvent()
			if err != nil || event == 0 {
				break
			}
			d.applyEvent(event)
		}

		if d.ctrl.clearKeyInterrupt() != nil {
			return
		}
		if d.irq.Get() {
			return
		}
	}
}

func (d *Device) applyEvent(event uint8) {
	pressed := (event & 0x80) != 0
	key := int(event&0x7F) - 1
	if key < 0 {
		return
	}

	rawRow := key / 10
	rawCol := key % 10
	row, col, ok := remapTCA8418(rawRow, rawCol)
	if !ok {
		return
	}

	mask := buttonMask(row, col)
	if mask == 0 {
		return
	}

	prev := d.state
	if pressed {
		d.state |= mask
		p := pressedBits(prev, d.state)
		if p != 0 && d.EventPressCallback != nil {
			d.EventPressCallback(p)
		}
		return
	}

	d.state &^= mask
	r := releasedBits(prev, d.state)
	if r != 0 && d.EventReleaseCallback != nil {
		d.EventReleaseCallback(r)
	}
}

func remapTCA8418(rawRow, rawCol int) (row, col int, ok bool) {
	if rawRow < 0 || rawRow > 6 || rawCol < 0 || rawCol > 7 {
		return 0, 0, false
	}
	col = rawRow * 2
	if rawCol > 3 {
		col++
	}
	row = (rawCol + 4) % 4
	if row < 0 || row > 3 || col < 0 || col > 13 {
		return 0, 0, false
	}
	return row, col, true
}

func buttonMask(row, col int) int64 {
	if row < 0 || row > 3 || col < 0 || col > 13 {
		return 0
	}
	return int64(1) << (row*14 + col)
}

func pressedBits(a, b int64) int64 {
	return b - (a & b)
}

func releasedBits(a, b int64) int64 {
	return a - (a & b)
}

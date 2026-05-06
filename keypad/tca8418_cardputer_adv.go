//go:build cardputer_adv

package keypad

import "machine"

const (
	tca8418Addr = 0x34

	tca8418RegCFG       = 0x01
	tca8418RegIntStat   = 0x02
	tca8418RegKeyLckEC  = 0x03
	tca8418RegKeyEventA = 0x04
	tca8418RegKPGPIO1   = 0x1D
	tca8418RegKPGPIO2   = 0x1E
	tca8418RegKPGPIO3   = 0x1F

	tca8418CFGKeyEventInt = 0x01
	tca8418StatKeyInt     = 0x01
)

type tca8418 struct {
	bus     *machine.I2C
	address uint16
}

func newTCA8418(bus *machine.I2C, address uint16) *tca8418 {
	return &tca8418{bus: bus, address: address}
}

func (t *tca8418) configureMatrix(rows, cols uint8) error {
	var reg1, reg2, reg3 uint8

	for i := uint8(0); i < rows && i < 8; i++ {
		reg1 |= 1 << i
	}
	for i := uint8(0); i < cols && i < 8; i++ {
		reg2 |= 1 << i
	}
	if cols > 8 {
		for i := uint8(8); i < cols && i < 10; i++ {
			reg3 |= 1 << (i - 8)
		}
	}

	if err := t.writeRegister(tca8418RegKPGPIO1, reg1); err != nil {
		return err
	}
	if err := t.writeRegister(tca8418RegKPGPIO2, reg2); err != nil {
		return err
	}
	if err := t.writeRegister(tca8418RegKPGPIO3, reg3); err != nil {
		return err
	}
	return nil
}

func (t *tca8418) enableKeyInterrupts() error {
	return t.writeRegister(tca8418RegCFG, tca8418CFGKeyEventInt)
}

func (t *tca8418) eventCount() (uint8, error) {
	v, err := t.readRegister(tca8418RegKeyLckEC)
	if err != nil {
		return 0, err
	}
	return v & 0x0F, nil
}

func (t *tca8418) getEvent() (uint8, error) {
	return t.readRegister(tca8418RegKeyEventA)
}

func (t *tca8418) flush() error {
	for {
		n, err := t.eventCount()
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := t.getEvent(); err != nil {
			return err
		}
	}
	return t.clearKeyInterrupt()
}

func (t *tca8418) clearKeyInterrupt() error {
	return t.writeRegister(tca8418RegIntStat, tca8418StatKeyInt)
}

func (t *tca8418) readRegister(reg uint8) (uint8, error) {
	buf := []byte{reg}
	out := []byte{0}
	if err := t.bus.Tx(t.address, buf, out); err != nil {
		return 0, err
	}
	return out[0], nil
}

func (t *tca8418) writeRegister(reg, value uint8) error {
	return t.bus.Tx(t.address, []byte{reg, value}, nil)
}

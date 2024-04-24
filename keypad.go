package cardputer

import (
	"github.com/sparques/cardputer/keypad"
)

// KP initializes the keypad. set KP.Receiver to get key-press events.
var KP = &kp{
	Device: keypad.New(),
}

type kp struct {
	*keypad.Device
}

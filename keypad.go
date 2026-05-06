package cardputer

import (
	"github.com/sparques/cardputer/keypad"
)

// KP exposes the built-in keyboard driver using the board-default wiring.
// Set KP.Receiver or callbacks as needed, then call KP.Start().
var KP = &kp{
	Device: keypad.New(),
}

type kp struct {
	*keypad.Device
}

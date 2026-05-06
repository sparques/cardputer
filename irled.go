package cardputer

import (
	"github.com/sparques/irtrx"
)

/*
	Example: The below is equiavlent to pressing the forward button on a hexbug remote set to channel 1.

	var cmd hexbug.Cmd = hexbug.CH1 | hexbug.CmdFwdMask
	IRLED.SendFrame(cmd)
*/

// IRLED exposes the built-in IR transmitter as an irtrx.TxDevice.
// Currently only 38KHz is supported. To send signals, you can manually send
// on/off pairs or use a value that implements irtrx.FrameMarshaller.
//
// IRLED is initialized during package startup, so it is ready to use immediately.
var IRLED = &irled{
	TxDevice: irtrx.NewTxDevice(IRPin),
}

type irled struct {
	*irtrx.TxDevice
}

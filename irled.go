package cardputer

import (
	"github.com/sparques/irtrx"
)

/*
	Example: The below is equiavlent to pressing the forward button on a hexbug remote set to channel 1.

	var cmd hexbug.Cmd = hexbug.CH1 | hexbug.CmdFwdMask
	IRLED.SendFrame(cmd)
*/

// IRLED composites an irtrx.TxDevice, exposing SendPair(), SendPairs(), SendFrame(), and SendFrames().
// Currently only 38KHz is supported. To send signals, you can manually send on/off pairs using SendPair() or
// use something that implements irtrx.FrameMarshaller. For example, the github.com/sparques/irtrx/hexbug package
// has hexbug.Cmd which implements irtrx.FrameMarshller; with this you can control Hexbug robots.
//
// IRLED is initialized at start up, so it is ready to use immediately.
var IRLED = &irled{
	TxDevice: irtrx.NewTxDevice(IRPin),
}

type irled struct {
	*irtrx.TxDevice
}

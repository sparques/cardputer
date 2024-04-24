package cardputer

import "machine"

var BatteryADC = machine.ADC{Pin: BatterySense}

func init() {
	// TODO: Figure out what settings actually work best here
	BatterySense.Configure(machine.PinConfig{Mode: machine.PinAnalog})
	BatteryADC.Configure(machine.ADCConfig{Resolution: 8, Samples: 32})
}

// BatteryLevel returns the battery level as a percent, 0 to 100.
func BatteryLevel() int {
	// As I understand the docs for the cardputer, the battery voltage is passed to a voltage divider
	// that divides voltage in half. If we expect the fully charged battery to be at 3.7 volts and our ref
	// voltage is 3.3 volts, than our ADC (in 8 bit resolution) will read ~144 at full charge.
	// (3.7/2) / 3.3 * 246 ≈ 144
	// I've rounded down to 143 for the pleasure of actually seeing a 100% charged battery
	// TODO: after confirming this works, add a max() in there so we don't show MORE than 100%
	// TODO: make this measurement linear if it isn't; that is, if the battery follows a logarthmic
	// discharge, correct that to be linear.
	return int(BatteryADC.Get()) * 100 / 144
}

package cardputer

import "machine"

const (
	batteryADCResolutionBits = 12
	batteryADCMax            = (1 << batteryADCResolutionBits) - 1
	batteryADCReferenceMV    = 3300
	batteryDividerNumerator  = 2
	batteryEmptyMV           = 3300
	batteryFullMV            = 4200
)

// BatteryADC exposes the ADC channel used to sample the battery divider.
// It is configured during package initialization.
var BatteryADC = machine.ADC{Pin: BatterySense}

func init() {
	BatterySense.Configure(machine.PinConfig{Mode: machine.PinAnalog})
	BatteryADC.Configure(machine.ADCConfig{Resolution: batteryADCResolutionBits, Samples: 32})
}

// BatteryMilliVolts returns the measured battery voltage in millivolts.
// This is derived from the ADC reading and the board's 2:1 divider.
func BatteryMilliVolts() int {
	raw := int(BatteryADC.Get())
	return raw * batteryADCReferenceMV * batteryDividerNumerator / batteryADCMax
}

// BatteryLevel returns an approximate battery percentage in the range 0..100.
func BatteryLevel() int {
	mv := BatteryMilliVolts()
	switch {
	case mv <= batteryEmptyMV:
		return 0
	case mv >= batteryFullMV:
		return 100
	default:
		return (mv - batteryEmptyMV) * 100 / (batteryFullMV - batteryEmptyMV)
	}
}

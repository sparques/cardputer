//go:build (esp32 || esp32s3) && cardputer_adv

package adv

import "machine"

const (
	SharedI2CSDA = machine.GPIO8
	SharedI2CSCL = machine.GPIO9
)

var (
	sharedI2C           = machine.I2C0
	sharedI2CConfigured bool
)

// SharedI2C returns the Cardputer-Adv shared peripheral bus used by the
// keyboard controller, ES8311 codec, and IMU.
func SharedI2C() (*machine.I2C, error) {
	if sharedI2CConfigured {
		return sharedI2C, nil
	}

	err := sharedI2C.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       SharedI2CSDA,
		SCL:       SharedI2CSCL,
	})
	if err != nil {
		return nil, err
	}

	sharedI2CConfigured = true
	return sharedI2C, nil
}

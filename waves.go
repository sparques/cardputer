//go:build waves

package main

import (
	"fmt"
	"math"
)

func main() {
	samplerate := float64(8000)
	freq := float64(440)
	period := 1 / freq

	amplitude := float64(^uint16(0)) / 2
	center := float64(^uint16(0)) / 2

	fmt.Printf("// Samplerate: \t%f\n", samplerate)
	fmt.Printf("// Frequency: \t%f\n", freq)
	fmt.Printf("// Amplitude: \t%f\n", amplitude)
	fmt.Printf("// center: \t%f\n", center)

	samples := make([]uint16, int(math.Round(samplerate*period)))
	for i := range samples {

		samples[i] = uint16(math.Round(amplitude*math.Sin(2*math.Pi*freq*float64(i)/samplerate) + center))
	}
	fmt.Printf("// SinFullWave: \t%+v\n", samples)

	fmt.Printf("Sin%dFullWave := %#v\n", int(freq), samples)
}

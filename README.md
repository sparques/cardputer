[![Coverage Status](https://coveralls.io/repos/github/sparques/cardputer/badge.svg?branch=master)](https://coveralls.io/github/sparques/cardputer?branch=master)
[![Go ReportCard](https://goreportcard.com/badge/sparques/cardputer)](https://goreportcard.com/report/sparques/cardputer)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/sparques/cardputer)

# Cardputer

A TinyGo package for working with the hardware on the [M5 Cardputer](https://shop.m5stack.com/products/m5stack-cardputer-kit-w-m5stamps3).

This package targets the ESP32-S3 based Cardputer family hardware.

Build with `make build` for the original Cardputer, or `make build BOARD=cardputer-adv` for the Cardputer-Adv. The checked-in `Makefile` forces TinyGo 0.37 to use Go 1.24.6, because TinyGo currently rejects the system `go1.26.x` toolchain.

The original Cardputer remains the default build. Cardputer-Adv support currently shares the common peripherals and board pin map, but still needs board-specific keypad and audio drivers on top of the new build-tag split.

# Project Goals

I want to be able to use everything the cardputer has to offer, including the microphone and amplifier/speaker.

Much of the hardware on the cardputer already has drivers. The [display](https://github.com/tinygo-org/drivers/tree/release/st7789) and the [SD card](https://github.com/tinygo-org/drivers/tree/release/sdcard) both already have drivers.

Some of the simple "peripherals" just need some convenience wrappers, such as sensing the battery level.

I intend to use my [IR package](https://github.com/sparques/irtrx) for sending IR signals. I may add an IR receiver via the grove port for receiving signals.

# Progress

 - ☑️ Pin definitions
 - ☑️ IR LED
 - ☑️ Keypad driver
 - 🔄 Screen (just a thin wrapper around the existing st7789 driver)
 - 🔄 SD Card (pins are initialized, but still working on filesystem support)
 - 🔄 Audio support
 - ☑️ Battery  Level

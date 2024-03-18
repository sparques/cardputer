# Cardputer

A TinyGo package for working with the hardware on the [M5 Cardputer](https://shop.m5stack.com/products/m5stack-cardputer-kit-w-m5stamps3).

This should work with the ESP32 / Stamp S3 board as well as the RP2040-based board.

# Project Goals

I want to be able to use everything the cardputer has to offer, including the i2s microphone and amplifier/speaker. The best way to do this is probably figuring out how to get PIO to work under TinyGo (it's possible, per several TinyGo issue discussions) and use the i2s PIO programs from [here](https://github.com/malacalypse/rp2040_i2s_example).

Much of the hardware on the cardputer already has drivers. The [display](https://github.com/tinygo-org/drivers/tree/release/st7789) and the [SD card](https://github.com/tinygo-org/drivers/tree/release/sdcard) both already have drivers.

Some of the simple "peripherals" just need some convenience wrappers, such as sensing the battery level.

I intend to use my [IR package](https://github.com/sparques/irtrx) for sending IR signals. I may add an IR receiver via the grove port for receiving signals.

# Progress

 - keypad driver written
 
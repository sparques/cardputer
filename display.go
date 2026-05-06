package cardputer

import (
	"image"
	"image/color"
	"machine"

	"tinygo.org/x/drivers/st7789"
)

const (
	dispWidth  = 240
	dispHeight = 135
)

// Display provides access to the built-in ST7789 LCD.
// Call Init before using it.
var Display = &display{}

type display struct {
	device st7789.Device
	bus    machine.SPI
	scroll int16
}

func (d *display) Init() {
	// Configure SPI pins
	LCDMOSI.Configure(machine.PinConfig{Mode: machine.PinSPI})
	//LCDMISO.Configure(machine.PinConfig{Mode: machine.PinSPI})

	// configure SPI
	machine.SPI0.Configure(machine.SPIConfig{
		SCK:       LCDSCK,
		SDO:       LCDMOSI,
		Frequency: 40 * machine.MHz,
	})

	d.device = st7789.New(machine.SPI0, LCDReset, LCDRS, LCDCS, LCDBacklight)

	d.device.Configure(st7789.Config{
		Width:  dispWidth,
		Height: dispHeight,
		//Rotation:
		//RowOffset    int16
		//ColumnOffset int16
		//FrameRate    FrameRate
		//VSyncLines   int16
	})

}

func (*display) Bounds() image.Rectangle {
	return image.Rect(0, 0, dispWidth, dispHeight)

	// if the display doesn't use 0,0 as the upper left corner, might have to change this.
	// if 0,0 is bottom left
	// return image.Rect(0,-135,240,0)
}

func (d *display) Set(x, y int, c color.Color) {
	d.device.SetPixel(int16(x), int16(y), colorToRGBA(c))
}

func (d *display) ColorModel() color.Model {
	return color.RGBAModel
}

func (d *display) Fill(r image.Rectangle, c color.Color) {
	d.device.FillRectangle(int16(r.Min.X), int16(r.Min.Y), int16(r.Dx()), int16(r.Dy()), colorToRGBA(c))
}

// Blit copies pixels from img into display, aligning img.Bounds().Min to 'at' within display.
func (d *display) Blit(img image.Image, at image.Point) {
	// must convert img to a slice of []color.RGBA first

	d.device.FillRectangleWithBuffer(int16(at.X), int16(at.Y), int16(img.Bounds().Dx()), int16(img.Bounds().Dy()), imageToRGBASlice(img))

	// would be better to make st7789.Device actually implement Blit and have it iterate over pixel data
	// so we still get hte performance benefit of Blit operation, without the cost of memory
}

// we can't implement this (can't read the display!), but need to fake it so we implement draw.Draw
// we could buffer the Set() calls, but pretty sure we'll run out of RAM.
func (d *display) At(x, y int) color.Color {
	return color.Alpha{0}
}

func (d *display) Scroll(amount int) {
	d.scroll = (d.scroll + int16(amount)) % dispHeight
	if d.scroll < 0 {
		d.scroll += dispHeight
	}
	d.device.SetScroll(d.scroll)
	// and then clear scrolled area? Leave that to fansiterm...
}

func colorToRGBA(c color.Color) color.RGBA {
	if rgba, ok := c.(color.RGBA); ok {
		return rgba
	}
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r / 0x101), uint8(g / 0x101), uint8(b / 0x101), uint8(a / 0x101)}
}

func imageToRGBASlice(img image.Image) []color.RGBA {
	bounds := img.Bounds()
	out := make([]color.RGBA, bounds.Dx()*bounds.Dy())
	stride := bounds.Dx()
	if rgbaImg, ok := img.(*image.RGBA); ok {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				dst := (y-bounds.Min.Y)*stride + (x - bounds.Min.X)
				out[dst] = rgbaImg.RGBAAt(x, y)
			}
		}
		return out
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst := (y-bounds.Min.Y)*stride + (x - bounds.Min.X)
			out[dst] = colorToRGBA(img.At(x, y))
		}
	}

	return out
}

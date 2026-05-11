package cardputer

import (
	"image"
	"image/color"
	"machine"

	"tinygo.org/x/drivers/st7789"
)

const (
	dispWidth      = 240
	dispHeight     = 135
	panelWidth     = 135
	panelHeight    = 240
	panelRowOffset = 40
	panelColOffset = 52
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
	// configure SPI
	machine.SPI0.Configure(machine.SPIConfig{
		SCK:       LCDSCK,
		SDO:       LCDMOSI,
		Frequency: 40 * machine.MHz,
	})

	d.device = st7789.New(machine.SPI0, LCDReset, LCDRS, LCDCS, LCDBacklight)

	d.device.Configure(st7789.Config{
		Width:        panelWidth,
		Height:       panelHeight,
		Rotation:     st7789.ROTATION_270,
		RowOffset:    panelRowOffset,
		ColumnOffset: panelColOffset,
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
	if !image.Pt(x, y).In(d.Bounds()) {
		return
	}
	hwX, hwY := mapLogicalPoint(x, y)
	d.device.SetPixel(hwX, hwY, colorToRGBA(c))
}

func (d *display) ColorModel() color.Model {
	return color.RGBAModel
}

func (d *display) Fill(r image.Rectangle, c color.Color) {
	r = r.Intersect(d.Bounds())
	if r.Empty() {
		return
	}
	hw := mapLogicalRect(r)
	d.device.FillRectangle(int16(hw.Min.X), int16(hw.Min.Y), int16(hw.Dx()), int16(hw.Dy()), colorToRGBA(c))
}

// Blit copies pixels from img into display, aligning img.Bounds().Min to 'at' within display.
func (d *display) Blit(img image.Image, at image.Point) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			d.Set(at.X+x-bounds.Min.X, at.Y+y-bounds.Min.Y, img.At(x, y))
		}
	}
}

// we can't implement this (can't read the display!), but need to fake it so we implement draw.Draw
// we could buffer the Set() calls, but pretty sure we'll run out of RAM.
func (d *display) At(x, y int) color.Color {
	return color.Alpha{0}
}

func (d *display) Scroll(amount int) {
	// Hardware scroll moves in panel coordinates. The display wrapper exposes a
	// transformed logical surface, so scrolling must be handled by the caller.
}

func mapLogicalPoint(x, y int) (int16, int16) {
	return int16(dispWidth - 1 - x), int16(dispHeight - 1 - y)
}

func mapLogicalRect(r image.Rectangle) image.Rectangle {
	return image.Rect(
		dispWidth-r.Max.X,
		dispHeight-r.Max.Y,
		dispWidth-r.Min.X,
		dispHeight-r.Min.Y,
	)
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

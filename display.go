package cardputer

import (
	"image"
	"image/color"
	"machine"
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
	device st7789RGB565
	bus    machine.SPI
	scroll int16
	pix    []RGB565
	line   []RGB565
}

type RGB565 uint16

var RGB565Model = color.ModelFunc(func(c color.Color) color.Color {
	if c, ok := c.(RGB565); ok {
		return c
	}
	return colorToRGB565(c)
})

func (d *display) Init() {
	d.InitWithBuffer(true)
}

func (d *display) InitWithBuffer(buffered bool) {
	// configure SPI
	machine.SPI0.Configure(machine.SPIConfig{
		SCK:       LCDSCK,
		SDO:       LCDMOSI,
		Frequency: 40 * machine.MHz,
	})

	d.device = newST7789RGB565(machine.SPI0, LCDReset, LCDRS, LCDCS, LCDBacklight)

	d.device.Configure(st7789Config{
		Width:        panelWidth,
		Height:       panelHeight,
		Rotation:     st7789Rotation270,
		RowOffset:    panelRowOffset,
		ColumnOffset: panelColOffset,
		Buffered:     false,
		//FrameRate    FrameRate
		//VSyncLines   int16
	})

	if buffered {
		d.pix = make([]RGB565, dispWidth*dispHeight)
		d.line = make([]RGB565, dispWidth)
	} else {
		d.pix = nil
		d.line = nil
	}
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
	p := colorToRGB565(c)
	if d.pix != nil {
		d.pix[d.pixOffset(x, y)] = p
	}
	hwX, hwY := mapLogicalPoint(x, y)
	d.device.Set(int(hwX), int(hwY), p)
}

func (d *display) ColorModel() color.Model {
	return RGB565Model
}

func (d *display) Fill(r image.Rectangle, c color.Color) {
	r = r.Intersect(d.Bounds())
	if r.Empty() {
		return
	}
	p := colorToRGB565(c)
	if d.pix != nil {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			row := y * dispWidth
			for x := r.Min.X; x < r.Max.X; x++ {
				d.pix[row+x] = p
			}
		}
	}
	hw := mapLogicalRect(r)
	d.device.Fill(hw, p)
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

func (d *display) At(x, y int) color.Color {
	if d.pix != nil && image.Pt(x, y).In(d.Bounds()) {
		return d.pix[d.pixOffset(x, y)]
	}
	return color.Alpha{0}
}

func (d *display) Scroll(amount int) {
	d.RegionScroll(d.Bounds(), amount)
}

func (d *display) RegionScroll(region image.Rectangle, amount int) {
	if d.pix == nil {
		return
	}
	region = region.Intersect(d.Bounds())
	if region.Empty() || amount == 0 {
		return
	}
	height := region.Dy()
	if amount >= height || amount <= -height {
		return
	}

	if amount > 0 {
		for y := region.Min.Y; y < region.Max.Y-amount; y++ {
			dst := d.pixOffset(region.Min.X, y)
			src := d.pixOffset(region.Min.X, y+amount)
			copy(d.pix[dst:dst+region.Dx()], d.pix[src:src+region.Dx()])
		}
	} else {
		amount = -amount
		for y := region.Max.Y - 1; y >= region.Min.Y+amount; y-- {
			dst := d.pixOffset(region.Min.X, y)
			src := d.pixOffset(region.Min.X, y-amount)
			copy(d.pix[dst:dst+region.Dx()], d.pix[src:src+region.Dx()])
		}
	}
	d.flush(region)
}

func (d *display) flush(r image.Rectangle) {
	r = r.Intersect(d.Bounds())
	if r.Empty() || d.pix == nil {
		return
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			d.line[r.Dx()-1-(x-r.Min.X)] = d.pix[d.pixOffset(x, y)]
		}
		hw := mapLogicalRect(image.Rect(r.Min.X, y, r.Max.X, y+1))
		d.device.Draw(hw, d.line[:r.Dx()])
	}
}

func (*display) pixOffset(x, y int) int {
	return y*dispWidth + x
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

func colorToRGB565(c color.Color) RGB565 {
	if c, ok := c.(RGB565); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return NewRGB565(uint8(r/0x101), uint8(g/0x101), uint8(b/0x101))
}

func NewRGB565(r, g, b uint8) RGB565 {
	return RGB565(uint16(r&0xf8)<<8 | uint16(g&0xfc)<<3 | uint16(b)>>3)
}

func (c RGB565) Bytes() (hi, lo byte) {
	return byte(c >> 8), byte(c)
}

func (c RGB565) RGBA8() color.RGBA {
	r := uint8((c >> 11) & 0x1f)
	g := uint8((c >> 5) & 0x3f)
	b := uint8(c & 0x1f)
	return color.RGBA{
		R: (r << 3) | (r >> 2),
		G: (g << 2) | (g >> 4),
		B: (b << 3) | (b >> 2),
		A: 0xff,
	}
}

func (c RGB565) RGBA() (r, g, b, a uint32) {
	rgba := c.RGBA8()
	return uint32(rgba.R) * 0x101, uint32(rgba.G) * 0x101, uint32(rgba.B) * 0x101, 0xffff
}

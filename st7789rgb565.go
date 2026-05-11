package cardputer

import (
	"errors"
	"image"
	"machine"
	"time"

	"tinygo.org/x/drivers"
)

type st7789Rotation = drivers.Rotation

const (
	st7789Rotation0   = drivers.Rotation0
	st7789Rotation90  = drivers.Rotation90
	st7789Rotation180 = drivers.Rotation180
	st7789Rotation270 = drivers.Rotation270
)

type st7789FrameRate uint8

type st7789Config struct {
	Width        int16
	Height       int16
	Rotation     st7789Rotation
	RowOffset    int16
	ColumnOffset int16
	FrameRate    st7789FrameRate
	VSyncLines   int16
	PVGAMCTRL    []uint8
	NVGAMCTRL    []uint8
	Buffered     bool
}

type st7789RGB565 struct {
	bus             drivers.SPI
	dcPin           machine.Pin
	resetPin        machine.Pin
	csPin           machine.Pin
	blPin           machine.Pin
	width           int16
	height          int16
	columnOffsetCfg int16
	rowOffsetCfg    int16
	columnOffset    int16
	rowOffset       int16
	rotation        st7789Rotation
	frameRate       st7789FrameRate
	vSyncLines      int16
	pix             []RGB565
	tx              []byte
	cmdBuf          [1]byte
	buf             [6]byte
}

const (
	st7789SWRESET   = 0x01
	st7789SLPOUT    = 0x11
	st7789NORON     = 0x13
	st7789INVON     = 0x21
	st7789DISPON    = 0x29
	st7789CASET     = 0x2a
	st7789RASET     = 0x2b
	st7789RAMWR     = 0x2c
	st7789COLMOD    = 0x3a
	st7789MADCTL    = 0x36
	st7789MADCTLMY  = 0x80
	st7789MADCTLMX  = 0x40
	st7789MADCTLMV  = 0x20
	st7789MADCTLBGR = 0x08
	st7789FRCTRL2   = 0xc6
	st7789PORCTRL   = 0xb2
	st7789GMCTRP1   = 0xe0
	st7789GMCTRN1   = 0xe1

	st7789ColorRGB565 st7789ColorFormat = 0b101
	st7789FrameRate60 st7789FrameRate   = 0x0f
)

type st7789ColorFormat uint8

func newST7789RGB565(bus drivers.SPI, resetPin, dcPin, csPin, blPin machine.Pin) st7789RGB565 {
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	resetPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	blPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return st7789RGB565{
		bus:      bus,
		dcPin:    dcPin,
		resetPin: resetPin,
		csPin:    csPin,
		blPin:    blPin,
	}
}

func (d *st7789RGB565) Configure(cfg st7789Config) {
	if cfg.Width != 0 {
		d.width = cfg.Width
	} else {
		d.width = 240
	}
	if cfg.Height != 0 {
		d.height = cfg.Height
	} else {
		d.height = 240
	}
	d.rotation = cfg.Rotation
	d.rowOffsetCfg = cfg.RowOffset
	d.columnOffsetCfg = cfg.ColumnOffset
	if cfg.FrameRate != 0 {
		d.frameRate = cfg.FrameRate
	} else {
		d.frameRate = st7789FrameRate60
	}
	if cfg.VSyncLines >= 2 && cfg.VSyncLines <= 254 {
		d.vSyncLines = cfg.VSyncLines
	} else {
		d.vSyncLines = 16
	}

	w, h := d.Size()
	if cfg.Buffered {
		d.pix = make([]RGB565, int(w)*int(h))
	} else {
		d.pix = nil
	}
	if d.tx == nil {
		d.tx = make([]byte, 512)
	}

	d.resetPin.High()
	time.Sleep(50 * time.Millisecond)
	d.resetPin.Low()
	time.Sleep(50 * time.Millisecond)
	d.resetPin.High()
	time.Sleep(50 * time.Millisecond)

	d.startWrite()
	d.sendCommand(st7789SWRESET, nil)
	d.endWrite()
	time.Sleep(150 * time.Millisecond)
	d.startWrite()
	d.sendCommand(st7789SLPOUT, nil)
	d.setColorFormat(st7789ColorRGB565)
	time.Sleep(10 * time.Millisecond)
	d.setRotation(d.rotation)
	d.setWindow(0, 0, w, h)
	d.endWrite()
	d.Fill(image.Rect(0, 0, int(w), int(h)), RGB565(0))

	d.startWrite()
	d.sendCommand(st7789FRCTRL2, []byte{byte(d.frameRate)})
	fp := uint8(d.vSyncLines / 2)
	bp := uint8(d.vSyncLines - int16(fp))
	d.sendCommand(st7789PORCTRL, []byte{bp, fp, 0x00, 0x22, 0x22})
	d.sendCommand(st7789INVON, nil)
	time.Sleep(10 * time.Millisecond)
	if len(cfg.PVGAMCTRL) == 14 {
		d.sendCommand(st7789GMCTRP1, cfg.PVGAMCTRL)
	}
	if len(cfg.NVGAMCTRL) == 14 {
		d.sendCommand(st7789GMCTRN1, cfg.NVGAMCTRL)
	}
	d.sendCommand(st7789NORON, nil)
	time.Sleep(10 * time.Millisecond)
	d.sendCommand(st7789DISPON, nil)
	time.Sleep(10 * time.Millisecond)
	d.endWrite()
	d.blPin.High()
}

func (d *st7789RGB565) Bounds() image.Rectangle {
	w, h := d.Size()
	return image.Rect(0, 0, int(w), int(h))
}

func (d *st7789RGB565) Size() (w, h int16) {
	if d.rotation == st7789Rotation0 || d.rotation == st7789Rotation180 {
		return d.width, d.height
	}
	return d.height, d.width
}

func (d *st7789RGB565) Set(x, y int, c RGB565) {
	if !image.Pt(x, y).In(d.Bounds()) {
		return
	}
	if d.pix != nil {
		d.pix[d.pixOffset(x, y)] = c
	}
	d.startWrite()
	d.setWindow(int16(x), int16(y), 1, 1)
	d.txRGB565(c)
	d.endWrite()
}

func (d *st7789RGB565) At(x, y int) (RGB565, bool) {
	if d.pix == nil || !image.Pt(x, y).In(d.Bounds()) {
		return 0, false
	}
	return d.pix[d.pixOffset(x, y)], true
}

func (d *st7789RGB565) Fill(r image.Rectangle, c RGB565) error {
	r = r.Intersect(d.Bounds())
	if r.Empty() {
		return nil
	}
	if d.pix != nil {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			offset := d.pixOffset(r.Min.X, y)
			for x := r.Min.X; x < r.Max.X; x++ {
				d.pix[offset] = c
				offset++
			}
		}
	}
	d.startWrite()
	err := d.fill(r, c)
	d.endWrite()
	return err
}

func (d *st7789RGB565) Draw(r image.Rectangle, pix []RGB565) error {
	r = r.Intersect(d.Bounds())
	if r.Empty() {
		return nil
	}
	if len(pix) != r.Dx()*r.Dy() {
		return errors.New("buffer length does not match rectangle size")
	}
	if d.pix != nil {
		src := 0
		for y := r.Min.Y; y < r.Max.Y; y++ {
			dst := d.pixOffset(r.Min.X, y)
			copy(d.pix[dst:dst+r.Dx()], pix[src:src+r.Dx()])
			src += r.Dx()
		}
	}
	d.startWrite()
	d.setWindow(int16(r.Min.X), int16(r.Min.Y), int16(r.Dx()), int16(r.Dy()))
	d.txRGB565Slice(pix)
	d.endWrite()
	return nil
}

func (d *st7789RGB565) sendCommand(command uint8, data []byte) error {
	d.cmdBuf[0] = command
	d.dcPin.Low()
	err := d.bus.Tx(d.cmdBuf[:1], nil)
	d.dcPin.High()
	if len(data) != 0 {
		err = d.bus.Tx(data, nil)
	}
	return err
}

func (d *st7789RGB565) startWrite() {
	if d.csPin != machine.NoPin {
		d.csPin.Low()
	}
}

func (d *st7789RGB565) endWrite() {
	if d.csPin != machine.NoPin {
		d.csPin.High()
	}
}

func (d *st7789RGB565) setWindow(x, y, w, h int16) {
	x += d.columnOffset
	y += d.rowOffset
	copy(d.buf[:4], []uint8{uint8(x >> 8), uint8(x), uint8((x + w - 1) >> 8), uint8(x + w - 1)})
	d.sendCommand(st7789CASET, d.buf[:4])
	copy(d.buf[:4], []uint8{uint8(y >> 8), uint8(y), uint8((y + h - 1) >> 8), uint8(y + h - 1)})
	d.sendCommand(st7789RASET, d.buf[:4])
	d.sendCommand(st7789RAMWR, nil)
}

func (d *st7789RGB565) fill(r image.Rectangle, c RGB565) error {
	if r.Empty() {
		return nil
	}
	d.setWindow(int16(r.Min.X), int16(r.Min.Y), int16(r.Dx()), int16(r.Dy()))
	count := r.Dx() * r.Dy()
	hi, lo := c.Bytes()
	for count > 0 {
		pixels := len(d.tx) / 2
		if pixels > count {
			pixels = count
		}
		for i := 0; i < pixels; i++ {
			d.tx[i*2] = hi
			d.tx[i*2+1] = lo
		}
		d.bus.Tx(d.tx[:pixels*2], nil)
		count -= pixels
	}
	return nil
}

func (d *st7789RGB565) txRGB565(c RGB565) {
	hi, lo := c.Bytes()
	d.buf[0] = hi
	d.buf[1] = lo
	d.bus.Tx(d.buf[:2], nil)
}

func (d *st7789RGB565) txRGB565Slice(pix []RGB565) {
	for len(pix) > 0 {
		count := len(d.tx) / 2
		if count > len(pix) {
			count = len(pix)
		}
		for i, c := range pix[:count] {
			hi, lo := c.Bytes()
			d.tx[i*2] = hi
			d.tx[i*2+1] = lo
		}
		d.bus.Tx(d.tx[:count*2], nil)
		pix = pix[count:]
	}
}

func (d *st7789RGB565) setColorFormat(format st7789ColorFormat) {
	d.sendCommand(st7789COLMOD, []byte{byte(format) | 0x50})
}

func (d *st7789RGB565) setRotation(rotation st7789Rotation) error {
	madctl := uint8(0)
	switch rotation % 4 {
	case st7789Rotation0:
		d.rowOffset = 0
		d.columnOffset = 0
	case st7789Rotation90:
		madctl = st7789MADCTLMX | st7789MADCTLMV
		d.rowOffset = 0
		d.columnOffset = 0
	case st7789Rotation180:
		madctl = st7789MADCTLMX | st7789MADCTLMY
		d.rowOffset = d.rowOffsetCfg
		d.columnOffset = d.columnOffsetCfg
	case st7789Rotation270:
		madctl = st7789MADCTLMY | st7789MADCTLMV
		d.rowOffset = d.columnOffsetCfg
		d.columnOffset = d.rowOffsetCfg
	}
	return d.sendCommand(st7789MADCTL, []byte{madctl})
}

func (d *st7789RGB565) pixOffset(x, y int) int {
	w, _ := d.Size()
	return y*int(w) + x
}

package vnc

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"

	"github.com/golang/glog"
	vnclib "github.com/kward/go-vnc"
	"github.com/kward/venue/venuelib"
)

// Framebuffer maintains a local copy of the remote VNC image.
type Framebuffer struct {
	Width, Height int
	fb            *image.RGBA
}

// NewFramebuffer returns a new Framebuffer object.
func NewFramebuffer(w, h int) *Framebuffer {
	return &Framebuffer{
		w, h,
		image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}}),
	}
}

// Paint accepts a Rectangle and Color data, and paints the framebuffer with it.
func (f *Framebuffer) Paint(r vnclib.Rectangle, colors []vnclib.Color) {
	if glog.V(4) {
		glog.Info(venuelib.FnNameWithArgs(r.String(), "colors"))
	}
	// TODO(kward): Implement double or triple buffering to reduce paint
	// interference.
	for x := 0; x < int(r.Width); x++ {
		for y := 0; y < int(r.Height); y++ {
			c := colors[x+y*int(r.Width)]
			f.fb.SetRGBA(x+int(r.X), y+int(r.Y), color.RGBA{uint8(c.R), uint8(c.G), uint8(c.B), 255})
		}
	}
}

// PNG converts the framebuffer into a base64 encoded PNG string.
func (f *Framebuffer) PNG() (string, error) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	var buf bytes.Buffer
	err := png.Encode(&buf, f.fb)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

package vnc

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"log"

	vnclib "github.com/kward/go-vnc"
)

type Framebuffer struct {
	Width, Height int
	fb            *image.RGBA
}

func NewFramebuffer(w, h int) *Framebuffer {
	return &Framebuffer{
		w, h,
		image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}}),
	}
}

func (f *Framebuffer) Paint(r vnclib.Rectangle, colors []vnclib.Color) {
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
	var buf bytes.Buffer
	err := png.Encode(&buf, f.fb)
	if err != nil {
		log.Println("Frambuffer.PNG() error;", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

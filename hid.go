package venue

import (
	"image"
	"time"

	vnc "github.com/kward/go-vnc"
)

func (v *Venue) MouseMove(p image.Point) {
	v.conn.PointerEvent(vnc.ButtonNone, uint16(p.X), uint16(p.Y))
	time.Sleep(uiSettle) // Give mouse some time to "settle".
}

func (v *Venue) MouseLeftClick(p image.Point) {
	v.MouseMove(p)
	v.conn.PointerEvent(vnc.ButtonLeft, uint16(p.X), uint16(p.Y))
	v.conn.PointerEvent(vnc.ButtonNone, uint16(p.X), uint16(p.Y))
}

func (v *Venue) MouseDrag(p, d image.Point) {
	v.MouseMove(p)
	v.conn.PointerEvent(vnc.ButtonLeft, uint16(p.X), uint16(p.Y))
	p = p.Add(d) // Add delta.
	v.conn.PointerEvent(vnc.ButtonLeft, uint16(p.X), uint16(p.Y))
	time.Sleep(uiSettle) // Give mouse some time to "settle".
	v.conn.PointerEvent(vnc.ButtonNone, uint16(p.X), uint16(p.Y))
}

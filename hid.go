package venue

import (
	"image"
	"log"
	"time"

	vnc "github.com/kward/go-vnc"
)

// KeyPress presses a key on the VENUE console.
func (v *Venue) KeyPress(key uint32) error {
	log.Printf("KeyPress key=0x%x\n", key)
	if err := v.conn.KeyEvent(key, true); err != nil {
		return err
	}
	if err := v.conn.KeyEvent(key, false); err != nil {
		return err
	}
	return nil
}

// MouseMove moves the mouse.
func (v *Venue) MouseMove(p image.Point) {
	v.conn.PointerEvent(vnc.ButtonNone, uint16(p.X), uint16(p.Y))
	time.Sleep(uiSettle) // Give mouse some time to "settle".
}

// MouseLeftClick moves the mouse to a position and left clicks.
func (v *Venue) MouseLeftClick(p image.Point) {
	v.MouseMove(p)
	v.conn.PointerEvent(vnc.ButtonLeft, uint16(p.X), uint16(p.Y))
	v.conn.PointerEvent(vnc.ButtonNone, uint16(p.X), uint16(p.Y))
}

// MouseDrag moves the mouse, clicks, and drags to a new position.
func (v *Venue) MouseDrag(p, d image.Point) {
	v.MouseMove(p)
	v.conn.PointerEvent(vnc.ButtonLeft, uint16(p.X), uint16(p.Y))
	p = p.Add(d) // Add delta.
	v.conn.PointerEvent(vnc.ButtonLeft, uint16(p.X), uint16(p.Y))
	time.Sleep(uiSettle) // Give mouse some time to "settle".
	v.conn.PointerEvent(vnc.ButtonNone, uint16(p.X), uint16(p.Y))
}

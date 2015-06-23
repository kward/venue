package venue

import "log"

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

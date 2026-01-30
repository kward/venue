package venue

import (
	"testing"

	"github.com/kward/venue/internal/router/signals"
)

func TestSignalControlName(t *testing.T) {
	for _, tt := range []struct {
		sig   signals.Signal
		sigNo signals.SignalNo
		name  string
	}{
		{signals.Aux, 1, "Aux 1"},
		{signals.Group, 2, "Group 2"},
		{signals.Direct, 3, "Invalid"},
	} {
		if got, want := signalControlName(tt.sig, tt.sigNo), tt.name; got != want {
			t.Errorf("signalControlName(%s, %d) = %s, want %s", tt.sig, tt.sigNo, got, want)
			continue
		}
	}
}

package venue

import (
	"testing"

	"github.com/kward/venue/router/signals"
)

// This test references Output fields to indicate intent and avoid U1000 until
// Output is fully wired in production code.
func TestOutputFieldsReferenced(t *testing.T) {
	_ = Output{sig: signals.Aux, sigNo: 1}
}

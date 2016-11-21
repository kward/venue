package workflow

import "testing"

func TestWorkflowAddr(t *testing.T) {
	for _, tt := range []struct {
		desc string
		wf   Workflow
		addr string
	}{
		{"gain", Workflow{
			version: "0.1",
			layout:  "th",
			page:    "soundcheck",
			control: "input",
			verb:    "gain"},
			"/0.1/th/soundcheck/input/gain"},
		{"select", Workflow{
			version: "0.1",
			layout:  "th",
			page:    "soundcheck",
			control: "input",
			verb:    "select",
			ref1:    1,
			ref2:    2},
			"/0.1/th/soundcheck/input/select/1/2"},
		{"select_label", Workflow{
			version: "0.1",
			layout:  "th",
			page:    "soundcheck",
			control: "input",
			verb:    "select",
			ref1:    1,
			label:   true},
			"/0.1/th/soundcheck/input/select/1/label"},
	} {
		if got, want := tt.wf.Addr(), tt.addr; got != want {
			t.Errorf("%s: Addr() = %s, want = %s", tt.desc, got, want)
		}
	}
}

package oscparse

import "testing"

const venueReq = "/venue/version/layout/page/control/command/x/y/label"

type lexTest struct {
	name  string
	input string
	items []item
}

var (
	tEOF           = item{itemEOF, "EOF"}
	tErrInvalidReq = item{itemError, "invalid request"}
	tLabel         = item{itemLabel, "label"}
	tPingReq       = item{itemPingReq, "ping"}
	tVenueReq      = item{itemVenueReq, "venue"}
)

var lexTests = []lexTest{
	{"empty", "", []item{tErrInvalidReq}},
	{"ping", "/ping", []item{tPingReq, tEOF}},
	{"generic", "/venue/version/layout/page/control/command/1/2/label", []item{
		tVenueReq,
		item{itemVersion, "version"},
		item{itemLayout, "layout"},
		item{itemPage, "page"},
		item{itemControl, "control"},
		item{itemCommand, "command"},
		item{itemPositionX, "1"},
		item{itemPositionY, "2"},
		tLabel,
		tEOF,
	}},
	{"cmd_w/o_label", "/venue/version/layout/page/control/command/1/2", []item{
		tVenueReq,
		item{itemVersion, "version"},
		item{itemLayout, "layout"},
		item{itemPage, "page"},
		item{itemControl, "control"},
		item{itemCommand, "command"},
		item{itemPositionX, "1"},
		item{itemPositionY, "2"},
		tEOF,
	}},
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.items) {
			t.Errorf("%s: got = %v, want = %v", test.name, items, test.items)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	for _, tt := range []struct {
		s  string
		is bool
	}{
		{"", false},
		{"abc", false},
		{"123", true},
		{"abc123", false},
		{"abc<123>", false},
	} {
		if got, want := isNumeric(tt.s), tt.is; got != want {
			t.Errorf("isAlphaNumeric(%s) = %v; want = %v", tt.s, got, want)
		}
	}
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []item) {
	l := lex(t.name, t.input)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func equal(i1, i2 []item) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
	}
	return true
}

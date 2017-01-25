/*
lex takes an OSC string, and converts it into tokens according to the position
of the elements.

The current lexable OSC string looks like:
/request/version/layout/page/control/command[/position][/position/][/label]

Lexing like this:
- The lexer l is instantiated.
- A go routine is started that repeatedly sends the OSC string out for lexing,
  using the previous token as a reference of how to lex the next token.
- To parse a string:
  - l.next() is called to provide the next token, minus the delimiter.
  - If the token can be parsed, the token is emitted.
  - If the token cannot be parsed, the EOF token is emitted.
- Once EOF is received by the go routine, the go routine ends.

Rather than parse the OSC string using a regular expression, lexing was chosen
to provide a proof-of-concept of how to do it as it is anticipated that future
versions of Venue will require modification, and lexing is easier to maintain
long-term.
*/
package oscparse

import (
	"fmt"
	"strings"
	"unicode"
)

// itemType identifies the type of lex items.
type itemType int

// item represents a text string returned from the scanner.
type item struct {
	typ itemType // The type of this item.
	val string   // The value of this item.
}

const (
	eof   = "EOF"
	label = "label"
)

func (i item) String() string {
	return fmt.Sprintf("<%d:%s>", i.typ, i.val)
}

const (
	itemError itemType = iota // error occurred; value is text of error
	itemCommand
	itemControl
	itemEOF
	itemLabel
	itemLayout
	itemPage
	itemPositionX
	itemPositionY
	itemPingReq
	itemVenueReq
	itemVersion
)

// stateFn represents the state of the scanner as a function that returns the
// next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name  string    // The name of the input; used only for error reporting.
	input string    // The string being scanned.
	state stateFn   // The next lexing function to enter.
	pos   int       // The current position.
	xy    int       // Bit field of Whether we have emitted none (0), x (1) or y (2).
	items chan item // Channel of scanned items.
}

const oscDelim = "/"

func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

// lexOSC parses an OSC path if possible.
func lexOSC(l *lexer) stateFn {
	if strings.HasPrefix(l.input[0:], oscDelim) {
		l.pos = 1
		return lexRequest
	}
	return l.errorf("invalid request")
}

func lexRequest(l *lexer) stateFn {
	if l.eof() {
		return l.errorf(eof)
	}
	switch s := l.next(); {
	case s == "ping":
		l.emit(itemPingReq, s)
		l.emit(itemEOF, eof)
		return nil
	case s == "venue":
		l.emit(itemVenueReq, s)
		return lexVersion
	}
	return l.errorf("unrecognized request")
}

func lexVersion(l *lexer) stateFn {
	if l.eof() {
		return l.errorf(eof)
	}
	l.emit(itemVersion, l.next())
	return lexLayout
}

func lexLayout(l *lexer) stateFn {
	if l.eof() {
		return l.errorf(eof)
	}
	l.emit(itemLayout, l.next())
	return lexPage
}

func lexPage(l *lexer) stateFn {
	if l.eof() {
		return l.errorf(eof)
	}
	l.emit(itemPage, l.next())
	return lexControl
}

func lexControl(l *lexer) stateFn {
	if l.eof() {
		return l.errorf(eof)
	}
	l.emit(itemControl, l.next())
	return lexCommand
}

func lexCommand(l *lexer) stateFn {
	if l.eof() {
		return l.errorf(eof)
	}
	l.emit(itemCommand, l.next())
	return lexPosition
}

func lexPosition(l *lexer) stateFn {
	switch s := l.next(); {
	case s == label:
		return lexLabel
	case isNumeric(s):
		switch l.xy {
		case 0:
			l.emit(itemPositionX, s)
			l.xy++
			return lexPosition
		case 1:
			l.emit(itemPositionY, s)
			l.xy++
			return lexLabel
		}
	case s == "":
		return lexEOF
	}
	return l.errorf("invalid position")
}

func lexLabel(l *lexer) stateFn {
	switch s := l.next(); {
	case s == label:
		l.emit(itemLabel, label)
		return lexEOF
	case s == "":
		return lexEOF
	}
	return l.errorf("invalid label")
}

func lexEOF(l *lexer) stateFn {
	l.emit(itemEOF, eof)
	return nil
}

func (l *lexer) emit(t itemType, s string) {
	if t == itemEOF {
		l.items <- item{t, s}
		return
	}
	l.items <- item{t, s}
}

func (l *lexer) eof() bool {
	if l.pos >= len(l.input) {
		return true
	}
	return false
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) next() string {
	if l.eof() {
		return ""
	}
	s := strings.SplitN(l.input[l.pos:], oscDelim, 2)
	if len(s) == 0 {
		return ""
	}
	l.pos += len(s[0]) + len(oscDelim)
	return s[0]
}

func (l *lexer) nextItem() item {
	item := <-l.items
	return item
}

func (l *lexer) run() {
	for l.state = lexOSC; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

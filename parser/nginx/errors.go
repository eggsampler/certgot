package nginx

import (
	"bytes"
	"fmt"
	"strings"
)

// ErrorLister is the public interface to access the inner errors
// included in a errList
type ErrorLister interface {
	Errors() []error
}

func (e errList) Errors() []error {
	return e
}

// ParserError is the public interface to errors of type parserError
type ParserError interface {
	Error() string
	InnerError() error
	Pos() (int, int, int)
	Expected() []string
}

func (p *parserError) InnerError() error {
	return p.Inner
}

func (p *parserError) Pos() (line, col, offset int) {
	return p.pos.line, p.pos.col, p.pos.offset
}

func (p *parserError) Expected() []string {
	return p.expected
}

func caretError(err error, input string) string {
	if el, ok := err.(ErrorLister); ok {
		var buffer bytes.Buffer
		for _, e := range el.Errors() {
			if parserErr, ok := e.(ParserError); ok {
				_, col, off := parserErr.Pos()
				line := extractLine(input, off)
				if col >= len(line) {
					col = len(line) - 1
				} else {
					if col > 0 {
						col--
					}
				}
				if col < 0 {
					col = 0
				}
				pos := col
				for _, chr := range line[:col] {
					if chr == '\t' {
						pos += 7
					}
				}
				buffer.WriteString(fmt.Sprintf("%s\n%s\n%s\n", line, strings.Repeat(" ", pos)+"^", err.Error()))
			} else {
				return err.Error()
			}
		}
		return buffer.String()
	}
	return err.Error()
}

func extractLine(input string, initPos int) string {
	if initPos < 0 {
		initPos = 0
	}
	if initPos >= len(input) && len(input) > 0 {
		initPos = len(input) - 1
	}
	startPos := initPos
	endPos := initPos
	for ; startPos > 0; startPos-- {
		if input[startPos] == '\n' {
			if startPos != initPos {
				startPos++
				break
			}
		}
	}
	for ; endPos < len(input); endPos++ {
		if input[endPos] == '\n' {
			if endPos == initPos {
				endPos++
			}
			break
		}
	}
	return input[startPos:endPos]
}

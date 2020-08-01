package parser

import (
	"testing"
)

func TestLetter(t *testing.T) {
	testList := []parserTest{
		{
			testName:  "letter",
			parser:    Letter(),
			input:     NewStringInput("a"),
			success:   true,
			resultStr: "a",
		},
		{
			testName:  "letter remainder",
			parser:    Letter(),
			input:     NewStringInput("ab"),
			success:   true,
			resultStr: "a",
			remainStr: "b",
		},
		{
			testName:  "number",
			parser:    Letter(),
			input:     NewStringInput("1"),
			success:   false,
			remainStr: "1",
		},
	}

	runTests(testList, t)
}

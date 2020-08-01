package parser

import (
	"fmt"
	"testing"
)

type parserTest struct {
	testName  string
	parser    Parser
	input     Input
	success   bool
	resultStr string
	remainStr string
}

func runTests(testList []parserTest, t *testing.T) {
	for _, currentTest := range testList {
		result := currentTest.parser(currentTest.input)
		if result.Success != currentTest.success {
			t.Fatalf("test %q: expected success %t, got: %t",
				currentTest.testName, currentTest.success, result.Success)
		}
		if result.String != currentTest.resultStr {
			t.Fatalf("test %q: expected result %q, got: %q",
				currentTest.testName, currentTest.resultStr, result.String)
		}
		if result.Remaining == nil {
			if currentTest.remainStr != "" {
				t.Fatalf("test %q: expected remaining %q, got nothing",
					currentTest.testName, currentTest.remainStr)
			}
		} else {
			remaining := fmt.Sprintf("%s", result.Remaining)
			if remaining != currentTest.remainStr {
				t.Fatalf("test %q: expected remaining %q, got: %q",
					currentTest.testName, currentTest.remainStr, remaining)
			}
		}
	}
}

func TestZeroOrMore(t *testing.T) {
	testList := []parserTest{
		{
			testName:  "zero+ letters w/ numbers only",
			parser:    ZeroOrMore(Letter()),
			input:     NewStringInput("123"),
			success:   true,
			resultStr: "",
			remainStr: "123",
		},
		{
			testName:  "zero+ letters w/ letters",
			parser:    ZeroOrMore(Letter()),
			input:     NewStringInput("abc"),
			success:   true,
			resultStr: "abc",
			remainStr: "",
		},
		{
			testName:  "zero+ letters w/ both",
			parser:    ZeroOrMore(Letter()),
			input:     NewStringInput("abc123"),
			success:   true,
			resultStr: "abc",
			remainStr: "123",
		},
	}

	runTests(testList, t)
}

func TestOneOrMore(t *testing.T) {
	testList := []parserTest{
		{
			testName:  "1+ letters w/ numbers only",
			parser:    OneOrMore(Letter()),
			input:     NewStringInput("123"),
			success:   false,
			resultStr: "",
			remainStr: "123",
		},
		{
			testName:  "1+ letters w/ letters",
			parser:    OneOrMore(Letter()),
			input:     NewStringInput("abc"),
			success:   true,
			resultStr: "abc",
			remainStr: "",
		},
		{
			testName:  "1+ letters w/ both",
			parser:    OneOrMore(Letter()),
			input:     NewStringInput("abc123"),
			success:   true,
			resultStr: "abc",
			remainStr: "123",
		},
	}

	runTests(testList, t)
}

func TestOr(t *testing.T) {
	testList := []parserTest{
		{
			testName:  "or letter",
			parser:    Or(Letter()),
			input:     NewStringInput("123"),
			success:   false,
			resultStr: "",
			remainStr: "123",
		},
		{
			testName:  "or number",
			parser:    Or(Number()),
			input:     NewStringInput("123"),
			success:   true,
			resultStr: "1",
			remainStr: "23",
		},
		{
			testName:  "or both",
			parser:    Or(Number(), Letter()),
			input:     NewStringInput("123"),
			success:   true,
			resultStr: "1",
			remainStr: "23",
		},
		{
			testName:  "or both",
			parser:    Or(Number(), Letter()),
			input:     NewStringInput("abc"),
			success:   true,
			resultStr: "a",
			remainStr: "bc",
		},
	}

	runTests(testList, t)
}

func TestAnd(t *testing.T) {
	testList := []parserTest{
		{
			testName:  "and letter",
			parser:    And(Letter()),
			input:     NewStringInput("123"),
			success:   false,
			resultStr: "",
			remainStr: "123",
		},
		{
			testName:  "and number",
			parser:    And(Number()),
			input:     NewStringInput("123"),
			success:   true,
			resultStr: "1",
			remainStr: "23",
		},
		{
			testName:  "and number letter fail",
			parser:    And(Number(), Letter()),
			input:     NewStringInput("123"),
			success:   false,
			resultStr: "",
			remainStr: "123",
		},
		{
			testName:  "and both letters",
			parser:    And(Letter(), Letter()),
			input:     NewStringInput("abc"),
			success:   true,
			resultStr: "ab",
			remainStr: "c",
		},
	}

	runTests(testList, t)
}

func TestOptional(t *testing.T) {
	testList := []parserTest{
		{
			testName:  "optional 1",
			parser:    Optional(Letter()),
			input:     NewStringInput("1"),
			success:   true,
			resultStr: "",
			remainStr: "1",
		},
		{
			testName:  "optional 2",
			parser:    Optional(Letter()),
			input:     NewStringInput("a"),
			success:   true,
			resultStr: "a",
			remainStr: "",
		},
	}

	runTests(testList, t)
}

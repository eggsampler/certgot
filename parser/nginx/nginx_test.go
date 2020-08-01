package nginx

import (
	"testing"

	"github.com/eggsampler/certgot/parser"
)

func Test_simpleDirectiveParser(t *testing.T) {
	testList := []struct {
		testName  string
		input     string
		success   bool
		resultStr string
		remainStr string
	}{
		{
			testName:  "succeed",
			input:     "testy mctest;",
			success:   true,
			resultStr: "testy mctest;",
		},
	}

	for _, currentTest := range testList {
		result := simpleDirectiveParser(parser.NewStringInput(currentTest.input))
		if result.Success != currentTest.success {
			t.Fatalf("test %q: expected success %t, got: %t (%+v)",
				currentTest.testName, currentTest.success, result.Success, result)
		}
		if result.String != currentTest.resultStr {
			t.Fatalf("test %q: expected result %q, got: %q",
				currentTest.testName, currentTest.resultStr, result.String)
		}
	}
}

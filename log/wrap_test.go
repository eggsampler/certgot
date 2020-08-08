package log

import "testing"

func TestPrintWrap(t *testing.T) {
	lorem := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et " +
		"dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex " +
		"ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu " +
		"fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt " +
		"mollit anim id est laborum."

	testList := []struct {
		testName string
		input    string
		len      int
		prefix   string
		output   string
	}{
		{
			testName: "nothing",
		},
		{
			testName: "lorem ipsum",
			input:    lorem,
			len:      80,
			output: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor\n" +
				"incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis\n" +
				"nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.\n" +
				"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu\n" +
				"fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in\n" +
				"culpa qui officia deserunt mollit anim id est laborum.",
		},
		{
			testName: "lorem ipsum prefix",
			input:    lorem,
			len:      80,
			prefix:   "  ",
			output: "  Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod\n" +
				"  tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,\n" +
				"  quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo\n" +
				"  consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse\n" +
				"  cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non\n" +
				"  proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		},
	}

	for _, currentTest := range testList {
		t.Run(currentTest.testName, func(t *testing.T) {
			out := Wrap(currentTest.input, currentTest.len, currentTest.prefix)
			if out != currentTest.output {
				t.Fatalf("output mismatch\n expected(%d): %q\n got(%d): %q",
					len(currentTest.output), currentTest.output, len(out), out)
			}
		})
	}
}

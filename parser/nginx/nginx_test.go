package nginx

import (
	"fmt"
	"reflect"
	"testing"
)

var exampleConfig = `
server {
 listen       80;
 server_name  domain2.com www.domain2.com;
 access_log   logs/domain2.access.log  main;
 
 location ~ ^/(images|javascript|js|css|flash|media|static)/
 {
  root    /var/www/virtual/big.server.com/htdocs;
  expires 30d;
 }

 location / {
  proxy_pass      http://127.0.0.1:8080;
 }
}
`

func TestParse(t *testing.T) {
	testList := []struct {
		testName   string
		input      string
		fileName   string
		hasError   bool
		expectsDir bool
		checkFunc  func(Directive) error
	}{
		{
			testName:   "no input",
			expectsDir: true, // should still return a main directive, i guess
		},
		{
			testName: "invalid directive",
			input:    "hello",
			hasError: true,
		},
		{
			testName:   "simplest directive",
			input:      "hello;",
			expectsDir: true,
			checkFunc: func(d Directive) error {
				od := Directive{
					Name:          "main",
					HasDirectives: true,
					Directives: []Directive{
						{Name: "hello"},
					},
				}
				if !reflect.DeepEqual(d, od) {
					return fmt.Errorf("directive mismatch\n expects: %+v\n got: %+v", od, d)
				}
				return nil
			},
		},
		{
			testName:   "simple directive",
			input:      "hello_world;",
			expectsDir: true,
		},
		{
			testName:   "simplest directive spces",
			input:      " hello_world ; ",
			expectsDir: true,
		},
		{
			testName:   "simple directive with parameter",
			input:      "hello world;",
			expectsDir: true,
		},
		{
			testName:   "simple directive with parameter and spaces",
			input:      " hello_world value ; ",
			expectsDir: true,
		},
		{
			testName:   "multiple simple directives",
			input:      "hello world;foo bar;",
			expectsDir: true,
		},
		{
			testName:   "empty block directive",
			input:      "hello {}",
			expectsDir: true,
		},
		{
			testName:   "block directive spaced",
			input:      " hello { } ",
			expectsDir: true,
		},
		{
			testName:   "block directive with simple",
			input:      "hello world { foo bar; }",
			expectsDir: true,
		},
		{
			testName:   "block directive with simple",
			input:      "hello world {\n\tfoo bar;\n}",
			expectsDir: true,
		},
		{
			testName:   "block directives nested",
			input:      "hello world { location { foo bar; } }",
			expectsDir: true,
		},
		{
			testName:   "real config i guess",
			input:      exampleConfig,
			expectsDir: true,
		},
	}

	for _, currentTest := range testList {
		output, err := Parse(currentTest.fileName, []byte(currentTest.input))
		if currentTest.hasError == (err == nil) {
			if err != nil {
				fmt.Println(caretError(err, currentTest.input))
			}
			t.Fatalf("%q: expected error %v, got: %v", currentTest.testName, currentTest.hasError, err)
		}
		directive, ok := output.(Directive)
		if currentTest.expectsDir != ok {
			t.Fatalf("test %q: expects directive %t, got: %t", currentTest.testName, currentTest.expectsDir, ok)
		}
		if currentTest.checkFunc != nil {
			if err := currentTest.checkFunc(directive); err != nil {
				t.Fatalf("test %q: check: %v", currentTest.testName, err)
			}
		}
	}
}

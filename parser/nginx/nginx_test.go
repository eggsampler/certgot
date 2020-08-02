package nginx

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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
		equalCheck []interface{}
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
			equalCheck: []interface{}{
				SimpleDirective{Name: "hello"},
			},
		},
		{
			testName:   "simple directive",
			input:      "hello_world;",
			expectsDir: true,
		},
		{
			testName:   "simplest directive spaces",
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
		{
			testName:   "comment",
			input:      "# hello world",
			expectsDir: true,
			equalCheck: []interface{}{
				CommentDirective("hello world"),
			},
		},
		{
			testName:   "comment spaced",
			input:      " # hello world ",
			expectsDir: true,
		},
		{
			testName:   "no comment",
			input:      "#",
			expectsDir: true,
		},
		{
			testName:   "no comment spaces 1",
			input:      " #",
			expectsDir: true,
		},
		{
			testName:   "no comment spaces 2",
			input:      " # ",
			expectsDir: true,
		},
		{
			testName:   "no comment next line 1",
			input:      " # \r\n hello world;",
			expectsDir: true,
		},
		{
			testName:   "no comment next line 2",
			input:      " # \n hello world;",
			expectsDir: true,
		},
		{
			testName:   "single quote semi colon",
			input:      "hello 'foo;bar' world;",
			expectsDir: true,
			equalCheck: []interface{}{
				SimpleDirective{Name: "hello", Parameter: `'foo;bar' world`},
			},
		},
		{
			testName:   "double quote semi colon",
			input:      `hello "foo;bar" world;`,
			expectsDir: true,
			equalCheck: []interface{}{
				SimpleDirective{Name: "hello", Parameter: `"foo;bar" world`},
			},
		},
		{
			testName:   "single quote brace",
			input:      "hello 'foo{bar' world;",
			expectsDir: true,
			equalCheck: []interface{}{
				SimpleDirective{Name: "hello", Parameter: `'foo{bar' world`},
			},
		},
		{
			testName:   "double quote brace",
			input:      `hello "foo;bar" world;`,
			expectsDir: true,
			equalCheck: []interface{}{
				SimpleDirective{Name: "hello", Parameter: `"foo{bar" world`},
			},
		},
	}

	for _, currentTest := range testList {
		output, err := Parse(currentTest.fileName, []byte(currentTest.input))
		if currentTest.hasError == (err == nil) {
			if err != nil {
				fmt.Println(caretError(err, currentTest.input))
			}
			t.Fatalf("test %q: expected error %v, got: %v", currentTest.testName, currentTest.hasError, err)
		}
		directives, ok := output.([]interface{})
		if currentTest.expectsDir != ok {
			t.Fatalf("test %q: expects directive %t, got: %t", currentTest.testName, currentTest.expectsDir, ok)
		}
		if currentTest.equalCheck != nil {
			if !reflect.DeepEqual(directives, currentTest.equalCheck) {
				t.Fatalf("test %q: directive mismatch\n expects: %+v\n got: %+v",
					currentTest.testName, currentTest.equalCheck, directives)
			}
		}
	}
}

func TestParseFile(t *testing.T) {
	fileList := []struct {
		fileName string
		hasError bool
	}{
		{
			fileName: filepath.Join("testdata", "nginx.conf"),
		},
		{
			fileName: filepath.Join("testdata", "broken.conf"),
			hasError: true,
		},
		{
			fileName: filepath.Join("testdata", "comment_in_file.conf"),
		},
		{
			fileName: filepath.Join("testdata", "edge_cases.conf"),
		},
	}

	for _, currentTest := range fileList {
		out, err := ParseFile(currentTest.fileName)
		if currentTest.hasError == (err == nil) {
			if err != nil {
				input, err2 := ioutil.ReadFile(currentTest.fileName)
				if err2 != nil {
					panic(err)
				}
				fmt.Println(caretError(err, string(input)))
			}
			t.Fatalf("test %q: expected error %v, got: %v\n output: %+v",
				currentTest.fileName, currentTest.hasError, err, out)
		}
	}
}

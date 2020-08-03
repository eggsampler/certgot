package nginx

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
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
		errorStr   string
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
			testName:   "simple directive with parameters",
			input:      "hello 'foo bar' world;",
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
				SimpleDirective{Name: "hello", Parameters: []string{"'foo;bar'", "world"}},
			},
		},
		{
			testName:   "double quote semi colon",
			input:      `hello "foo;bar" world;`,
			expectsDir: true,
			equalCheck: []interface{}{
				SimpleDirective{Name: "hello", Parameters: []string{`"foo;bar"`, "world"}},
			},
		},
		{
			testName:   "single quote brace",
			input:      "hello 'foo{bar' world;",
			expectsDir: true,
			equalCheck: []interface{}{
				SimpleDirective{Name: "hello", Parameters: []string{"'foo{bar'", "world"}},
			},
		},
		{
			testName:   "double quote brace",
			input:      `hello "foo{bar" world;`,
			expectsDir: true,
			equalCheck: []interface{}{
				SimpleDirective{Name: "hello", Parameters: []string{`"foo{bar"`, "world"}},
			},
		},
		{
			testName:   "edge case",
			input:      `add_header  Cache-Control  'public, must-revalidate, proxy-revalidate' "test,;{}" foo;`,
			expectsDir: true,
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
		if err != nil && !strings.Contains(err.Error(), currentTest.errorStr) {
			t.Fatalf("test %q: expected %q in error: %v", currentTest.testName, currentTest.errorStr, err)
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
		fileName   string
		hasError   bool
		errorStr   string
		expectsDir bool
		equalCheck []interface{}
	}{
		{
			fileName: filepath.Join("testdata", "broken.conf"),
			hasError: true,
		},
		{
			fileName:   filepath.Join("testdata", "comment_in_file.conf"),
			expectsDir: true,
			equalCheck: []interface{}{
				CommentDirective("a comment inside a file"),
			},
		},
		{
			fileName:   filepath.Join("testdata", "edge_cases.conf"),
			expectsDir: true,
			equalCheck: []interface{}{
				CommentDirective("This is not a valid nginx config file but it tests edge cases in valid nginx syntax"),
				BlockDirective{
					Name: "server",
					Children: []interface{}{
						SimpleDirective{
							Name:       "server_name",
							Parameters: []string{"simple"},
						},
					},
				},
				BlockDirective{
					Name: "server",
					Children: []interface{}{
						SimpleDirective{
							Name:       "server_name",
							Parameters: []string{"with.if"},
						},
						BlockDirective{
							Name:       "location",
							Parameters: []string{"~", "^/services/.+$"},
							Children: []interface{}{
								BlockDirective{
									Name:       "if",
									Parameters: []string{"($request_filename", "~*", `\.(ttf|woff)$)`},
									Children: []interface{}{
										SimpleDirective{
											Name:       "add_header",
											Parameters: []string{"Access-Control-Allow-Origin", `"*"`},
										},
									},
								},
							},
						},
					},
				},
				BlockDirective{
					Name: "server",
					Children: []interface{}{
						SimpleDirective{
							Name:       "server_name",
							Parameters: []string{"with.complicated.headers"},
						},
						BlockDirective{
							Name:       "location",
							Parameters: []string{"~*", "\\.(?:gif|jpe?g|png)$"},
							Children: []interface{}{
								SimpleDirective{
									Name:       "add_header",
									Parameters: []string{"Pragma", "public"},
								},
								SimpleDirective{
									Name:       "add_header",
									Parameters: []string{"Cache-Control", `'public, must-revalidate, proxy-revalidate'`, `"test,;{}"`, "foo"},
								},
								SimpleDirective{
									Name:       "blah",
									Parameters: []string{`"hello;world"`},
								},
								SimpleDirective{
									Name:       "try_files",
									Parameters: []string{"$uri", "@rewrites"},
								},
							},
						},
					},
				},
			},
		},
		{
			fileName:   filepath.Join("testdata", "foo.conf"),
			expectsDir: true,
		},
		{
			fileName:   filepath.Join("testdata", "invalid_unicode_comments.conf"),
			hasError:   true,
			errorStr:   "invalid encoding",
			expectsDir: true, // TODO: should this actually "succeed" to parse?
		},
		{
			fileName:   filepath.Join("testdata", "minimalistic_comments.conf"),
			expectsDir: true,
		},
		{
			fileName:   filepath.Join("testdata", "multiline_quotes.conf"),
			expectsDir: true,
			equalCheck: []interface{}{
				CommentDirective("Test nginx configuration file with multiline quoted strings."),
				CommentDirective("Good example of usage for multilined quoted values is when"),
				CommentDirective("using Openresty's Lua directives and you wish to keep the"),
				CommentDirective("inline Lua code readable."),
				BlockDirective{
					Name: "http",
					Children: []interface{}{
						BlockDirective{
							Name: "server",
							Children: []interface{}{
								SimpleDirective{
									Name:       "listen",
									Parameters: []string{"*:443"},
								},
								CommentDirective("because there should be no other port open."),
								BlockDirective{
									Name:       "location",
									Parameters: []string{"/"},
									Children: []interface{}{
										SimpleDirective{
											Name: "body_filter_by_lua",
											Parameters: []string{`'ngx.ctx.buffered = (ngx.ctx.buffered or "") .. string.sub(ngx.arg[1], 1, 1000)
                            if ngx.arg[2] then
                              ngx.var.resp_body = ngx.ctx.buffered
                            end'`},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			fileName:   filepath.Join("testdata", "nginx.conf"),
			expectsDir: true,
		},
		{
			fileName:   filepath.Join("testdata", "server.conf"),
			expectsDir: true,
		},
		{
			fileName:   filepath.Join("testdata", "valid_unicode_comments.conf"),
			expectsDir: true,
		},
	}

	for _, currentTest := range fileList {
		output, err := ParseFile(currentTest.fileName)
		if currentTest.hasError == (err == nil) {
			if err != nil {
				input, err2 := ioutil.ReadFile(currentTest.fileName)
				if err2 != nil {
					panic(err)
				}
				fmt.Println(caretError(err, string(input)))
			}
			t.Fatalf("test %q: expected error %v, got: %v\n output: %+v",
				currentTest.fileName, currentTest.hasError, err, output)
		}
		if err != nil && !strings.Contains(err.Error(), currentTest.errorStr) {
			t.Fatalf("test %q: expected %q in error: %v", currentTest.fileName, currentTest.errorStr, err)
		}
		directives, ok := output.([]interface{})
		if currentTest.expectsDir != ok {
			t.Fatalf("test %q: expects directive %t, got: %t\n output: %#v",
				currentTest.fileName, currentTest.expectsDir, ok, directives)
		}
		if currentTest.equalCheck != nil {
			if !reflect.DeepEqual(directives, currentTest.equalCheck) {
				t.Fatalf("test %q: directive mismatch\n expects: %#v\n got: %#v",
					currentTest.fileName, currentTest.equalCheck, directives)
			}
		}
	}
}

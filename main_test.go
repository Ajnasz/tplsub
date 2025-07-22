package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		data        any
		expected    string
		expectError bool
	}{
		// Basic template tests
		{
			name:     "simple text template",
			template: "hello world",
			data:     map[string]any{},
			expected: "hello world",
		},
		{
			name:     "basic variable substitution",
			template: "Hello {{ .FirstName }} {{ .LastName }}",
			data:     map[string]any{"FirstName": "John", "LastName": "Doe"},
			expected: "Hello John Doe",
		},
		{
			name:     "template with data file format",
			template: "Data file: {{ .FirstName }} {{ .LastName }}",
			data:     map[string]any{"FirstName": "John", "LastName": "Doe"},
			expected: "Data file: John Doe",
		},

		// String helper tests
		{
			name:     "string lower function",
			template: "Param tpl: {{ .FirstName }} {{ .LastName | lower }}",
			data:     map[string]any{"FirstName": "John", "LastName": "Doe"},
			expected: "Param tpl: John doe",
		},
		{
			name:     "string upper function",
			template: "Upper: {{ .text | upper }}",
			data:     map[string]any{"text": "Hello World"},
			expected: "Upper: HELLO WORLD",
		},
		{
			name:     "string trim function",
			template: "Trim: \"{{ .text | trim }}\"",
			data:     map[string]any{"text": "  hello  "},
			expected: "Trim: \"hello\"",
		},
		{
			name:     "string replace function",
			template: "Replace: \"{{ .text | replace \"World\" \"Gopher\" }}\"",
			data:     map[string]any{"text": "Hello World"},
			expected: "Replace: \"Hello Gopher\"",
		},
		{
			name:     "string split function",
			template: "Split: {{ split \",\" .text }}",
			data:     map[string]any{"text": "hello,world,test"},
			expected: "Split: [hello world test]",
		},
		{
			name:     "string join function",
			template: "Join: {{ .items | toStrings | join \"-\" }}",
			data:     map[string]any{"items": []any{"a", "b", "c"}},
			expected: "Join: a-b-c",
		},
		{
			name:     "string contains function",
			template: "Contains \"world\": {{ contains \"world\" .text }}",
			data:     map[string]any{"text": "hello world"},
			expected: "Contains \"world\": true",
		},
		{
			name:     "string repeat function",
			template: "Repeat: {{ .FirstName | repeat 3 }} {{ .LastName | repeat 2 }}",
			data:     map[string]any{"FirstName": "John", "LastName": "Doe"},
			expected: "Repeat: JohnJohnJohn DoeDoe",
		},

		// Math helper tests
		{
			name:     "math add function",
			template: "Add: {{ add .a .b }}",
			data:     map[string]any{"a": 10, "b": 3},
			expected: "Add: 13",
		},
		{
			name:     "math sub function",
			template: "Sub: {{ sub .a .b }}",
			data:     map[string]any{"a": 10, "b": 3},
			expected: "Sub: 7",
		},
		{
			name:     "math mul function",
			template: "Mul: {{ mul .a .b }}",
			data:     map[string]any{"a": 10, "b": 3},
			expected: "Mul: 30",
		},
		{
			name:     "math div function",
			template: "Div: {{ div .a .b }}",
			data:     map[string]any{"a": 10, "b": 3},
			expected: "Div: 3",
		},
		{
			name:     "math mod function",
			template: "Mod: {{ mod .a .b }}",
			data:     map[string]any{"a": 10, "b": 3},
			expected: "Mod: 1",
		},

		// Float math helper tests
		{
			name:     "float add function",
			template: "Add float: {{ addf .a .b }}",
			data:     map[string]any{"a": 10.5, "b": 3.2},
			expected: "Add float: 13.7",
		},
		{
			name:     "float sub function",
			template: "Sub float: {{ subf .a .b }}",
			data:     map[string]any{"a": 10.5, "b": 3.2},
			expected: "Sub float: 7.3",
		},
		{
			name:     "float mul function",
			template: "Mul float: {{ mulf .a .b }}",
			data:     map[string]any{"a": 10.5, "b": 3.2},
			expected: "Mul float: 33.6",
		},
		{
			name:     "float div function",
			template: "Div float: {{ divf .a .b }}",
			data:     map[string]any{"a": 10.5, "b": 3.2},
			expected: "Div float: 3.28125",
		},
		{
			name:     "mixed types float conversion",
			template: "Mixed types: {{ addf .a .b }} (string + int)",
			data:     map[string]any{"a": "15.75", "b": 4},
			expected: "Mixed types: 19.75 (string + int)",
		},
		{
			name:     "to float conversion",
			template: "To float: {{ toFloat .value }} (converted to float)",
			data:     map[string]any{"value": 42},
			expected: "To float: 42 (converted to float)",
		},
		{
			name:     "precise division comparison",
			template: "Precise division: {{ divf .a .b }} vs {{ div .a .b }}",
			data:     map[string]any{"a": 22, "b": 7},
			expected: "Precise division: 3.142857142857143 vs 3",
		},

		// Date helper tests (these will vary by execution time, so we'll test format only)
		{
			name:     "date year extraction",
			template: "min: {{ .time | parseDate \"2006-01-02 15:04:05 MST\" | year }}",
			data:     map[string]any{"time": "2025-07-20 17:17:00 CEST"},
			expected: "min: 2025",
		},

		// Collection helper tests
		{
			name:     "collection length",
			template: "Length: {{ len .items }}",
			data:     map[string]any{"items": []any{"first", "second", "third"}},
			expected: "Length: 3",
		},
		{
			name:     "collection first",
			template: "First: {{ first .items }}",
			data:     map[string]any{"items": []any{"first", "second", "third"}},
			expected: "First: first",
		},
		{
			name:     "collection last",
			template: "Last: {{ last .items }}",
			data:     map[string]any{"items": []any{"first", "second", "third"}},
			expected: "Last: third",
		},

		// Conditional helper tests
		{
			name:     "default with empty value",
			template: "Default: {{ default \"Anonymous\" .name }}",
			data:     map[string]any{"name": ""},
			expected: "Default: Anonymous",
		},
		{
			name:     "default with non-empty value",
			template: "Default: {{ default \"Anonymous\" .name }}",
			data:     map[string]any{"name": "John"},
			expected: "Default: John",
		},
		{
			name:     "empty check with empty value",
			template: "Empty check: {{ if empty .name }}Name is empty{{ else }}Name: {{ .name }}{{ end }}",
			data:     map[string]any{"name": ""},
			expected: "Empty check: Name is empty",
		},

		// File helper tests
		{
			name:     "file basename",
			template: "Basename: {{ basename .path }}",
			data:     map[string]any{"path": "/home/user/document.txt"},
			expected: "Basename: document.txt",
		},
		{
			name:     "file dirname",
			template: "Dirname: {{ dirname .path }}",
			data:     map[string]any{"path": "/home/user/document.txt"},
			expected: "Dirname: /home/user",
		},
		{
			name:     "file extension",
			template: "Extension: {{ ext .path }}",
			data:     map[string]any{"path": "/home/user/document.txt"},
			expected: "Extension: .txt",
		},
		{
			name:     "path join",
			template: "Join: {{ pathjoin .dir .subdir .file }}",
			data:     map[string]any{"dir": "/home", "subdir": "user", "file": "doc.txt"},
			expected: "Join: /home/user/doc.txt",
		},

		// Hashing helper tests
		{
			name:     "md5 hash",
			template: "MD5: {{ md5 .FirstName }} {{ md5 .LastName }}",
			data:     map[string]any{"FirstName": "John", "LastName": "Doe"},
			expected: "MD5: 61409aa1fd47d4a5332de23cbf59a36f ad695f53ae7569fb981fc95598e27e67",
		},
		{
			name:     "sha1 hash",
			template: "SHA1: {{ sha1 .text }}",
			data:     map[string]any{"text": "hello world"},
			expected: "SHA1: 2aae6c35c94fcfb415dbe95f408b9ce91ee846ed",
		},
		{
			name:     "sha256 hash",
			template: "SHA256: {{ sha256 .text }}",
			data:     map[string]any{"text": "hello world"},
			expected: "SHA256: b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			name:     "base64 encode",
			template: "Base64: {{ base64Encode .text }}",
			data:     map[string]any{"text": "hello world"},
			expected: "Base64: aGVsbG8gd29ybGQ=",
		},
		{
			name:     "base64 decode",
			template: "Decoded: {{ base64Decode .encoded }}",
			data:     map[string]any{"encoded": "aGVsbG8gd29ybGQ="},
			expected: "Decoded: hello world",
		},

		// JSON helper tests
		{
			name:     "toPrettyJSON",
			template: "{{ . | toPrettyJSON }}",
			data:     map[string]any{"FirstName": "John", "LastName": "Doe"},
			expected: "{\n  \"FirstName\": \"John\",\n  \"LastName\": \"Doe\"\n}",
		},

		// Error cases
		{
			name:        "invalid template syntax",
			template:    "Hello {{ .Name",
			data:        map[string]any{},
			expectError: true,
		},
		{
			name:        "invalid function call",
			template:    "{{ invalidFunc .text }}",
			data:        map[string]any{"text": "hello"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := executeTemplate(&buf, tt.template, tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("%s: expected error but got none %s", tt.name, buf.String())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			result := buf.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestExecuteTemplateWithNilData(t *testing.T) {
	var buf bytes.Buffer
	err := executeTemplate(&buf, "Hello World", nil)
	if err != nil {
		t.Errorf("unexpected error with nil data: %v", err)
	}

	result := buf.String()
	expected := "Hello World"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestExecuteTemplateWithEmptyTemplate(t *testing.T) {
	var buf bytes.Buffer
	err := executeTemplate(&buf, "", map[string]any{})
	if err != nil {
		t.Errorf("unexpected error with empty template: %v", err)
	}

	result := buf.String()
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestExecuteTemplateSequenceFunction(t *testing.T) {
	var buf bytes.Buffer
	template := "Sequence: {{ range seq 1 5 }}{{ . }} {{ end }}"
	data := map[string]any{}

	err := executeTemplate(&buf, template, data)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	result := strings.TrimSpace(buf.String())
	expected := "Sequence: 1 2 3 4 5"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

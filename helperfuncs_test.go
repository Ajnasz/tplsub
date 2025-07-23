package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"text/template"
	"time"
)

// Test string manipulation functions
func TestStringFunctions(t *testing.T) {
	funcs := createHelperFuncs()

	tests := []struct {
		name     string
		funcName string
		args     []any
		expected any
	}{
		{"upper", "upper", []any{"hello"}, "HELLO"},
		{"lower", "lower", []any{"HELLO"}, "hello"},
		{"trim", "trim", []any{"  hello  "}, "hello"},
		{"replace", "replace", []any{"o", "0", "hello"}, "hell0"},
		{"split", "split", []any{",", "a,b,c"}, []string{"a", "b", "c"}},
		{"contains true", "contains", []any{"ell", "hello"}, true},
		{"contains false", "contains", []any{"xyz", "hello"}, false},
		{"hasPrefix true", "hasPrefix", []any{"he", "hello"}, true},
		{"hasPrefix false", "hasPrefix", []any{"lo", "hello"}, false},
		{"hasSuffix true", "hasSuffix", []any{"lo", "hello"}, true},
		{"hasSuffix false", "hasSuffix", []any{"he", "hello"}, false},
		{"repeat", "repeat", []any{3, "ab"}, "ababab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := funcs[tt.funcName]
			if fn == nil {
				t.Fatalf("function %s not found", tt.funcName)
			}

			result := callFunc(fn, tt.args...)
			if !equalValues(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestJoinFunction(t *testing.T) {
	funcs := createHelperFuncs()
	joinFunc := funcs["join"]

	tests := []struct {
		name     string
		sep      string
		elems    []string
		expected string
	}{
		{"strings", ",", []string{"a", "b", "c"}, "a,b,c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := callFunc(joinFunc, tt.sep, tt.elems)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Test type conversion functions
func TestTypeConversionFunctions(t *testing.T) {
	funcs := createHelperFuncs()

	t.Run("toFloat", func(t *testing.T) {
		toFloatFunc := funcs["toFloat"]

		tests := []struct {
			input    any
			expected float64
			hasError bool
		}{
			{42, 42.0, false},
			{42.5, 42.5, false},
			{"42.5", 42.5, false},
			{"invalid", 0, true},
		}

		for _, tt := range tests {
			result, err := callFuncWithError(toFloatFunc, tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("expected error for input %v", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %v: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		}
	})

	t.Run("toInt", func(t *testing.T) {
		toIntFunc := funcs["toInt"]

		tests := []struct {
			input    any
			expected int
			hasError bool
		}{
			{42, 42, false},
			{42.0, 42, false},
			{"42", 42, false},
			{"invalid", 0, true},
		}

		for _, tt := range tests {
			result, err := callFuncWithError(toIntFunc, tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("expected error for input %v", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %v: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		}
	})

	t.Run("toString", func(t *testing.T) {
		toStringFunc := funcs["toString"]

		tests := []struct {
			input    any
			expected string
		}{
			{42, "42"},
			{42.5, "42.5"},
			{"hello", "hello"},
			{true, "true"},
		}

		for _, tt := range tests {
			result := callFunc(toStringFunc, tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		}
	})

	t.Run("toStrings", func(t *testing.T) {
		toStringsFunc := funcs["toStrings"]

		tests := []struct {
			input    []any
			expected []string
			hasError bool
		}{
			{[]any{1, 2.5, "hello"}, []string{"1", "2.5", "hello"}, false},
			{[]any{"invalid", 42}, []string{"invalid", "42"}, false},
			{[]any{true, nil}, []string{"true", ""}, false},
			{[]any{}, []string{}, false},
		}

		for _, tt := range tests {
			result, err := callFuncWithError(toStringsFunc, tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("expected error for input %v", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %v: %v", tt.input, err)
				}
				if !equalValues(result, tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		}
	})
}

// Test math operations
func TestMathFunctions(t *testing.T) {
	funcs := createHelperFuncs()

	t.Run("integer math", func(t *testing.T) {
		tests := []struct {
			name     string
			funcName string
			a, b     any
			expected int
			hasError bool
		}{
			{"add", "add", 5, 3, 8, false},
			{"sub", "sub", 3, 5, 2, false},
			{"mul", "mul", 5, 3, 15, false},
			{"div", "div", 3, 15, 5, false},
			{"mod", "mod", 3, 7, 1, false},
			{"add invalid", "add", "invalid", 3, 0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				fn := funcs[tt.funcName]
				result, err := callFuncWithError(fn, tt.a, tt.b)

				if tt.hasError {
					if err == nil {
						t.Errorf("expected error")
					}
				} else {
					if err != nil {
						t.Errorf("unexpected error: %v", err)
					}
					if result != tt.expected {
						t.Errorf("expected %v, got %v", tt.expected, result)
					}
				}
			})
		}
	})

	t.Run("float math", func(t *testing.T) {
		tests := []struct {
			name     string
			funcName string
			a, b     any
			expected float64
			hasError bool
		}{
			{"addf", "addf", 5.5, 3.2, 8.7, false},
			{"subf", "subf", 3.2, 5.5, 2.3, false},
			{"mulf", "mulf", 5.0, 3.0, 15.0, false},
			{"divf", "divf", 3.0, 15.0, 5.0, false},
			{"divf by zero", "divf", 0.0, 15.0, 0.0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				fn := funcs[tt.funcName]
				result, err := callFuncWithError(fn, tt.a, tt.b)

				if tt.hasError {
					if err == nil {
						t.Errorf("expected error")
					}
				} else {
					if err != nil {
						t.Errorf("unexpected error: %v", err)
					}
					if fmt.Sprintf("%.1f", result) != fmt.Sprintf("%.1f", tt.expected) {
						t.Errorf("expected %v, got %v", tt.expected, result)
					}
				}
			})
		}
	})
}

// Test date/time functions
func TestDateTimeFunctions(t *testing.T) {
	funcs := createHelperFuncs()

	t.Run("now", func(t *testing.T) {
		nowFunc := funcs["now"]
		result := callFunc(nowFunc)
		if _, ok := result.(time.Time); !ok {
			t.Errorf("expected time.Time, got %T", result)
		}
	})

	t.Run("parseDate", func(t *testing.T) {
		parseDateFunc := funcs["parseDate"]

		result, err := callFuncWithError(parseDateFunc, "2006-01-02", "2023-12-25")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		parsed := result.(time.Time)
		if parsed.Year() != 2023 || parsed.Month() != 12 || parsed.Day() != 25 {
			t.Errorf("date not parsed correctly: %v", parsed)
		}
	})

	t.Run("formatDate", func(t *testing.T) {
		formatDateFunc := funcs["formatDate"]
		testTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)

		result := callFunc(formatDateFunc, "2006-01-02", testTime)
		if result != "2023-12-25" {
			t.Errorf("expected 2023-12-25, got %s", result)
		}
	})

	t.Run("timestamp", func(t *testing.T) {
		timestampFunc := funcs["timestamp"]
		testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

		result := callFunc(timestampFunc, testTime)
		if result != int64(1672531200) {
			t.Errorf("expected 1672531200, got %v", result)
		}
	})

	t.Run("year/month/day", func(t *testing.T) {
		testTime := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)

		yearFunc := funcs["year"]
		monthFunc := funcs["month"]
		dayFunc := funcs["day"]

		year := callFunc(yearFunc, testTime)
		month := callFunc(monthFunc, testTime)
		day := callFunc(dayFunc, testTime)

		if year != 2023 {
			t.Errorf("expected year 2023, got %v", year)
		}
		if month != time.December {
			t.Errorf("expected December, got %v", month)
		}
		if day != 25 {
			t.Errorf("expected day 25, got %v", day)
		}
	})
}

// Test collection helpers
func TestCollectionHelpers(t *testing.T) {
	funcs := createHelperFuncs()

	t.Run("len", func(t *testing.T) {
		lenFunc := funcs["len"]

		tests := []struct {
			input    any
			expected int
		}{
			{[]any{"a", "b", "c"}, 3},
			{map[string]any{"a": 1, "b": 2}, 2},
			{"hello", 5},
			{123, 0}, // unsupported type
		}

		for _, tt := range tests {
			result := callFunc(lenFunc, tt.input)
			if result != tt.expected {
				t.Errorf("len(%v) expected %d, got %d", tt.input, tt.expected, result)
			}
		}
	})

	t.Run("first", func(t *testing.T) {
		firstFunc := funcs["first"]

		slice := []any{"a", "b", "c"}
		result := callFunc(firstFunc, slice)
		if result != "a" {
			t.Errorf("expected 'a', got %v", result)
		}

		empty := []any{}
		result = callFunc(firstFunc, empty)
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})

	t.Run("last", func(t *testing.T) {
		lastFunc := funcs["last"]

		slice := []any{"a", "b", "c"}
		result := callFunc(lastFunc, slice)
		if result != "c" {
			t.Errorf("expected 'c', got %v", result)
		}
	})

	t.Run("slice", func(t *testing.T) {
		sliceFunc := funcs["slice"]

		slice := []any{"a", "b", "c", "d", "e"}
		result := callFunc(sliceFunc, 1, 3, slice)
		expected := []any{"b", "c"}

		if !equalSlices(result.([]any), expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

// Test conditional helpers
func TestConditionalHelpers(t *testing.T) {
	funcs := createHelperFuncs()

	t.Run("default", func(t *testing.T) {
		defaultFunc := funcs["default"]

		tests := []struct {
			defaultVal any
			val        any
			expected   any
		}{
			{"default", nil, "default"},
			{"default", "", "default"},
			{"default", "value", "value"},
			{42, nil, 42},
		}

		for _, tt := range tests {
			result := callFunc(defaultFunc, tt.defaultVal, tt.val)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		}
	})

	t.Run("empty", func(t *testing.T) {
		emptyFunc := funcs["empty"]

		tests := []struct {
			input    any
			expected bool
		}{
			{nil, true},
			{"", true},
			{[]any{}, true},
			{map[string]any{}, true},
			{"hello", false},
			{[]any{"a"}, false},
			{map[string]any{"a": 1}, false},
			{42, false},
		}

		for _, tt := range tests {
			result := callFunc(emptyFunc, tt.input)
			if result != tt.expected {
				t.Errorf("empty(%v) expected %v, got %v", tt.input, tt.expected, result)
			}
		}
	})
}

// Test file path helpers
func TestFilePathHelpers(t *testing.T) {
	funcs := createHelperFuncs()

	tests := []struct {
		name     string
		funcName string
		input    string
		expected string
	}{
		{"basename", "basename", "/path/to/file.txt", "file.txt"},
		{"dirname", "dirname", "/path/to/file.txt", "/path/to"},
		{"ext", "ext", "/path/to/file.txt", ".txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := funcs[tt.funcName]
			result := callFunc(fn, tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}

	t.Run("pathjoin", func(t *testing.T) {
		pathjoinFunc := funcs["pathjoin"]
		result := callFunc(pathjoinFunc, "path", "to", "file.txt")
		expected := "path/to/file.txt"
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})
}

// Test JSON functions
func TestJSONFunctions(t *testing.T) {
	funcs := createHelperFuncs()

	t.Run("toJSON", func(t *testing.T) {
		toJSONFunc := funcs["toJSON"]

		input := map[string]any{"name": "John", "age": 30}
		result, err := callFuncWithError(toJSONFunc, input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var parsed map[string]any
		if err := json.Unmarshal([]byte(result.(string)), &parsed); err != nil {
			t.Errorf("result is not valid JSON: %v", err)
		}
	})

	t.Run("toPrettyJSON", func(t *testing.T) {
		toPrettyJSONFunc := funcs["toPrettyJSON"]

		input := map[string]any{"name": "John", "age": 30}
		result, err := callFuncWithError(toPrettyJSONFunc, input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !strings.Contains(result.(string), "\n") {
			t.Errorf("expected pretty JSON with newlines")
		}
	})
}

// Test hashing functions
func TestHashingFunctions(t *testing.T) {
	funcs := createHelperFuncs()

	tests := []struct {
		name     string
		funcName string
		input    string
		expected string
	}{
		{"sha256", "sha256", "hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"md5", "md5", "hello", "5d41402abc4b2a76b9719d911017c592"},
		{"sha1", "sha1", "hello", "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{"base64Encode", "base64Encode", "hello", "aGVsbG8="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := funcs[tt.funcName]
			result := callFunc(fn, tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}

	t.Run("base64Decode", func(t *testing.T) {
		base64DecodeFunc := funcs["base64Decode"]

		result, err := callFuncWithError(base64DecodeFunc, "aGVsbG8=")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != "hello" {
			t.Errorf("expected 'hello', got %s", result)
		}

		// Test invalid base64
		_, err = callFuncWithError(base64DecodeFunc, "invalid!")
		if err == nil {
			t.Errorf("expected error for invalid base64")
		}
	})
}

// Integration tests with Go template engine
func TestTemplateIntegration(t *testing.T) {
	funcs := createHelperFuncs()

	tests := []struct {
		name     string
		template string
		data     any
		expected string
	}{
		{
			name:     "string manipulation",
			template: `{{.name | upper}} is {{len .name}} characters long`,
			data:     map[string]any{"name": "john"},
			expected: "JOHN is 4 characters long",
		},
		{
			name:     "string lower",
			template: `{{.name | lower}} is {{len .name}} characters long`,
			data:     map[string]any{"name": "JOHN"},
			expected: "john is 4 characters long",
		},
		{
			name:     "string trim",
			template: `{{.name | trim}} is {{len (.name | trim)}} characters long`,
			data:     map[string]any{"name": "  john  "},
			expected: "john is 4 characters long",
		},
		{
			name:     "string replace",
			template: `{{.name | replace "o" "0"}}`,
			data:     map[string]any{"name": "john"},
			expected: "j0hn",
		},
		{
			name:     "string split, join",
			template: `{{ .name | split "," | join ", " }}`,
			data:     map[string]any{"name": "apple,banana,cherry"},
			expected: "apple, banana, cherry",
		},
		{
			name:     "string contains",
			template: `{{if .name | contains "ell"}}Yes{{else}}No{{end}}`,
			data:     map[string]any{"name": "hello"},
			expected: "Yes",
		},
		{
			name:     "string hasPrefix",
			template: `{{if .name | hasPrefix "he"}}Yes{{else}}No{{end}}`,
			data:     map[string]any{"name": "hello"},
			expected: "Yes",
		},
		{
			name:     "string hasSuffix",
			template: `{{if .name | hasSuffix "lo"}}Yes{{else}}No{{end}}`,
			data:     map[string]any{"name": "hello"},
			expected: "Yes",
		},
		{
			name:     "string repeat",
			template: `{{.name | repeat 3}}`,
			data:     map[string]any{"name": "ab"},
			expected: "ababab",
		},
		{
			name:     "math operations",
			template: `{{add .a .b}} and {{mulf .x .y}}`,
			data:     map[string]any{"a": 5, "b": 3, "x": 2.5, "y": 4.0},
			expected: "8 and 10",
		},
		{
			name:     "conditional and default",
			template: `{{if empty .missing}}No value{{else}}{{.missing}}{{end}} | {{default "N/A" .missing}}`,
			data:     map[string]any{},
			expected: "No value | N/A",
		},
		{
			name:     "collection operations",
			template: `First: {{first .items}}, Last: {{last .items}}, Count: {{len .items}}`,
			data:     map[string]any{"items": []any{"a", "b", "c"}},
			expected: "First: a, Last: c, Count: 3",
		},
		{
			name:     "join function",
			template: `{{join ", " .items}}`,
			data:     map[string]any{"items": []string{"apple", "banana", "cherry"}},
			expected: "apple, banana, cherry",
		},
		{
			name:     "date formatting",
			template: `Year: {{year .date}}, Formatted: {{formatDate "2006-01-02" .date}}`,
			data:     map[string]any{"date": time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)},
			expected: "Year: 2023, Formatted: 2023-12-25",
		},
		{
			name:     "file path operations",
			template: `File: {{basename .path}}, Dir: {{dirname .path}}, Ext: {{ext .path}}`,
			data:     map[string]any{"path": "/home/user/document.pdf"},
			expected: "File: document.pdf, Dir: /home/user, Ext: .pdf",
		},
		{
			name:     "hashing",
			template: `SHA256: {{sha256 .text}}, Base64: {{base64Encode .text}}`,
			data:     map[string]any{"text": "test"},
			expected: "SHA256: 9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08, Base64: dGVzdA==",
		},
		{
			name:     "range with seq",
			template: `{{range seq 1 3}}{{.}} {{end}}`,
			data:     nil,
			expected: "1 2 3 ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New("test").Funcs(funcs).Parse(tt.template)
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			var buf strings.Builder
			if err := tmpl.Execute(&buf, tt.data); err != nil {
				t.Fatalf("failed to execute template: %v", err)
			}

			result := buf.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test error handling in templates
func TestTemplateErrorHandling(t *testing.T) {
	funcs := createHelperFuncs()

	tests := []struct {
		name        string
		template    string
		data        any
		shouldError bool
	}{
		{
			name:        "division by zero",
			template:    `{{div 0 10}}`,
			data:        nil,
			shouldError: true,
		},
		{
			name:        "invalid type conversion",
			template:    `{{add "invalid" 5}}`,
			data:        nil,
			shouldError: true,
		},
		{
			name:        "invalid base64",
			template:    `{{base64Decode "invalid!"}}`,
			data:        nil,
			shouldError: true,
		},
		{
			name:        "valid operations",
			template:    `{{add 5 3}}`,
			data:        nil,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New("test").Funcs(funcs).Parse(tt.template)
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			var buf strings.Builder
			err = tmpl.Execute(&buf, tt.data)

			if tt.shouldError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// Helper functions for tests
func callFunc(fn any, args ...any) any {
	result, _ := callFuncWithError(fn, args...)
	return result
}

func callFuncWithError(fn any, args ...any) (any, error) {
	// Use reflection to call the function
	// This is a simplified version - in practice you might want more robust reflection
	switch f := fn.(type) {
	case func() time.Time:
		return f(), nil
	case func(string) string:
		return f(args[0].(string)), nil
	case func(string, string) string:
		return f(args[0].(string), args[1].(string)), nil
	case func(string, string) []string:
		return f(args[0].(string), args[1].(string)), nil
	case func(string, string) bool:
		return f(args[0].(string), args[1].(string)), nil
	case func(string, int) string:
		return f(args[0].(string), args[1].(int)), nil
	case func(string, []string) string:
		return f(args[0].(string), args[1].([]string)), nil
	case func(string, string, string) string:
		return f(args[0].(string), args[1].(string), args[2].(string)), nil
	case func(string, []any) string:
		return f(args[0].(string), args[1].([]any)), nil
	case func(int, string) string:
		return f(args[0].(int), args[1].(string)), nil
	case func(any) (float64, error):
		return f(args[0])
	case func(any) (int, error):
		return f(args[0])
	case func(any) string:
		return f(args[0]), nil
	case func(any, any) (int, error):
		return f(args[0], args[1])
	case func(any, any) (float64, error):
		return f(args[0], args[1])
	case func(any) int:
		return f(args[0]), nil
	case func([]any) []string:
		return f(args[0].([]any)), nil
	case func([]any) any:
		return f(args[0].([]any)), nil
	case func(int, int, []any) []any:
		return f(args[0].(int), args[1].(int), args[2].([]any)), nil
	case func(any, any) any:
		return f(args[0], args[1]), nil
	case func(any) bool:
		return f(args[0]), nil
	case func(string, string) (time.Time, error):
		return f(args[0].(string), args[1].(string))
	case func(string, time.Time) string:
		return f(args[0].(string), args[1].(time.Time)), nil
	case func(time.Time) int64:
		return f(args[0].(time.Time)), nil
	case func(time.Time) int:
		return f(args[0].(time.Time)), nil
	case func(time.Time) time.Month:
		return f(args[0].(time.Time)), nil
	case func(any) (string, error):
		return f(args[0])
	case func([]any) ([]string, error):
		return f(args[0].([]any))
	case func(string) (string, error):
		return f(args[0].(string))
	case func(...string) string:
		strArgs := make([]string, len(args))
		for i, arg := range args {
			strArgs[i] = arg.(string)
		}
		return f(strArgs...), nil
	default:
		return nil, fmt.Errorf("unsupported function type: %T", fn)
	}
}

func equalValues(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func equalSlices(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !equalValues(a[i], b[i]) {
			return false
		}
	}

	return true
}

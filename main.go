package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/mattn/go-isatty"
)

func toInt(s any) (int, error) {
	switch v := s.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		var i int
		_, err := fmt.Sscanf(v, "%d", &i)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string '%s' to int: %w", v, err)
		}
		return i, nil
	default:
		return 0, fmt.Errorf("unsupported type for conversion to int: %T", s)
	}
}

func toFloat(s any) (float64, error) {
	switch v := s.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		var f float64
		_, err := fmt.Sscanf(v, "%f", &f)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string '%s' to float: %w", v, err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("unsupported type for conversion to float: %T", s)
	}
}

func toFloatPair(a any, b any) (float64, float64, error) {
	aFloat, err := toFloat(a)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot convert first argument to float: %w", err)
	}

	bFloat, err := toFloat(b)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot convert second argument to float: %w", err)
	}
	return aFloat, bFloat, nil
}

func toIntPair(a any, b any) (int, int, error) {
	aInt, err := toInt(a)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot convert first argument to int: %w", err)
	}

	bInt, err := toInt(b)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot convert second argument to int: %w", err)
	}
	return aInt, bInt, nil
}

// Helper functions for the template engine
func createHelperFuncs() template.FuncMap {
	return template.FuncMap{
		// String manipulation
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"trim":  strings.TrimSpace,
		"replace": func(old, n, s string) string {
			return strings.Replace(s, old, n, -1)
		},
		"split": func(sep, str string) []string {
			return strings.Split(str, sep)
		},
		"join": func(sep string, elems []any) string {
			if len(elems) == 0 {
				return ""
			}

			// Convert []any to []string
			strElems := make([]string, len(elems))
			for i, elem := range elems {
				if str, ok := elem.(string); ok {
					strElems[i] = str
				} else {
					strElems[i] = fmt.Sprintf("%v", elem) // Fallback to string conversion
				}
			}
			// Join the string elements
			return strings.Join(strElems, sep)
		},
		"contains": func(substr string, s string) bool {
			return strings.Contains(s, substr)
		},
		"hasPrefix": func(prefix, str string) bool {
			return strings.HasPrefix(str, prefix)
		},
		"hasSuffix": func(suffix, str string) bool {
			return strings.HasSuffix(str, suffix)
		},
		"repeat": strings.Repeat,

		// Type conversion
		"toFloat": func(v any) (float64, error) {
			return toFloat(v)
		},

		"toInt": func(v any) (int, error) {
			return toInt(v)
		},

		"toString": func(v any) string {
			return fmt.Sprintf("%v", v)
		},

		// Math operations
		"add": func(a, b any) (int, error) {
			aInt, bInt, err := toIntPair(a, b)
			if err != nil {
				return 0, err
			}
			return aInt + bInt, nil
		},
		"sub": func(a, b any) (int, error) {
			aInt, bInt, err := toIntPair(a, b)
			if err != nil {
				return 0, err
			}
			return aInt - bInt, nil
		},
		"mul": func(a, b any) (int, error) {
			aInt, bInt, err := toIntPair(a, b)
			if err != nil {
				return 0, err
			}

			return aInt * bInt, nil
		},
		"div": func(a, b any) (int, error) {
			aInt, bInt, err := toIntPair(a, b)
			if err != nil {
				return 0, err
			}

			return aInt / bInt, nil
		},
		"mod": func(a, b any) (int, error) {

			aInt, bInt, err := toIntPair(a, b)
			if err != nil {
				return 0, err
			}
			return aInt % bInt, nil
		},
		// Float math operations
		"addf": func(a, b any) (float64, error) {
			aFloat, bFloat, err := toFloatPair(a, b)
			if err != nil {
				return 0, err
			}
			return aFloat + bFloat, nil
		},
		"subf": func(a, b any) (float64, error) {
			aFloat, bFloat, err := toFloatPair(a, b)
			if err != nil {
				return 0, err
			}
			return aFloat - bFloat, nil
		},
		"mulf": func(a, b any) (float64, error) {
			aFloat, bFloat, err := toFloatPair(a, b)
			if err != nil {
				return 0, err
			}
			return aFloat * bFloat, nil
		},
		"divf": func(a, b any) (float64, error) {
			aFloat, bFloat, err := toFloatPair(a, b)
			if err != nil {
				return 0, err
			}
			if bFloat == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return aFloat / bFloat, nil
		},

		// Date/time formatting
		"now": func() time.Time {
			return time.Now()
		},
		"parseDate": func(format, dateStr string) (time.Time, error) {
			t, err := time.Parse(format, dateStr)
			if err != nil {
				return time.Time{}, fmt.Errorf("failed to parse date '%s' with format '%s': %w", dateStr, format, err)
			}
			return t, nil
		},
		"formatDate": func(format string, t time.Time) string {
			return t.Format(format)
		},
		"timestamp": func(t time.Time) int64 {
			return t.Unix()
		},
		"year": func(t time.Time) int {
			return t.Year()
		},
		"month": func(t time.Time) time.Month {
			return t.Month()
		},
		"day": func(t time.Time) int {
			return t.Day()
		},

		// Collection helpers
		"len": func(v any) int {
			switch val := v.(type) {
			case []any:
				return len(val)
			case map[string]any:
				return len(val)
			case string:
				return len(val)
			default:
				return 0
			}
		},
		"first": func(v []any) any {
			if len(v) > 0 {
				return v[0]
			}
			return nil
		},
		"last": func(v []any) any {
			if len(v) > 0 {
				return v[len(v)-1]
			}
			return nil
		},
		"slice": func(start, end int, v []any) []any {
			if start < 0 || end > len(v) || start > end {
				return []any{}
			}
			return v[start:end]
		},

		// Conditional helpers
		"default": func(defaultVal, val any) any {
			if val == nil || val == "" {
				return defaultVal
			}
			return val
		},
		"empty": func(v any) bool {
			switch val := v.(type) {
			case nil:
				return true
			case string:
				return val == ""
			case []any:
				return len(val) == 0
			case map[string]any:
				return len(val) == 0
			default:
				return false
			}
		},

		// File path helpers
		"basename": filepath.Base,
		"dirname":  filepath.Dir,
		"ext":      filepath.Ext,
		"pathjoin": filepath.Join,

		// Environment variables
		"env": os.Getenv,

		// Loop helpers
		"seq": func(start, end int) iter.Seq[int] {
			return func(yield func(int) bool) {
				if start <= end {
					for i := start; i <= end; i++ {
						if !yield(i) {
							return
						}
					}
				} else {
					for i := start; i >= end; i-- {
						if !yield(i) {
							return
						}
					}
				}
			}
		},

		"toJSON": func(v any) (string, error) {
			data, err := json.Marshal(v)
			if err != nil {
				return "", fmt.Errorf("failed to marshal to JSON: %w", err)
			}
			return string(data), nil
		},

		"toPrettyJSON": func(v any) (string, error) {
			data, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				return "", fmt.Errorf("failed to marshal to pretty JSON: %w", err)
			}
			return string(data), nil
		},

		// hashing functions
		"sha256": func(s string) string {
			hash := sha256.Sum256([]byte(s))
			return fmt.Sprintf("%x", hash)
		},

		"md5": func(s string) string {
			hash := md5.Sum([]byte(s))
			return fmt.Sprintf("%x", hash)
		},

		"sha1": func(s string) string {
			hash := sha1.Sum([]byte(s))
			return fmt.Sprintf("%x", hash)
		},

		"base64Encode": func(s string) string {
			return base64.StdEncoding.EncodeToString([]byte(s))
		},

		"base64Decode": func(s string) (string, error) {
			data, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				return "", fmt.Errorf("failed to decode base64: %w", err)
			}
			return string(data), nil
		},
	}
}

func showHelp() {
	fmt.Printf(`tplsub - A Go template processor with JSON data input

USAGE:
    %s [OPTIONS] <template-file> [data-file]
    %s [OPTIONS] -t <template-string> [data-file]

OPTIONS:
    -h, --help              Show this help message
    -t, --template <string> Use template string instead of file

ARGUMENTS:
    <template-file>         Path to the Go template file
    <template-string>       Template string to execute directly
    [data-file]             Optional JSON file containing template data
                           If not provided, data is read from stdin

DATA INPUT:
    1. From file:    %s template.tmpl data.json
    2. From stdin:   echo '{"name":"John"}' | %s template.tmpl
    3. No data:      %s -t 'Hello {{ env "USER" }}'

EXAMPLES:
    # Basic template with JSON data
    echo '{"name":"John","age":30}' | %s -t 'Hello {{ .name }}, age {{ .age }}'

    # String manipulation
    echo '{"text":"hello world"}' | %s -t '{{ .text | upper | replace "WORLD" "GO" }}'

    # Math operations
    echo '{"a":10,"b":3}' | %s -t 'Sum: {{ add .a .b }}, Float: {{ divf .a .b }}'

    # Date/time functions
    %s -t 'Now: {{ now | formatDate "2006-01-02 15:04:05" }}'

    # File operations
    echo '{"path":"/home/user/doc.txt"}' | %s -t 'File: {{ basename .path }}'

    # Hashing and encoding
    echo '{"text":"hello"}' | %s -t 'Hash: {{ sha256 .text }}'

AVAILABLE FUNCTIONS:
    String:     upper, lower, trim, split, join, contains, replace, repeat
    Math:       add, sub, mul, div, mod (integers)
    Float:      addf, subf, mulf, divf, toFloat
    Date:       now, parseDate, formatDate, timestamp, year, month, day
    Collection: len, first, last, slice, seq
    Condition:  default, empty
    File:       basename, dirname, ext, pathjoin
    System:     env
    JSON:       toJSON, toPrettyJSON
    Hash:       md5, sha1, sha256, base64Encode, base64Decode
    Convert:    toString, toInt, toFloat

For detailed documentation and more examples, visit:
https://gitea.lorien.space/Ajnasz/tplsub

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	var templateContent, dataFile string

	// Check for help flag
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" {
			showHelp()
			os.Exit(0)
		}
	}

	// Parse command-line arguments
	switch {
	case len(os.Args) < 2:
		fmt.Fprintf(os.Stderr, "Usage: %s [-t template_string | template_file] [data_file]\n", os.Args[0])
		os.Exit(1)
	case os.Args[1] == "-t" || os.Args[1] == "--template":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: template string is missing after %s\n", os.Args[1])
			os.Exit(1)
		}
		templateContent = os.Args[2]
		if len(os.Args) > 3 {
			dataFile = os.Args[3]
		}
	default:
		templateFile := os.Args[1]
		content, err := os.ReadFile(templateFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading template file: %v", err)
			os.Exit(1)
		}
		templateContent = string(content)
		if len(os.Args) > 2 {
			dataFile = os.Args[2]
		}
	}

	// Load data
	var data any
	var dataReader io.Reader = os.Stdin

	if dataFile != "" {
		file, err := os.Open(dataFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening data file: %v", err)
			os.Exit(1)
		}
		defer file.Close()
		dataReader = file
	}

	if dataFile == "" && isatty.IsTerminal(os.Stdin.Fd()) {
		data = make(map[string]any)
	} else {
		decoder := json.NewDecoder(dataReader)
		if err := decoder.Decode(&data); err != nil {
			// Allow empty data if stdin is a TTY and no data is piped
			if dataFile == "" {
				f, ok := dataReader.(*os.File)
				if ok && isatty.IsTerminal(f.Fd()) {
					data = make(map[string]any)
				} else if err == io.EOF {
					data = make(map[string]any)
				} else {
					fmt.Fprintf(os.Stderr, "Error reading JSON data: %v", err)
					os.Exit(1)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Error reading JSON data from %s: %v", dataFile, err)
				os.Exit(1)
			}
		}
	}

	// Create and execute template
	tmpl, err := template.New("gotpl").Funcs(createHelperFuncs()).Parse(templateContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template: %v", err)
		os.Exit(1)
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v", err)
		os.Exit(1)
	}
}

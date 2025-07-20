package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/mattn/go-isatty"
)

// Helper functions for the template engine
func createHelperFuncs() template.FuncMap {
	return template.FuncMap{
		// String manipulation
		"upper":     strings.ToUpper,
		"lower":     strings.ToLower,
		"trim":      strings.TrimSpace,
		"trimLeft":  strings.TrimLeft,
		"trimRight": strings.TrimRight,
		"replace":   strings.Replace,
		"split":     strings.Split,
		"join":      strings.Join,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"repeat":    strings.Repeat,

		// Type conversion
		"toString": func(v any) string {
			return fmt.Sprintf("%v", v)
		},

		// Math operations
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
		"seq": func(start, end int) []int {
			var result []int
			if start <= end {
				for i := start; i <= end; i++ {
					result = append(result, i)
				}
			} else {
				for i := start; i >= end; i-- {
					result = append(result, i)
				}
			}
			return result
		},
	}
}

func main() {
	var templateContent, dataFile string

	// Parse command-line arguments
	switch {
	case len(os.Args) < 2:
		log.Fatalf("Usage: %s [-t template_string | template_file] [data_file]\n", os.Args[0])
	case os.Args[1] == "-t" || os.Args[1] == "--template":
		if len(os.Args) < 3 {
			log.Fatalf("Error: template string is missing after %s\n", os.Args[1])
		}
		templateContent = os.Args[2]
		if len(os.Args) > 3 {
			dataFile = os.Args[3]
		}
	default:
		templateFile := os.Args[1]
		content, err := os.ReadFile(templateFile)
		if err != nil {
			log.Fatalf("Error reading template file: %v", err)
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
			log.Fatalf("Error opening data file: %v", err)
		}
		defer file.Close()
		dataReader = file
	}

	decoder := json.NewDecoder(dataReader)
	if err := decoder.Decode(&data); err != nil {
		// Allow empty data if stdin is a TTY and no data is piped
		if dataFile == "" {
			if f, ok := dataReader.(*os.File); ok && isatty.IsTerminal(f.Fd()) {
				data = make(map[string]any)
			} else if err == io.EOF {
				data = make(map[string]any)
			} else {
				log.Fatalf("Error reading JSON data: %v", err)
			}
		} else {
			log.Fatalf("Error reading JSON data from %s: %v", dataFile, err)
		}
	}

	// Create and execute template
	tmpl, err := template.New("gotpl").Funcs(createHelperFuncs()).Parse(templateContent)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

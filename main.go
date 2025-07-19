package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
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
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a % b
		},
		"max": func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},

		// Comparison
		"eq": func(a, b any) bool { return a == b },
		"ne": func(a, b any) bool { return a != b },
		"lt": func(a, b int) bool { return a < b },
		"le": func(a, b int) bool { return a <= b },
		"gt": func(a, b int) bool { return a > b },
		"ge": func(a, b int) bool { return a >= b },

		// Date/time formatting
		"now": func() time.Time {
			return time.Now()
		},
		"formatDate": func(format string, t time.Time) string {
			return t.Format(format)
		},
		"timestamp": func() int64 {
			return time.Now().Unix()
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
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <template-file> [data-file]\n", os.Args[0])
	}

	fileName := os.Args[1]

	// Check if the file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		panic("File does not exist: " + fileName)
	}

	// Load data from JSON file if provided
	var data any
	if len(os.Args) > 2 {
		dataFile := os.Args[2]
		if content, err := os.ReadFile(dataFile); err == nil {
			json.Unmarshal(content, &data)
		} else {
			log.Printf("Warning: Could not read data file %s: %v\n", dataFile, err)
		}
	}

	// Create template with helper functions
	tmpl := template.New(filepath.Base(fileName)).Funcs(createHelperFuncs())
	tmpl = template.Must(tmpl.ParseFiles(fileName))

	if err := tmpl.ExecuteTemplate(os.Stdout, filepath.Base(fileName), data); err != nil {
		log.Fatalf("Error executing template: %v\n", err)
	}
}

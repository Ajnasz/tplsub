package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func mmap[T any](m []any, f func(any) (T, error)) ([]T, error) {
	result := make([]T, len(m))
	for i, v := range m {
		res, err := f(v)
		if err != nil {
			return nil, fmt.Errorf("error applying function to element %d: %w", i, err)
		}
		result[i] = res
	}
	return result, nil
}

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

func toString(x any) (string, error) {
	if x == nil {
		return "", nil
	}
	switch v := x.(type) {
	case string:
		return v, nil
	case int, int64, float64:
		return fmt.Sprintf("%v", v), nil
	default:
		return fmt.Sprintf("%v", x), nil
	}
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
		"join": func(sep string, elems []string) string {
			return strings.Join(elems, sep)
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
		"repeat": func(count int, s string) string {
			return strings.Repeat(s, count)
		},

		// Type conversion
		"toFloat": func(v any) (float64, error) {
			return toFloat(v)
		},
		"toInt": func(v any) (int, error) {
			return toInt(v)
		},
		"toString": func(v any) string {
			s, _ := toString(v)
			return s
		},
		"toStrings": func(v []any) ([]string, error) {
			return mmap(v, toString)
		},
		"toInts": func(v []any) ([]int, error) {
			return mmap(v, toInt)
		},
		"toFloats": func(v []any) ([]float64, error) {
			return mmap(v, toFloat)
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

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/mattn/go-isatty"
)

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
    Convert:    toString, toStrings, toInt, toInts, toFloat, toFloats

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
				if err == io.EOF {
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
	if err := executeTemplate(os.Stdout, templateContent, data); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func executeTemplate(out io.Writer, templateContent string, data any) error {
	tmpl, err := template.New("gotpl").Funcs(createHelperFuncs()).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	if err := tmpl.Execute(out, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

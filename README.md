# tplsub

A simple Go command-line tool for executing Go text templates with JSON data input.

## Description

`tplsub` is a lightweight utility that takes a Go template and JSON data as input, executes the template with the provided data, and outputs the result to stdout. It's useful for template processing and text substitution tasks with dynamic data.

## Installation

```bash
go install gitea.lorien.space/Ajnasz/tplsub@latest
```

Or clone and build locally:

```bash
git clone https://gitea.lorien.space/Ajnasz/tplsub.git
cd tplsub
go build
```

## Usage

```bash
# Using template file
tplsub <template-file> [data-file]

# Using inline template
tplsub -t <template-string> [data-file]
tplsub --template <template-string> [data-file]
```

### Arguments

- `<template-file>`: Path to the Go template file to execute
- `-t, --template <template-string>`: Template string to execute directly
- `[data-file]`: Optional JSON file containing template data. If not provided, data is read from stdin

### Data Input

You can provide JSON data in three ways:

1. **From a file**: `tplsub template.tmpl data.json`
2. **From stdin**: `echo '{"name": "John"}' | tplsub template.tmpl`
3. **No data**: `tplsub template.tmpl` (empty data object will be used)

### Examples

```bash
# Execute a template file with data from stdin
echo '{"FirstName": "John", "LastName": "Doe"}' | tplsub template.tmpl

# Execute inline template with data from stdin
echo '{"FirstName": "John", "LastName": "Doe"}' | tplsub -t 'Hello {{ .FirstName }} {{ .LastName }}'

# Execute template with data from file
tplsub template.tmpl data.json

# Execute template without data
tplsub -t 'Current time: {{ now | formatDate "2006-01-02 15:04:05" }}'
```

## Template Format

The tool uses Go's `text/template` package with additional helper functions. Your template files should follow the standard Go template syntax.

Example template file (`example.tmpl`):
```
Hello {{ .FirstName }} {{ .LastName | upper }}!
Your data hash: {{ . | toJSON | md5 }}
```

## Available Helper Functions

### String Manipulation
- `upper` - Convert to uppercase: `{{ "hello" | upper }}` → `HELLO`
- `lower` - Convert to lowercase: `{{ "HELLO" | lower }}` → `hello`
- `trim` - Trim whitespace: `{{ " hello " | trim }}` → `hello`
- `replace` - Replace strings: `{{ replace "hello" "hi" "hello world" }}` → `hi world`
- `split` - Split string: `{{ split "," "a,b,c" }}` → `[a b c]`
- `join` - Join array: `{{ join "," (split ";" "a;b;c") }}` → `a,b,c`
- `contains` - Check if contains: `{{ contains "ell" "hello" }}` → `true`
- `hasPrefix` - Check prefix: `{{ hasPrefix "he" "hello" }}` → `true`
- `hasSuffix` - Check suffix: `{{ hasSuffix "lo" "hello" }}` → `true`
- `repeat` - Repeat string: `{{ repeat 3 "hi" }}` → `hihihi`

### Type Conversion
- `toString` - Convert to string: `{{ toString 123 }}` → `123`
- `toInt` - Convert to int: `{{ toInt "123" }}` → `123`
- `toFloat` - Convert to float: `{{ toFloat "3.14" }}` → `3.14`
- `toStrings` - Convert array to strings: `{{ toStrings [1 2 3] }}` → `["1" "2" "3"]`
- `toInts` - Convert array to ints: `{{ toInts ["1" "2" "3"] }}` → `[1 2 3]`
- `toFloats` - Convert array to floats: `{{ toFloats ["1.1" "2.2" "3.3"] }}` → `[1.1 2.2 3.3]`

### Math Operations
- `add` - Addition: `{{ add 5 3 }}` → `8`
- `sub` - Subtraction: `{{ sub 5 3 }}` → `2`
- `mul` - Multiplication: `{{ mul 5 3 }}` → `15`
- `div` - Division: `{{ div 6 3 }}` → `2`
- `mod` - Modulo: `{{ mod 7 3 }}` → `1`

### Float Math Operations
- `addf` - Float addition: `{{ addf 5.5 3.2 }}` → `8.7`
- `subf` - Float subtraction: `{{ subf 5.5 3.2 }}` → `2.3`
- `mulf` - Float multiplication: `{{ mulf 5.5 3.2 }}` → `17.6`
- `divf` - Float division: `{{ divf 22 7 }}` → `3.142857142857143`

### Date/Time Functions
- `now` - Current time: `{{ now }}`
- `parseDate` - Parse date: `{{ parseDate "2006-01-02" "2023-12-25" }}`
- `formatDate` - Format date: `{{ now | formatDate "2006-01-02 15:04:05" }}`
- `timestamp` - Unix timestamp: `{{ now | timestamp }}`
- `year` - Get year: `{{ now | year }}`
- `month` - Get month: `{{ now | month }}`
- `day` - Get day: `{{ now | day }}`

### Collection Helpers
- `len` - Get length: `{{ len .items }}` or `{{ len "hello" }}` → `5`
- `first` - First element: `{{ first .items }}`
- `last` - Last element: `{{ last .items }}`
- `slice` - Slice array: `{{ slice 1 3 .items }}`

### Conditional Helpers
- `default` - Default value: `{{ default "N/A" .name }}`
- `empty` - Check if empty: `{{ if empty .name }}No name{{ end }}`

### File Path Helpers
- `basename` - Get basename: `{{ basename "/path/to/file.txt" }}` → `file.txt`
- `dirname` - Get directory: `{{ dirname "/path/to/file.txt" }}` → `/path/to`
- `ext` - Get extension: `{{ ext "file.txt" }}` → `.txt`
- `pathjoin` - Join paths: `{{ pathjoin "/path" "to" "file.txt" }}` → `/path/to/file.txt`

### Environment Variables
- `env` - Get environment variable: `{{ env "HOME" }}`

### Loop Helpers
- `seq` - Generate sequence: `{{ range seq 1 5 }}{{ . }}{{ end }}` → `12345`

### JSON Functions
- `toJSON` - Convert to JSON: `{{ .data | toJSON }}`
- `toPrettyJSON` - Convert to pretty JSON: `{{ .data | toPrettyJSON }}`

### Hashing and Encoding
- `sha256` - SHA256 hash: `{{ sha256 "hello" }}`
- `md5` - MD5 hash: `{{ md5 "hello" }}`
- `sha1` - SHA1 hash: `{{ sha1 "hello" }}`
- `base64Encode` - Base64 encode: `{{ base64Encode "hello" }}`
- `base64Decode` - Base64 decode: `{{ base64Decode "aGVsbG8=" }}`

## Error Handling

- If no template is provided, the program will exit with usage information
- If the specified template file doesn't exist, the program will log an error and exit
- If there's an error parsing or executing the template, the program will log the error and exit
- If JSON data is malformed, the program will log an error and exit

## Requirements

- Go 1.24.3 or later

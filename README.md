# gotplsubst

A simple Go command-line tool for executing Go text templates.

## Description

`gotplsubst` is a lightweight utility that takes a Go template file as input and executes it, outputting the result to stdout. It's useful for template processing and text substitution tasks.

## Installation

```bash
go install gitea.lorien.space/Ajnasz/gotplsubst@latest
```

Or clone and build locally:

```bash
git clone https://gitea.lorien.space/Ajnasz/gotplsubst.git
cd gotplsubst
go build
```

## Usage

```bash
gotplsubst <template-file>
```

### Arguments

- `<template-file>`: Path to the Go template file to execute

### Example

```bash
# Execute a template file
gotplsubst my-template.tmpl

# Redirect output to a file
gotplsubst my-template.tmpl > output.txt
```

## Template Format

The tool uses Go's `text/template` package. Your template files should follow the standard Go template syntax.

Example template file (`example.tmpl`):
```
{{$name := "John"}}

The name is {{$name}}
```

## Error Handling

- If no template file is provided, the program will exit with usage information
- If the specified template file doesn't exist, the program will panic with an error message
- If there's an error executing the template, the program will log the error and exit

## Requirements

- Go 1.24.3 or later

## License

[Add your license information here]

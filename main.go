package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <template-file>\n", os.Args[0])
	}

	fileName := os.Args[1]

	// Check if the file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		panic("File does not exist: " + fileName)
	}
	templ := template.Must(template.New(fileName).ParseFiles(fileName))
	if err := templ.ExecuteTemplate(os.Stdout, fileName, nil); err != nil {
		log.Fatalf("Error executing template: %v\n", err)
	}
}

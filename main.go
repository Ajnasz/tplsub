package main

import (
	"os"
	"text/template"
)

func main() {
	fileName := "example.tmpl"
	templ := template.Must(template.New(fileName).ParseFiles(fileName))
	if err := templ.ExecuteTemplate(os.Stdout, "example.tmpl", nil); err != nil {
		panic(err)
	}
}

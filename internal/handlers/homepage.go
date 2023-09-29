package handlers

import (
	"net/http"
	"text/template"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./html/index.page.tmpl",
		"./html/base.layout.tmpl",
	}

	home := template.Must(template.ParseFiles(files...))

	home.Execute(w, nil)
}

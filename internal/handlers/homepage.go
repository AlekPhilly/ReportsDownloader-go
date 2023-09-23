package handlers

import (
	"net/http"
	"os"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	site, _ := os.ReadFile("./html/main.html")

	w.Write(site)
}

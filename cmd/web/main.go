package main

import (
	"log"
	"net/http"

	"github.com/AlekPhilly/ReportsDownloader-go/internal/handlers"
)

const PORT = ":8080"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomePage)
	mux.HandleFunc("/search", handlers.Search)
	mux.HandleFunc("/download", handlers.Download)

	srv := &http.Server{
		Addr:    PORT,
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}

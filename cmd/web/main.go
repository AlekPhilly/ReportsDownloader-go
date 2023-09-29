package main

import (
	"fmt"
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

	fileSrv := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileSrv))

	srv := &http.Server{
		Addr:    PORT,
		Handler: mux,
	}

	fmt.Printf("Starting web server at 127.0.0.1%s...\n", PORT)
	log.Fatal(srv.ListenAndServe())
}

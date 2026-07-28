package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API Key Service")
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "healthy")
	})

	router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ready")
	})

	log.Println("server listening on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

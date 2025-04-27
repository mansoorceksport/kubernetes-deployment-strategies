package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	version := os.Getenv("VERSION")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from version %s!", version)
	})
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}

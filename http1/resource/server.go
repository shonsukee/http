package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request:", r.Method, r.URL.Path, r.URL.RawQuery)
	fmt.Println("Received query parameters:", r.URL.Query())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Server", "Go HTTP Server")
	w.Header().Set("X-Custom-Header", "CustomValue")
	w.WriteHeader(http.StatusOK)
	// Write the response
	fmt.Fprintf(w, "<h1>hello world</h1>")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server starting on port 18888...")
	log.Fatal(http.ListenAndServe(":18888", nil))
}

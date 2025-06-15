package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", HandleVisited)
	http.HandleFunc("/form", HandleXmlHttpForm)
	http.HandleFunc("/upload", HandleMultipartUpload)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

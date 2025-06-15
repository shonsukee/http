package main

import (
	"encoding/xml"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HandleXmlHttpForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Name  string `xml:"name"`
		Value string `xml:"value"`
	}
	if err := xml.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("uploads", strings.ReplaceAll(data.Name, "/", "_"))
	if err := os.WriteFile(filePath, []byte(data.Value), 0644); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// curl -X POST -H "Content-Type: application/xml" -d '<?xml version="1.0" encoding="UTF-8"?><data><name>example.txt</name><value>これはテストデータです</value></data>' http://localhost:8080/form

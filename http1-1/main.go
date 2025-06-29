package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// This is a placeholder function to demonstrate the structure.
	// Actual implementation would handle HTTP requests.
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, "Error dumping request", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))

	// ファイル名を取得
	file_path := r.FormValue("file")
	if file_path == "" {
		file_path = "http1-1.txt"
	}

	// ファイルの内容を読み込む
	content, err := os.ReadFile(file_path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file %s: %v", file_path, err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File: %s\nContent:\n%s", file_path, string(content))
}

func main() {
	// 各パスに対応するハンドラを登録
	http.HandleFunc("/", handler)

	fmt.Println("Server running on https://localhost:8080")
	err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

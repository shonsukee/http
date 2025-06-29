package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// HandleMultipartUpload はmultipart/form-dataのアップロードを処理します
func HandleMultipartUpload(w http.ResponseWriter, r *http.Request) {
	// multipart form をパース
	err := r.ParseMultipartForm(10 << 20) // 10 MB のメモリ制限
	if err != nil {
		log.Println("Error parsing multipart form:", err)
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	// "attachment-file" フィールドからファイルを取得
	file, header, err := r.FormFile("attachment-file")
	if err != nil {
		log.Println("Error getting form file:", err)
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Println("Uploaded File:", header.Filename)
	log.Println("File Size:", header.Size)
	log.Println("MIME Header:", header.Header)

	// ファイルの中身を表示
	content, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading file:", err)
		http.Error(w, "Could not read file", http.StatusInternalServerError)
		return
	}
	log.Println("File content:\n" + string(content))

	fmt.Fprintf(w, "File %s received.\nContent:\n%s", header.Filename, content)
}

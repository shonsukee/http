package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

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
	log.Println("File content:\n" + string(content) + "\n")

	fmt.Fprintf(w, "File %s received.\nContent:\n%s", header.Filename, content)
}

// Terminal logs
// % curl --http1.0 -F attachment-file=@http1.txt http://localhost:8080/upload
// File http1.txt received.
// Content:
//     __  __________________ ___ ____     ______           _       _
//    / / / /_  __/_  __/ __ <  // __ \   /_  __/________ _(_)___  (_)___  ____ _
//   / /_/ / / /   / / / /_/ / // / / /    / / / ___/ __ `/ / __ \/ / __ \/ __ `/
//  / __  / / /   / / / ____/ // /_/ /    / / / /  / /_/ / / / / / / / / / /_/ /
// /_/ /_/ /_/   /_/ /_/   /_(_)____/    /_/ /_/   \__,_/_/_/ /_/_/_/ /_/\__, /
//                                                                      /____/

// docker logs
// http1  | 2025/06/14 06:47:02 Uploaded File: http1.txt
// http1  | 2025/06/14 06:47:02 File Size: 473
// http1  | 2025/06/14 06:47:02 MIME Header: map[Content-Disposition:[form-data; name="attachment-file"; filename="http1.txt"] Content-Type:[text/plain]]
// http1  | 2025/06/14 06:47:02 File content:
// http1  |     __  __________________ ___ ____     ______           _       _
// http1  |    / / / /_  __/_  __/ __ <  // __ \   /_  __/________ _(_)___  (_)___  ____ _
// http1  |   / /_/ / / /   / / / /_/ / // / / /    / / / ___/ __ `/ / __ \/ / __ \/ __ `/
// http1  |  / __  / / /   / / / ____/ // /_/ /    / / / /  / /_/ / / / / / / / / / /_/ /
// http1  | /_/ /_/ /_/   /_/ /_/   /_(_)____/    /_/ /_/   \__,_/_/_/ /_/_/_/ /_/\__, /
// http1  |                                                                      /____/

package main

import (
	"io"
	"log"
	"net/http"
)

func HandleGetMethod(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	log.Println("Request body:", string(body))
}

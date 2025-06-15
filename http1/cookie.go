package main

import (
	"log"
	"net/http"
)

func HandleVisited(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Set-Cookie", "VISIT=TRUE")
	if _, ok := r.Header["Cookie"]; ok {
		log.Println("おかえりなさい")
	} else {
		log.Println("ようこそ！")
	}
	http.ServeFile(w, r, "index.html")
}

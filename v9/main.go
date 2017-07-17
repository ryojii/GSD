package main

import (
	"log"
	"net/http"
)

func main() {
    log.Printf("Serveur started")
	Init()

	router := NewRouter()
    s := http.StripPrefix("/js/", http.FileServer(http.Dir("./js/")))
    router.PathPrefix("/js/").Handler(s)
    http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8188", router))
}

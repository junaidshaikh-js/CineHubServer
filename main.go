package main

import (
	"log"
	"net/http"
	"time"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("public")))

	const addr = ":5555"

	s := http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on %s", addr)

	err := s.ListenAndServe()

	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

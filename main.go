package main

import (
	"log"
	"net/http"
)

func main() {
	templ, err := loadTemplates()
	if err != nil {
		log.Fatalf("Can not load templates: %s", err)
	}

	questions, err := readDatabase()
	if err != nil {
		log.Fatalf("Can not read database: %s", err)
	}

	mux := setupHandlers(templ, questions)

	addr := ":8080"
	log.Printf("Listening on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

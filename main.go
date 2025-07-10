package main

import (
	"log"
	"net/http"

	"local-ai/handlers"
)

func main() {
	http.HandleFunc("/chat", handlers.ChatHandler)
	log.Println("Codex API running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

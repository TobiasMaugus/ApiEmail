package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/contato", Handler)

	// Porta que o Render fornece (ou fallback para 8080 localmente)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado na porta %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

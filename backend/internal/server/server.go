package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StartServer(router *chi.Mux, port string) {

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Printf("Server running on port %s\n", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server error:", err)
	}
}

package main

import (
	"log"

	"github.com/didrikolofsson/materials/internal/config"
	infrahttp "github.com/didrikolofsson/materials/internal/infra/http"
)

func main() {
	cfg := config.Load()
	srv := infrahttp.New(cfg.Port)

	if err := srv.Run(); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

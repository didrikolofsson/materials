package main

import (
	"log"

	"github.com/didrikolofsson/materials/internal/config"
	infrahttp "github.com/didrikolofsson/materials/internal/infra/http"
	infradb "github.com/didrikolofsson/materials/internal/infra/mysql"
)

func main() {
	cfg := config.Load()

	db, err := infradb.New(cfg.DBDsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	srv := infrahttp.New(cfg.Port, db)

	if err := srv.Run(); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

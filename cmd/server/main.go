package main

import (
	"log"

	"github.com/didrikolofsson/materials/internal/config"
	"github.com/didrikolofsson/materials/internal/domain/school"
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

	subjectsRepo := school.NewMySQLSubjectsRepository(db)
	subjectsHandler := school.NewSubjectsHandler(subjectsRepo)

	srv := infrahttp.New(cfg.Port, subjectsHandler)

	if err := srv.Run(); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

package main

import (
	"database/sql"
	"log"

	"github.com/didrikolofsson/materials/internal/config"
	infrahttp "github.com/didrikolofsson/materials/internal/infra/http"
	infradb "github.com/didrikolofsson/materials/internal/infra/mysql"
	"github.com/go-playground/validator/v10"
)

type Dependencies struct {
	DB       *sql.DB
	Validate *validator.Validate
}

func InitDependencies(cfg config.Config) Dependencies {
	db, err := infradb.New(cfg.DBDsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return Dependencies{
		DB:       db,
		Validate: validator.New(),
	}
}

func main() {
	cfg := config.Load()

	deps := InitDependencies(cfg)
	defer deps.DB.Close()

	srv := infrahttp.New(
		cfg.Port, deps.DB, deps.Validate,
	)

	if err := srv.Run(); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

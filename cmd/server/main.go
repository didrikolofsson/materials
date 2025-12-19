package main

import (
	"log"

	"github.com/didrikolofsson/materials/internal/config"
	"github.com/didrikolofsson/materials/internal/handlers"
	infrahttp "github.com/didrikolofsson/materials/internal/infra/http"
	"github.com/didrikolofsson/materials/internal/infra/mysql"
	"github.com/didrikolofsson/materials/internal/repositories"
	"github.com/didrikolofsson/materials/internal/services"
	"github.com/go-playground/validator/v10"
)

func main() {
	cfg := config.Load()

	db := mysql.New(cfg.DBDsn)
	validate := validator.New()

	repos := repositories.New(db)
	svc := services.New(db, repos)
	handlers := handlers.New(svc, validate)

	defer db.Close()

	srv := infrahttp.New(
		cfg.Port,
		handlers,
	)

	if err := srv.Run(); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

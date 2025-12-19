package main

import (
	"log"

	"github.com/didrikolofsson/materials/internal/config"
	"github.com/didrikolofsson/materials/internal/infra/mysql"
	"github.com/didrikolofsson/materials/internal/seed"
)

func main() {
	cfg := config.Load()
	db := mysql.New(cfg.DBDsn)

	if err := seed.Run(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	log.Println("Database seeded successfully!")
	db.Close()
}

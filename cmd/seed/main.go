package main

import (
	"context"
	"log"

	"github.com/didrikolofsson/materials/internal/config"
	"github.com/didrikolofsson/materials/internal/infra/mysql"
	"github.com/didrikolofsson/materials/internal/seed"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()
	db := mysql.New(cfg.DBDsn)

	if err := seed.Run(ctx, db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	log.Println("Database seeded successfully!")
	db.Close()
}

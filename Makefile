DB_URL="mysql://root:root@tcp(localhost:3306)/materials?parseTime=true"
MIGRATIONS_DIR="migrations"
DOCKER_COMPOSE_FILE="deployments/docker-compose.yml"

# Commands
run:
	go run cmd/server/main.go

# DB commands
migrate:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) down 1

migrate-reset:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) reset

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) $(name)

# Docker commands
docker-up:
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(DOCKER_COMPOSE_FILE) down

docker-restart:
	docker compose -f $(DOCKER_COMPOSE_FILE) restart

docker-logs:
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f
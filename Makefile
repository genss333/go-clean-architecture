.PHONY: docker-up docker-down sqlc-gen

COMPOSE_FILE := deployments/docker-compose.yml

docker-up:
	docker compose -f $(COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(COMPOSE_FILE) down

sqlc-gen:
	@echo "Generate SQLC (PostgreSQL)..."
	@cd db/sqlc && sqlc generate
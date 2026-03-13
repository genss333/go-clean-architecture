.PHONY: docker-up docker-down sqlc-gen test-unit test-integration

COMPOSE_FILE := deployments/docker-compose.yml

docker-up:
	docker compose -f $(COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(COMPOSE_FILE) down

sqlc-gen:
	@echo "Generate SQLC (PostgreSQL)..."
	@cd db/sqlc && sqlc generate

test-unit:
	go test -v -race ./test/unit-test/...

test-integration:
	go test -v -timeout 120s ./test/integration/...
.PHONY: docker-up docker-down sqlc-gen test-unit test-db-integration test-api-integration bench-unit bench-db-integration bench-api-integration

COMPOSE_FILE := deployments/docker-compose.yml

docker-up:
	docker compose -f $(COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(COMPOSE_FILE) down

sqlc-gen:
	@echo "Generate SQLC (PostgreSQL)..."
	@cd db/sqlc && sqlc generate

test-unit:
	go test -v -race -coverprofile=coverage-unit.out ./test/unit-test/...

test-db-integration:
	go test -v -timeout 120s -coverprofile=coverage-db.out ./test/database_integration/...

test-api-integration:
	go test -v -timeout 120s -coverprofile=coverage-api.out ./test/api_integration/...

bench-unit:
	go test -bench=. -benchmem -timeout 120s ./test/unit-test/...

bench-db-integration:
	go test -bench=. -benchmem -timeout 120s ./test/database_integration/...

bench-api-integration:
	go test -bench=. -benchmem -timeout 120s ./test/api_integration/...
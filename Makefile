.PHONY: docker-up docker-down sqlc-gen test-unit test-integration bench bench-unit bench-handler bench-integration

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

bench:
	go test -bench=. -benchtime=5x -benchmem ./test/unit-test/... ./test/api_test/...

bench-unit:
	go test -bench=. -benchtime=5x -benchmem ./test/unit-test/...

bench-handler:
	go test -bench=. -benchtime=5x -benchmem ./test/api_test/...

bench-integration:
	go test -bench=. -benchtime=5x -benchmem -timeout 120s ./test/integration/...
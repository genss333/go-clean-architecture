.PHONY: docker-up docker-down sqlc-gen test-unit test-integration test-api-integration bench bench-unit bench-handler bench-integration bench-api-integration

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
	go test -v -timeout 120s ./test/database_integration/...

bench:
	go test -bench=. -benchmem -timeout 120s ./test/unit-test/... ./test/api_integration/...

bench-unit:
	go test -bench=. -benchmem -timeout 120s ./test/unit-test/...

bench-handler:
	go test -bench=. -benchmem -timeout 120s ./test/api_test/...

test-api-integration:
	go test -v -timeout 120s ./test/api_integration/...

bench-integration:
	go test -bench=. -benchmem -timeout 120s ./test/database_integration/...

bench-api-integration:
	go test -bench=. -benchmem -timeout 120s ./test/api_integration/...
AUTH_SERVICE=services/auth-service
TENANT_SERVICE=services/tenant-service
LOCATION_SERVICE=services/location-service
STREAMING_SERVICE=services/streaming-service

include .env
export $(shell sed 's/=.*//' .env)

.PHONY: all up down build migrate run-auth run-tenant run-location run-streaming logs clean psql

# Start all services using Docker Compose
up:
	cd deployments && docker-compose up --build

# Stop all services
down:
	cd deployments && docker-compose down

# Build all Go binaries (optional)
build:
	cd $(AUTH_SERVICE) && go build -o ../../bin/auth-service
	cd $(TENANT_SERVICE) && go build -o ../../bin/tenant-service
	cd $(LOCATION_SERVICE) && go build -o ../../bin/location-service
	cd $(STREAMING_SERVICE) && go build -o ../../bin/streaming-service

psql:
	docker exec -it $$(docker ps --filter "ancestor=postgres" --format "{{.ID}}") \
	bash -c "psql -U postgres -d multi_tenant_db"



# Run database migrations
migrate:
	docker run --rm \
		--network deployments_default \
		-v $(PWD)/migrations:/migrations \
		--env-file scripts/.env \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://$$DB_USER:$$DB_PASSWORD@db:5432/$$DB_NAME?sslmode=disable" up


# Run services individually (for development)
run-auth:
	cd $(AUTH_SERVICE) && go run main.go

run-tenant:
	cd $(TENANT_SERVICE) && go run main.go

run-location:
	cd $(LOCATION_SERVICE) && go run main.go

run-streaming:
	cd $(STREAMING_SERVICE) && go run main.go

# View logs from all services
logs:
	cd deployments && docker-compose logs -f

# Clean up built binaries and Docker containers/volumes
clean:
	rm -rf bin/*
	cd deployments && docker-compose down -v

# (Optional) Health check stubs (implement /health endpoints in your services for these to work)
check-auth:
	curl -i http://localhost:8081/health || echo "No /health endpoint implemented"

check-tenant:
	curl -i http://localhost:8082/health || echo "No /health endpoint implemented"

check-location:
	curl -i http://localhost:8083/health || echo "No /health endpoint implemented"

check-streaming:
	curl -i http://localhost:8084/health || echo "No /health endpoint implemented"
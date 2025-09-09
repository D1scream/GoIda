# Build the application
build:
	go build -o server -v ./cmd/server

# Run the application
run:
	go run ./cmd/server

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Database migrations
migrate-up:
	docker compose run --rm liquibase liquibase update

migrate-status:
	docker compose run --rm liquibase liquibase status

migrate-validate:
	docker compose run --rm liquibase liquibase validate

migrate-reset:
	docker compose run --rm liquibase liquibase dropAll
	docker compose run --rm liquibase liquibase update

# Docker Compose commands
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose up --build -d

docker-logs:
	docker compose logs -f

docker-restart:
	docker compose restart

docker-status:
	docker compose ps

docker-clean:
	docker compose down -v --remove-orphans

# Build the application
build:
	go build -o main -v .

# Run the application
run:
	go run main.go

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
	docker-compose run --rm liquibase liquibase update

migrate-status:
	docker-compose run --rm liquibase liquibase status

migrate-validate:
	docker-compose run --rm liquibase liquibase validate

migrate-reset:
	docker-compose run --rm liquibase liquibase dropAll
	docker-compose run --rm liquibase liquibase update

# Docker Compose commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose up --build -d

docker-logs:
	docker-compose logs -f

docker-restart:
	docker-compose restart

docker-status:
	docker-compose ps

docker-clean:
	docker-compose down -v --remove-orphans

# Development commands
dev-setup:
	cp env.example .env
	docker-compose up -d

dev-start:
	docker-compose up -d
	go run main.go

dev-stop:
	docker-compose down

# Quick start
start:
	make dev-setup
	sleep 10
	make migrate-up

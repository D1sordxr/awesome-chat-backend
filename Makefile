# Build the application containers
build:
	docker-compose build

# Start only the environment
up:
	docker-compose up -d

# Stop everything
down:
	docker-compose down

# Restart the app container
restart:
	docker-compose restart app

# Remove containers and volumes
clean:
	docker-compose down -v

# Run the Go application locally (outside the container)
run:
	go run ./cmd/api/main.go

# Run all tests
test:
	go test ./...

# Show logs from all services
logs:
	docker-compose logs -f --tail=100

lint:
	golangci-lint run ./...
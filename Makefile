generate_doc:
	swag init -g ./cmd/main.go -o ./internal/docs

build:
	go build -o bin/main main.go

start_docker_compose:
	docker compose -f ./scripts/dev-docker-compose.yaml up -d

run: start_docker_compose generate_doc
	go run cmd/main.go
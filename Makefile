MIGRATIONS_DIR=./migrations

migrate-create-%:
	goose -dir $(MIGRATIONS_DIR) create $(subst migrate-create-,,$@) sql
docker-up:
	docker compose up --build
start:
	go run  cmd/main.go

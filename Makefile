# === Variables ===
COMPOSE_FILE := compose.yaml

help:
	@echo "Available commands:"
	@echo " make docker-build                 Build the services"
	@echo " make docker-up                    Run the services"
	@echo " make docker-down                  Drop the services"
	@echo " make docker-migrator-up           Apply all up database migrations"
	@echo " make docker-migrator-down         Apply all down database migrations"
	@echo " make docker-start                 Build and run the services, and apply up migrations"

# ---------- Docker Commands ----------
docker-build:
	docker compose -f $(COMPOSE_FILE) build

docker-up:
	docker compose -f $(COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(COMPOSE_FILE) down -v

docker-migrator-up:
	docker compose run --rm auth-migrator -up

docker-migrator-down:
	docker compose run --rm auth-migrator -down 0

docker-start:
	make docker-build
	make docker-up
	make docker-migrator-up

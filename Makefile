IMAGE_NAME = go
COMPOSE_FILE = docker-compose.yml

# Targets
.PHONY: build up down clean logs

up:
	docker build . -t $(IMAGE_NAME)
	docker-compose -f $(COMPOSE_FILE) up -d

down:
	docker-compose -f $(COMPOSE_FILE) down
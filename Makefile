# Targets
.PHONY: build up down clean logs

up:
	docker-compose up 

build:
	docker build . -t test
	docker run --rm -it -p 3000:3000 --network=docker_mynetwork test

down:
	docker-compose down
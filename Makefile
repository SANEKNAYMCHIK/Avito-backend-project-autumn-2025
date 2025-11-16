.PHONY: build up down lint clean logs help

DOCKER_COMPOSE = docker-compose
GO = go
GOLANGCI_LINT = golangci-lint

help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  up              - Start the application with docker-compose"
	@echo "  down            - Stop the application"
	@echo "  lint            - Run linter"
	@echo "  clean           - Clean up containers and volumes"
	@echo "  logs            - Show logs of the service"
	

build:
	$(DOCKER_COMPOSE) build

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

lint:
	$(GOLANGCI_LINT) run

clean:
	$(DOCKER_COMPOSE) down -v
	$(GO) clean

logs:
	$(DOCKER_COMPOSE) logs -f app

setup: build up

.DEFAULT_GOAL := help
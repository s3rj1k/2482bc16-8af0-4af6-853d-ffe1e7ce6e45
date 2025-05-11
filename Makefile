.PHONY: all build up down restart logs help

all: up

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down -v

restart: down up

logs:
	docker-compose logs -f app

help:
	@echo "Available commands:"
	@echo "  make build     - Build Docker images"
	@echo "  make up        - Start all services"
	@echo "  make down      - Stop services and remove volumes"
	@echo "  make restart   - Restart all services"
	@echo "  make logs      - View logs for app service"
	@echo "  make help      - Show this help message"

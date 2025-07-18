.PHONY: up logs test build

up:
	docker-compose up -d --build

logs:
	docker-compose logs -f

test:
	@echo "--- Running backend tests ---"
	@docker build --target dev -t utils-backend-dev ./backend
	@docker run --rm -w /app utils-backend-dev ./gradlew test

build:
	docker-compose build

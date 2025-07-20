.PHONY: up logs test build

up:
	docker-compose up -d --build

logs:
	docker-compose logs -f

test:
	@echo "--- Running backend tests ---"
	@docker build -t utils-backend-test ./backend
	@docker run --rm -w /app utils-backend-test sh -c "chmod +x ./gradlew && ./gradlew test --rerun-tasks"

build:
	docker-compose build

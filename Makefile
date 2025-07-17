.PHONY: up logs

up:
	docker-compose up -d --build

logs:
	docker-compose logs -f
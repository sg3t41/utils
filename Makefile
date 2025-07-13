.PHONY: up logs

up:
	docker-compose up -d

logs:
	docker-compose logs -f
up:
	docker compose -f docker-compose.yml up -d --build
down:
	docker compose -f docker-compose.yml down
logs:
	docker compose -f docker-compose.yml logs -f
build:
	docker compose -f docker-compose.yml build

restart:
	docker compose -f docker-compose.yml down && docker compose -f docker-compose.yml up -d
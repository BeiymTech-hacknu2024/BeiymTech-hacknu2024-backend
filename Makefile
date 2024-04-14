DOCKER_COMPOSE_FILE ?= docker-compose.yml
docker-up:
docker-up:
	docker compose up --build 
	
docker-down:
docker-down:
	docker compose down

migrate-up: 
migrate-up:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate up

migrate-down: 
migrate-down:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate down 



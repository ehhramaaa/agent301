
build:
	docker build -t cats .

up:
	docker-compose up -d

down:
	docker-compose down

delete:
	docker rmi cats --force

.PHONY: build up down delete
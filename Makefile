
build:
	docker build -t agent301 .

up:
	docker-compose up -d

down:
	docker-compose down

delete:
	docker rmi agent301 --force

.PHONY: build up down delete
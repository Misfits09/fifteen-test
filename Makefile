update: services
	docker-compose up --build -d
	docker image prune -f

destroy: services
	docker-compose down
	docker image prune -f

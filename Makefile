deploy: services
	docker-compose up -d

deploy-update: services
	docker-compose up --build -d
	docker image prune -f

destroy: services
	docker-compose down
	docker image prune -f

lint: services
	golangci-lint run services/*

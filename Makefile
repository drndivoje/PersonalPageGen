build:
	go run . example && cp -r output/* deployment/data/

run: build
	docker compose -f deployment/docker-compose.yml up -d
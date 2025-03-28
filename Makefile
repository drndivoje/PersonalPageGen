build:
	go run . example && cp -r output/* deployment/data/
clean:
	rm -rf output/* deployment/data/*
run: clean build
	docker compose -f deployment/docker-compose.yml up -d
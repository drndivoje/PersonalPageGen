build:
	go run ./cmd/ppg inputFiles && cp -r output/* deployments/data/
clean:
	rm -rf output/* deployments/data/*
run: clean build
	docker compose -f deployments/docker-compose.yml up -d
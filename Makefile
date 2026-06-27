.PHONY: build test coverage clean run

build:
	go build -o ppg ./cmd/ppg

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out

clean:
	rm -f ppg coverage.out
	rm -rf output/*

run: clean build
	./ppg examples/simple
	cp -r output/* deployments/data/
	docker compose -f deployments/docker-compose.yml up -d

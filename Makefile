.PHONY: build test coverage clean run

build:
	go clean -cache -testcache
	go build -o ppg ./cmd/ppg

test:
	go test ./... -v

coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out

clean:
	rm -f ppg coverage.out
	rm -rf output/*
	rm -rf deployments/data/*

run: clean build
	./ppg examples/simple
	cp -r output/* deployments/data/
	docker compose -f deployments/docker-compose.yml up -d

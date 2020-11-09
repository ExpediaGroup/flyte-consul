test:
	go test ./...
build:
	go build .
docker-build:
	docker build -t flyte-consul .

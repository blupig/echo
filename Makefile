all:
	@echo "make [run | test | docker-build | docker-push]"

# Run locally
run:
	@go run echo.go

test:
	@go test -v -cover

# Build docker image for Linux
docker-build:
	@docker build -t blupig/echo --build-arg SOURCE_COMMIT=$(shell git rev-parse HEAD) .

docker-push: all
	@docker push blupig/echo

all:
	@echo "make [run | test | docker-build | docker-push]"

# Run locally
run:
	@go run echo.go

test:
	@go test -v -cover

# Build / push docker image for blupig/echo repo
docker-build:
	@docker build -t blupig/echo --build-arg SOURCE_COMMIT=$(shell git rev-parse HEAD) .

docker-push:
	@docker push blupig/echo

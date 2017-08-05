all:
	@CGO_ENABLED=0 GOOS=linux go build -o build/echo -a -ldflags '-extldflags "-static"' .
	@docker build -t yunzhu/echo .
	@docker push yunzhu/echo

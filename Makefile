all:
	@docker build -t yunzhu/echo --build-arg SOURCE_COMMIT=$(shell git rev-parse HEAD) .

push: all
	@docker push yunzhu/echo

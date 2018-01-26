all:
	@docker build -t yunzhu/echo .

push: all
	@docker push yunzhu/echo

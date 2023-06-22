PG = postgresql
MEM = memory

all:
	@echo "Specify storage type: make postgresql or make memory"

#test:
#	@echo "Running tests...."
#    @go test ./internal/config_parser
#   	@go test ./internal/tools
#   	@go test ./internal/database

$(PG):
	docker-compose --profile postgresql up --build

$(MEM):
	docker-compose --profile memory up --build

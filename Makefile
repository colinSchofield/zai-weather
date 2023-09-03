all: help

help:     				## Show this help message.
	@echo 'Usage: make [TARGET]'
	@echo 'Targets:'
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | awk -F ':.*?## ' 'NF==2 {printf "\033[36m  %-25s\033[0m %s\n", $$1, $$2}'

run: clean build		## Build and Run (in Docker) the Zai weather service.
	docker run -p8080:8080  weather

lint:					## Run lint checks.
	golangci-lint run ./...

test:	  				## Test and Code Coverage.
	go test ./... -cover

build:	  				## Build Docker image.
	docker build -t weather -f deployment/Dockerfile .

shell:    				## Shell into Docker image.
	docker run -it weather /bin/sh

clean:					## Remove any transient build artifacts.
	docker rmi weather -f
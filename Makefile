# PID file will store the server process id when it's running on development mode
PROJECT_NAME=$(shell basename "$(PWD)")
PID=/tmp/.$(PROJECT_NAME).pid

init: ## Run me to download some of this project dependencies for coding normalization
	pip3 install pre-commit
	pre-commit install --hook-type pre-commit
	pre-commit install --hook-type pre-push
	git clone https://github.com/lintingzhen/commitizen-go.git && cd commitizen-go && make && sudo ./commitizen-go install && cd .. && rm -rf commitizen-go

commit: ## Commit changes using commitizen
	@git cz

install: ## Rebuild the go.mod and go.sum files (removing unused ones)
	@go mod tidy
.PHONY: install

generate: ## Run go generate in the project root
	@go generate ./...
.PHONY: generate

migrate-up: ## Run all migrations not yet applied to the database (the migrations are located in the ./migrations folder). Run `make migrate-up filename=some_file.yml` to run the migration only for this file
	@go build -v -o build/migrate cmd/migrate/main.go
	./build/migrate $(filename)

build: ## Build project binary
	@go build -v -o build/bin cmd/using-gin/main.go
.PHONY: build

run: build ## Run the application
	./build/bin 2>&1 & echo $$! >> $(PID) &
.PHONY: run

stop: ## Stop current running server
	@-touch $(PID)
	@-kill `cat $(PID)` 2> /dev/null || true
	@-rm $(PID)
.PHONY: stop

help: ## Display help screen
	@echo "Usage:"
	@echo "	make [COMMAND]"
	@echo "	make help \n"
	@echo "Commands: \n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

.DEFAULT_GOAL := help

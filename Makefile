# PID file will store the server process id when it's running on development mode
PROJECT_NAME=$(shell basename "$(PWD)")
PID=/tmp/.$(PROJECT_NAME).pid

init: ## Run me to download some of this project dependencies for coding normalization
	pip3 install pre-commit
	pre-commit install --hook-type pre-commit
	pre-commit install --hook-type pre-push
	git clone https://github.com/lintingzhen/commitizen-go.git && cd commitizen-go && make && sudo ./commitizen-go install && cd .. && rm -rf commitizen-go
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/vektra/mockery/v2/.../
.PHONY: init

commit: ## Commit changes using commitizen
	@git cz
.PHONY: commit

install: ## Rebuild the go.mod and go.sum files (removing unused ones)
	@go mod tidy
.PHONY: install

generate: ## Run go generate in the project root
	@go generate ./...
	@mockery --exported --dir ./pkg/domain/species/ --name repositorier --output ./pkg/domain/species --structname repositoryMock --filename repository_mock.go --inpackage
.PHONY: generate

generate-docs: ## Generate project documentation
	@swag init -g cmd/using-gin/main.go -o pkg/server/docs
.PHONY: generate-docs

migrate-up: ## Run all migrations not yet applied to the database (the migrations are located in the ./migrations folder). Run `make migrate-up filename=some_file.yml` to run the migration only for this file
	go build -v -o build/migrate cmd/migrate/main.go
ifndef filename
	./build/migrate up
else
	./build/migrate up -f $(filename)
endif
.PHONY: migrate-up

migrate-down: ## Undo all migrations already applied to the database (the migrations are located in the ./migrations folder). Run `make migrate-down filename=some_file.yml` to undo the migration only for this file
	go build -v -o build/migrate cmd/migrate/main.go
ifndef filename
	./build/migrate down
else
	./build/migrate down -f $(filename)
endif
.PHONY: migrate-down

migrate-create: ## Create an empty migration file at ./migrations. Run `make migrate-create name=migration-name` to customize the migration file
	go build -v -o build/migrate cmd/migrate/main.go
ifndef name
	./build/migrate create -n empty-migration
else
	./build/migrate create -n $(name)
endif
.PHONY: migrate-down

migrate-create-table: ## Create a create table migration file at ./migrations. Run `make migrate-create tablename=my-table` to indicate the table name
	go build -v -o build/migrate cmd/migrate/main.go
ifndef tablename
	./build/migrate create table -t table
else
	./build/migrate create table -t $(tablename)
endif
.PHONY: migrate-down

lint: ## Run lint
	golangci-lint run
.PHONY: lint

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

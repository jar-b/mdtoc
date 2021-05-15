.DEFAULT_GOAL:=help

.PHONY: build
build: clean ## Build binaries
	@go build -o mdtoc cmd/mdtoc/main.go

.PHONY: clean
clean: ## Clean up binaries
	@rm -f mdtoc

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'


.PHONY: test
test: ## Run unit tests
	@go test -v ./...

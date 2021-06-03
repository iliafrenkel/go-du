PROJECT_NAME := "go-du"
USER_NAME := "iliafrenkel"
PKG := "github.com/$(USER_NAME)/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

all: build

dep: ## Get the dependencies
	@go mod download

lint: ## Lint all Golang files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run all the unit tests
	@go test -short ${PKG_LIST}

test-coverage: ## Run all the unit tests with coverage report
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out >> coverage.txt

build: dep ## Build the binary
	@go build -o build/$(PROJECT_NAME) $(PKG)/app

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)/build

help: ## Print this help message
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

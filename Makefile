NAME="gopherlol"
PORT=8080

help: ## show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: run
run: ## run the application
	go run .

.PHONY: build
build: ## build the application
	go build -o bin/$(NAME) .

.PHONY: clean
clean: ## clean build artifacts
	rm -rf bin/

.PHONY: install-deps
install-deps: ## install dependencies
	go mod download
	go mod tidy

##@ Testing
.PHONY: test
test: ## run tests
	go test ./...

.PHONY: test-verbose
test-verbose: ## run tests with verbose output
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## run tests with coverage
	go test -cover ./...

.PHONY: test-coverage-html
test-coverage-html: ## run tests with coverage and generate HTML report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

##@ Code Quality
.PHONY: fmt
fmt: ## format code
	go fmt ./...

.PHONY: vet
vet: ## run go vet
	go vet ./...

.PHONY: lint
lint: ## run golangci-lint (requires installation)
	golangci-lint run

.PHONY: check
check: fmt vet test ## run all checks (format, vet, test)

##@ Utilities
.PHONY: dev
dev: ## run in development mode with auto-reload (requires air)
	air

.PHONY: open
open: ## open the application in browser
	open "http://localhost:$(PORT)"

.PHONY: status
status: ## check if application is running
	@curl -s http://localhost:$(PORT)/\?q\=help > /dev/null && echo "Application is running" || echo "Application is not running"

.PHONY: docker-build
docker-build: ## build docker image
	docker build -t $(NAME) .

.PHONY: docker-run
docker-run: ## run docker container
	docker run -p $(PORT):$(PORT) --name $(NAME) $(NAME)

##@ Release
.PHONY: release-build
release-build: ## build for multiple platforms
	GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/$(NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/$(NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/$(NAME)-windows-amd64.exe .
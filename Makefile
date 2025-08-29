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

.PHONY: tidy
tidy: ## tidy up go modules
	go mod tidy

##@ Testing
.PHONY: test
test: ## run tests
	go test ./...

.PHONY: test-verbose
test-verbose: ## run tests with verbose output and coverage
	go test -v -cover ./...

##@ Code Quality
.PHONY: check
check: ## run all checks (format, vet, test with coverage)
	go fmt ./...
	go vet ./...
	go test -cover ./...

##@ Utilities
.PHONY: open
open: ## open the application in browser
	open "http://localhost:$(PORT)"

.PHONY: status
status: ## check if application is running
	@curl -s http://localhost:$(PORT)/\?q\=help > /dev/null && echo "Application is running" || echo "Application is not running"

##@ Release
.PHONY: release-build
release-build: ## build for multiple platforms
	GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/$(NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/$(NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/$(NAME)-windows-amd64.exe .
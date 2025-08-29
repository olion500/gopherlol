NAME="gopherlol"
PORT=8080

help: ## show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: usage
usage: ## show browser setup instructions
	@echo ""
	@echo "\033[1m🔍 gopherlol Browser Setup Instructions\033[0m"
	@echo ""
	@echo "\033[36m1. Start the server:\033[0m"
	@echo "   make run"
	@echo ""
	@echo "\033[36m2. Add search engine to your browser:\033[0m"
	@echo ""
	@echo "\033[33m   Chrome:\033[0m"
	@echo "   • Go to Settings → Search engine → Manage search engines"
	@echo "   • Click 'Add' next to 'Other search engines'"
	@echo "   • Search engine: gopherlol"
	@echo "   • Keyword: gl (or any shortcut you prefer)"
	@echo "   • URL: http://localhost:$(PORT)/?q=%s"
	@echo "   • Click 'Add' and optionally 'Make default'"
	@echo ""
	@echo "\033[33m   Firefox:\033[0m"
	@echo "   • Go to Settings → Search → Search Shortcuts"
	@echo "   • Click 'Add' and enter the same details as Chrome"
	@echo ""
	@echo "\033[33m   Safari:\033[0m"
	@echo "   • Go to Safari → Settings → Search"
	@echo "   • Click 'Manage Search Engines' and add custom engine"
	@echo ""
	@echo "\033[36m3. Test it:\033[0m"
	@echo "   • Type in address bar: gl help"
	@echo "   • Try: gl g hello world"
	@echo "   • Try: gl gh pr typescript"
	@echo ""
	@echo "\033[36m4. View all commands:\033[0m"
	@echo "   • Visit: http://localhost:$(PORT)/?q=help"
	@echo "   • Or type in browser: gl help"
	@echo ""
	@echo "\033[32m✨ Pro tip: Set 'gl' as your keyword for quick access!\033[0m"
	@echo ""

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

##@ Release
.PHONY: release-build
release-build: ## build for multiple platforms
	GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/$(NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/$(NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/$(NAME)-windows-amd64.exe .
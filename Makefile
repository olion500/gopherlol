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
	@echo "\033[33m   Chrome (and Chromium-based browsers):\033[0m"
	@echo "   • Open: chrome://settings/searchEngines"
	@echo "   • Click 'Add' next to 'Other search engines'"
	@echo "   • Search engine: gopherlol"
	@echo "   • Keyword: gl (or any shortcut you prefer)"
	@echo "   • URL: http://localhost:$(PORT)/?q=%s"
	@echo "   • Click 'Add' and optionally 'Make default'"
	@echo ""
	@echo "\033[33m   Firefox:\033[0m"
	@echo "   • Open: about:preferences#search"
	@echo "   • Click 'Add' under 'Search Shortcuts'"
	@echo "   • Enter the same details as Chrome above"
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

##@ Dependencies
.PHONY: install-deps
install-deps: ## install Rust and Tauri dependencies
	asdf plugin add rust https://github.com/asdf-community/asdf-rust.git
	asdf plugin add golang https://github.com/asdf-community/asdf-golang.git
	asdf install

##@ Development
.PHONY: run
run: ## run the UI application (default)
	cd ui && cargo tauri dev

.PHONY: run-cli
run-cli: ## run the CLI server only
	go run .

.PHONY: ui-dev
ui-dev: ## run UI in development mode
	cd ui && cargo tauri dev

.PHONY: ui-build
ui-build: ## build UI application
	cd ui && cargo tauri build

.PHONY: build
build: ## build the application
	go build -o $(NAME) .

.PHONY: clean
clean: ## clean build artifacts and logs
	rm -f $(NAME) usage.log
	rm -rf dist/

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

##@ Analytics
.PHONY: analytics
analytics: ## show command usage analytics (add -overall for all-time stats)
	@go run cmd/analytics/main.go $(ARGS)



##@ Release
.PHONY: release-build
release-build: ## build for multiple platforms
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/$(NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o dist/$(NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o dist/$(NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/$(NAME)-windows-amd64.exe .
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

gopherlol is a smart bookmarking tool inspired by Facebook's bunnylol/bunny1. It's a modernized Go-based HTTP server that acts as a custom search engine, allowing users to define custom commands that redirect to specific URLs. When an unrecognized command is entered, it falls back to Google search.

## Development Commands

- **Run the application**: `make run` or `go run .`
- **Build**: `make build` or `go build`
- **Test**: `make test` or `go test ./...`
- **Test with coverage**: `make test-coverage`
- **Format code**: `make fmt`
- **Run all checks**: `make check` (format, vet, test)
- **Clean build artifacts**: `make clean`

The application runs on `http://localhost:8080` and expects search queries at `/?q=%s`.

## Modern Architecture

The codebase has been modernized with several key improvements:

### Project Setup
- **Go 1.23.0**: Latest Go version managed with asdf via `.tool-versions`
- **Comprehensive Makefile**: Development, testing, and deployment commands
- **Full Test Coverage**: Unit tests for all components with 100% coverage

### Core HTTP Server (`main.go`)
- Single HTTP handler that processes all incoming requests
- JSON-based configuration system instead of reflection-based discovery
- Parses query parameters and handles command/subcommand routing
- Special handling for `help`/`list` commands that generate rich HTML documentation
- Falls back to Google search for unrecognized commands

### JSON Configuration System (`commands.json`)
- All commands defined in a single JSON configuration file
- Easy to edit and maintain without code changes
- Supports rich metadata: descriptions, aliases, subcommands
- Template-based URL generation with `{{.Query}}` placeholders

### Command Registry (`internal/config/`)
- `CommandConfig`: JSON structure definition
- `CommandRegistry`: In-memory lookup system for commands and aliases
- `ExecuteURL()`: Template processing for dynamic URLs
- Support for multiple aliases per command
- Hierarchical subcommand support

## Key Features

### Multiple Aliases
Commands can have multiple aliases:
```json
{
  "name": "google",
  "aliases": ["g", "search"],
  "url": "https://www.google.com/#q={{.Query}}"
}
```

### Subcommands
Commands can have subcommands with their own aliases:
```json
{
  "name": "github",
  "aliases": ["gh"],
  "subcommands": [
    {
      "name": "pr",
      "aliases": ["pull", "pullrequest"],
      "url": "https://github.com/search?type=pullrequests&q={{.Query}}"
    }
  ]
}
```

### Built-in Commands
- **google/g/search**: Google search
- **stackoverflow/so/stack**: Stack Overflow search  
- **github/gh**: GitHub operations with subcommands (pr, issues, repo, user)
- **datadog/dd**: Datadog operations with subcommands (logs, metrics, dashboard)
- **gmail/mail/email**: Gmail search
- **vscode/code/vs**: VS Code marketplace with subcommands (extensions, themes)
- **youtube/yt**: YouTube search
- **twitter/tw/x**: Twitter/X search

## Usage Examples

- `g hello world` → Google search for "hello world"
- `gh pr typescript` → Search GitHub pull requests for "typescript"
- `dd logs error` → Search Datadog logs for "error"
- `so golang testing` → Search Stack Overflow for "golang testing"

## Adding New Commands

To add new commands:
1. Edit `commands.json` to add your command definition
2. Include name, aliases, description, URL template, and any subcommands
3. Use `{{.Query}}` in URLs where the search query should be inserted
4. Restart the application - no code changes needed!

Example:
```json
{
  "name": "myservice",
  "aliases": ["ms"],
  "description": "Search my service",
  "url": "https://myservice.com/search?q={{.Query}}",
  "requiresQuery": true
}
```

## Implementation Details

- Command names and aliases are case-insensitive
- URL templates use Go's `text/template` package
- Proper URL encoding is handled automatically
- Rich help page shows all commands, aliases, and subcommands
- Unrecognized commands fall back to Google search
- Comprehensive test suite ensures reliability

## Reference Documentation

### Tauri Framework
For any Tauri-related development tasks, refer to the comprehensive documentation in `docs/tauri-llms.txt`. This file contains the complete Tauri framework documentation index covering:

- **Start**: Project creation, prerequisites, frontend configuration
- **Concepts**: Architecture, process model, inter-process communication  
- **Security**: Capabilities, permissions, CSP, runtime authority
- **Development**: API calls, configuration, debugging, state management, plugin development
- **Distribution**: App store publishing, installers, code signing
- **Learning**: Tutorials, system tray, window customization, splashscreen
- **Plugins**: Extensibility and community plugins

When working with Tauri features, always consult this documentation first to understand proper implementation patterns, security considerations, and best practices.
- to memorize
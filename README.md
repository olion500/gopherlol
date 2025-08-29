# gopherlol 🔍

> A modernized [bunnylol](https://www.quora.com/What-is-Facebooks-bunnylol) / [bunny1](http://www.bunny1.org) -like smart bookmarking tool, written in Go

Transform your browser's address bar into a powerful command center for quick navigation to your favorite websites and services!

## ✨ Features

- 🚀 **JSON Configuration**: Easy-to-edit commands without code changes
- 🏷️ **Multiple Aliases**: `g`, `google`, `search` all work for Google
- 🌳 **Subcommands**: `gh pr` for GitHub pull requests, `dd logs` for Datadog logs
- 🎯 **Smart Fallback**: Unknown commands automatically search Google
- 📚 **Rich Help**: Type `help` to see all commands, aliases, and descriptions
- ⚡ **Lightning Fast**: Instant redirects to your destination

## 🚀 Quick Start

```bash
# Run the application
make run
# or
go run .

# Then add `http://localhost:8080/?q=%s` as a search engine to your browser
```

### Using asdf (recommended)
```bash
asdf install golang 1.23.0  # Install Go 1.23.0
make run                    # Run with latest Go
```

## 📖 Usage Examples

| Command | Description | Where it goes |
|---------|-------------|---------------|
| `g hello world` | Google search | https://google.com/#q=hello+world |
| `gh pr typescript` | GitHub PRs | https://github.com/search?type=pullrequests&q=typescript |
| `so golang testing` | Stack Overflow | https://stackoverflow.com/search?q=golang+testing |
| `dd logs error` | Datadog logs | https://app.datadoghq.com/logs?query=error |
| `yt funny cats` | YouTube search | https://youtube.com/results?search_query=funny+cats |

## 🛠️ Development

```bash
# All available commands
make help

# Common operations
make run           # Start the server
make test          # Run tests
make build         # Build binary
make check         # Run format, vet, and tests
```

## 🎛️ Built-in Commands

### Core Services
- **google** (`g`, `search`) - Google search
- **stackoverflow** (`so`, `stack`) - Stack Overflow search
- **youtube** (`yt`) - YouTube search
- **twitter** (`tw`, `x`) - Twitter/X search

### Developer Tools
- **github** (`gh`) - GitHub with subcommands:
  - `pr` (`pull`) - Pull requests
  - `issues` (`issue`) - Issues
  - `repo` (`repository`) - Repositories
  - `user` (`users`) - Users
- **vscode** (`code`, `vs`) - VS Code marketplace with subcommands:
  - `extensions` (`ext`) - Extensions
  - `themes` (`theme`) - Themes

### Enterprise Services
- **datadog** (`dd`) - Datadog with subcommands:
  - `logs` (`log`) - Log search
  - `metrics` (`metric`) - Metrics explorer
  - `dashboard` (`dash`) - Dashboard search
- **gmail** (`mail`, `email`) - Gmail search

## 🔧 Adding Custom Commands

Simply edit `commands.json`:

```json
{
  "name": "myservice",
  "aliases": ["ms", "service"],
  "description": "Search my internal service",
  "url": "https://myservice.company.com/search?q={{.Query}}",
  "requiresQuery": true,
  "subcommands": [
    {
      "name": "docs",
      "aliases": ["documentation"],
      "description": "Search documentation",
      "url": "https://myservice.company.com/docs?q={{.Query}}"
    }
  ]
}
```

Then restart the server - no code changes needed!

## 🌐 Browser Setup

### Chrome
1. Go to Settings → Search engine → Manage search engines
2. Add new search engine:
   - **Search engine**: gopherlol
   - **Keyword**: gopherlol (or any shortcut you prefer)
   - **URL**: `http://localhost:8080/?q=%s`

### Other Browsers
- [Instructions for all major browsers](https://www.howtogeek.com/114176/how-to-easily-create-search-plugins-add-any-search-engine-to-your-browser/)

## 💡 Why Use This?

Just like [Facebook's internal bunnylol](http://www.ccheever.com/blog/?p=74), gopherlol transforms your browser into a productivity powerhouse:

- **Speed**: No more bookmarks hunting or typing full URLs
- **Consistency**: Same commands work across all your devices  
- **Extensibility**: Add your company's internal tools and services
- **Muscle Memory**: Short, memorable commands become second nature

## 🏗️ Architecture

- **Modern Go**: Built with Go 1.23.0, latest language features
- **JSON Config**: All commands defined in `commands.json`
- **Template URLs**: Dynamic URL generation with `{{.Query}}` placeholders
- **Comprehensive Tests**: Full test coverage for reliability
- **Clean Code**: Well-structured, maintainable codebase

## 🤝 Contributing

1. Fork the repository
2. Make your changes
3. Add tests for new functionality
4. Run `make check` to ensure quality
5. Submit a pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

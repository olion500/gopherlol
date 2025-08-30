# Contributing to gopherlol 🤝

Thank you for your interest in contributing to gopherlol! We welcome contributions from everyone.

## 🚀 Quick Start

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/gopherlol.git
   cd gopherlol
   ```
3. **Install dependencies** (Go 1.23.0+):
   ```bash
   # Using asdf (recommended)
   asdf install golang 1.23.0
   
   # Or install Go from https://golang.org/dl/
   ```

## 🔧 Development Workflow

### Setting Up
```bash
# Copy the sample configuration
cp commands.json.sample commands.json

# Run the development server
make run

# In another terminal, run tests
make test
```

### Making Changes
1. **Create a branch** for your feature:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** and ensure they follow the project standards:
   ```bash
   # Format code
   make fmt
   
   # Run all checks (format, vet, test)
   make check
   ```

3. **Add tests** for any new functionality
4. **Update documentation** if needed

### Code Standards
- **Go Style**: Follow standard Go formatting (`gofmt`)
- **Tests**: Maintain 100% test coverage
- **Documentation**: Update README.md for user-facing changes
- **Commit Messages**: Use clear, descriptive commit messages

## 🧪 Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run all quality checks
make check
```

## 📝 Types of Contributions

### 🐛 Bug Reports
- Use GitHub Issues with the "bug" label
- Include steps to reproduce
- Provide Go version and OS information

### 💡 Feature Requests
- Use GitHub Issues with the "enhancement" label
- Describe the use case and expected behavior
- Consider backward compatibility

### 🔧 Code Contributions
- **New Commands**: Add to `commands.json.sample`
- **Core Features**: HTTP server, command parsing, URL generation
- **Developer Experience**: Makefile, tooling, documentation
- **Tests**: Always include comprehensive tests

## 📚 Project Structure

```
gopherlol/
├── main.go              # HTTP server & request routing
├── internal/config/     # Command registry & JSON parsing
│   ├── config.go        # Configuration types
│   └── registry.go      # Command lookup & execution
├── commands.json        # Runtime configuration
├── commands.json.sample # Template for new users
├── Makefile            # Build & development commands
└── *_test.go           # Test files
```

## 🎯 Areas That Need Help

- **Browser Integration**: Better setup instructions for different browsers
- **Command Templates**: More built-in commands for popular services  
- **Performance**: URL generation and redirect optimization
- **Documentation**: Examples, tutorials, advanced usage guides
- **Testing**: Edge cases, error conditions, integration tests

## ❓ Questions?

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and community chat
- **Email**: Open an issue first, we prefer public discussions!

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for making gopherlol better!** 🎉
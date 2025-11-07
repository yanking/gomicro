# GoMicro

A Go project with multiple MySQL instance support.

## Features

- Multiple MySQL instance support
- Configuration management
- Linting with golangci-lint
- Git pre-commit hooks

## Installation

```bash
go mod tidy
```

## Configuration

Configuration files are located in the `configs/` directory. The project supports multiple MySQL instances as defined in the configuration.

## Linting

This project uses golangci-lint for code quality checks.

### Running linting manually

```bash
# Using the provided script
./scripts/lint.sh

# Or directly
golangci-lint run ./...
```

### Git pre-commit hooks

The project includes a pre-commit hook that automatically runs linting before each commit. To install the hooks:

```bash
./scripts/install-hooks.sh
```

The hook will run automatically before each commit and will prevent the commit if linting fails. To bypass the check (not recommended), use:

```bash
git commit --no-verify
```

## Usage

Examples of how to use the various components can be found in the `examples/` directory.

## Documentation

- [MySQL Multi-instance Usage](pkg/client/database/README.md)
- [Configuration Management](pkg/conf/README.md)
- [Linting Guide](docs/linting.md)
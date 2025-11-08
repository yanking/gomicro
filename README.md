# GoMicro

A Go project with multiple MySQL instance support.

## Features

- Multiple MySQL instance support
- Multiple Redis instance support
- Configuration management
- Linting with golangci-lint
- Git pre-commit hooks
- HTTP transport layer based on Gin
- API Documentation with Swagger

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

## Transport Layer

The project includes an HTTP transport layer based on the Gin framework:

- [HTTP Transport Documentation](pkg/transport/http/README.md)

## API Documentation

This project uses Swagger for API documentation. The documentation is located in `docs/swagger/swagger.json`.

### Initializing Swagger Documentation

```bash
make swagger-init
```

This command creates the swagger documentation directory and ensures it exists.

### Generating Swagger Documentation

```bash
make swagger-generate
```

This command automatically generates Swagger documentation from source code annotations. The documentation is generated from the example files in the `examples/` directory.

### Validating Swagger Documentation

```bash
make swagger-validate
```

### Serving Swagger UI

```bash
make swagger-serve
```

After running this command, open your browser and navigate to `http://localhost:63150/docs` to view the interactive API documentation.

## Usage

Examples of how to use the various components can be found in the `examples/` directory.

## Documentation

- [MySQL Multi-instance Usage](pkg/client/database/README.md)
- [Redis Multi-instance Usage](pkg/client/database/README.md)
- [HTTP Transport Usage](pkg/transport/http/README.md)
- [Configuration Management](pkg/conf/README.md)
- [Linting Guide](docs/linting.md)
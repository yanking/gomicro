# Internal Packages

This directory contains the internal packages of the application. These packages are not intended to be used by external applications.

## Structure

- `config/` - Configuration management
- `handler/` - HTTP handlers
- `service/` - Business logic services
- `repository/` - Data access layer
- `model/` - Domain models
- `server/` - HTTP server setup
- `logger/` - Logging utilities
- `middleware/` - HTTP middleware
- `utils/` - Utility functions

## Package Organization

Each package should follow the single responsibility principle and have a clear purpose. Packages should not have circular dependencies.

## Dependency Injection

Dependencies between packages should be injected through interfaces to promote loose coupling and testability.
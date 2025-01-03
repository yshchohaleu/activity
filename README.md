# Activity Project

## Overview

The Activity project is a Go-based application designed to manage and track various activities. It leverages GraphQL for API interactions and integrates with Firebase for backend services.

## Project Structure

- **cmd/**: Contains the main application entry points.
- **internal/**: Houses the core application logic and internal packages.
- **graph/**: Contains GraphQL schema and resolver implementations.

## Key Technologies

- **Go**: The primary programming language used for this project.
- **GraphQL**: Utilized for API interactions, implemented using `gqlgen`.
- **Firebase**: Integrated for backend services.
- **GORM**: Used for ORM with support for SQLite and PostgreSQL.

## Dependencies

- `firebase.google.com/go/v4`
- `github.com/99designs/gqlgen`
- `github.com/glebarez/sqlite`
- `github.com/google/uuid`
- `github.com/stretchr/testify`
- `gorm.io/driver/postgres`
- `gorm.io/gorm`

## Setup

1. Ensure you have Go 1.22 or later installed.
2. Clone the repository.
3. Run `go mod tidy` to install dependencies.
4. Use `docker-compose` to set up the environment.

## Usage

To start the application, navigate to the `cmd` directory and run:

```bash
go run main.go
```

## License

This project is licensed under the MIT License. 
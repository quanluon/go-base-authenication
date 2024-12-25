# go-base-authenication

## Setup

1. Install `go`

```sh
brew install go
```

2. Install `sqlc`

```sh
go install github.com/kyleconroy/sqlc@latest
```

3. Copy .env.example to .env and fill in the values

### Installation Requirements

- Go version 1.16 or higher
- Access to a Postgres database

## DB Setup

1. Run `make create-migration name=<name>` to create a new migration
2. Run `make migrate` to apply the migrations
3. Run `make rollback` to rollback the migrations
4. Run `make sqlc` to generate the database models

## Run

1. Run `make run` to start the server

## Features

1. Authentication
2. Authorization
3. User Management
4. Role-Based Access Control (RBAC)
5. JWT
6. Refresh Token
7. Password Hashing
8. Error Handling
9. Logging
10. Middleware

### Usage Examples

- To log in, send a POST request to `/auth/login` with the required body.
- To register a new user, send a POST request to `/auth/register` with the necessary details.

## Structure

Source code structure is following Golang best practices

- cmd: The main entry point for the server
- internal: The internal packages
- pkg: The external packages
- middleware: The middleware

Internal packages are following Clean Architecture

- controllers: The HTTP handlers
- services: The business logic
- repositories: The database operations
- models: The database models
- db: The database migrations
- routes: The routes for the API
- middlewares: Middlewares of server

## API Documentation

Auth apis

- POST /auth/login
  Body: { email: string, password: string }
  Response: { accessToken: string, refreshToken: string, user: UserResponse, expiresIn: number, refreshExpiresIn: number }

- POST /auth/refresh
  Body: { refreshToken: string }
  Response: { accessToken: string, refreshToken: string, user: UserResponse, expiresIn: number, refreshExpiresIn: number }

- POST /auth/register
  Body: { email: string, password: string, name: string }
  Response: UserResponse

User apis

- GET /users
  Response: []UserResponse

- GET /users/:id
  Response: UserResponse

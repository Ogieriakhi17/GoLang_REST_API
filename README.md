# Go To-Do REST API with Authentication
A production-style REST API built in Go using Gin, PostgreSQL, and JWT authentication. This project demonstrates backend engineering fundamentals including secure authentication, middleware, database integration, and clean architecture.

# Overview
This API allows users to register, authenticate, and manage their personal to-do items securely. Each user can create, view, update, and delete their own tasks, with all protected routes requiring JWT authentication.

The project is designed using real-world backend practices such as dependency injection, middleware-based authentication, structured routing, and database connection pooling.

# Features
User registration and login

Secure JWT authentication

Protected routes using middleware

Full CRUD operations for to-do items

PostgreSQL database integration using pgxpool

Clean modular architecture

Environment-based configuration

Proper error handling and validation

Dependency injection for testability and scalability

# TechStack
## Language
Go

## Framework
Gin (HTTP web framework)

Database

PostgreSQL

pgx v5 connection pool

## Authentication
JWT (JSON Web Tokens)

## Architecture
REST API

Middleware-based authentication

Modular project structure

# Project Structure
```
todos_api/
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── database/
│   │   └── database.go
│   │
│   ├── handlers/
│   │   ├── auth.go
│   │   └── todos.go
│   │
│   ├── middleware/
│   │   └── auth.go
│   │
│   ├── models/
│       ├── user.go
│       └── todo.go
│   │
│   └── repository/
│   │   ├── user_repository.go
│       └── todo_repository.go
│
├── .env
├── go.mod
└── README.md

```
# Project Installation and Setup
## Clone repository
```
git clone https://github.com/YOUR_USERNAME/todos_api.git

cd todos_api
```

## .env File Set Up
```
PORT=8080

DATABASE_URL=postgres://username:password@localhost:5432/todos_db

JWT_SECRET=your_super_secret_key
```

## Create the PostgreSQL database

## Install Dependencies
```
go mod tidy
```
## Run the server
go run cmd/main.go


# License

MIT License







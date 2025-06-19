# GoCleanArch

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/celpung/gocleanarch)](https://goreportcard.com/report/github.com/celpung/gocleanarch)
[![GoDoc](https://godoc.org/github.com/celpung/gocleanarch?status.svg)](https://godoc.org/github.com/celpung/gocleanarch)

> **Version:** `v2.2.0`

## 📚 Introduction

**GoCleanArch** is a reference implementation of the **Clean Architecture** pattern in a Go (Golang) application. The project is structured to emphasize **separation of concerns**, **testability**, and **scalability**. By organizing the application into distinct layers, it promotes maintainable and readable codebases—especially useful for medium to large-scale systems.

## 📂 Project Structure

```bash
gocleanarch
├── cmd
│   ├── gin
│   │   └── main.go
│   └── http
│       └── main.go
├── configs
│   ├── database
│   │   ├── mysql
│   │   │   └── mysql_connection.go
│   │   └── sqlite
│   │       └── sqlite_conntection.go
│   ├── environment
│   │   └── environment.go
│   └── role
│       └── user_role.go
├── delivery
│   ├── dto
│   │   └── user_dto.go
│   ├── gin
│   │   └── user_delivery
│   │       ├── implementation
│   │       │   └── user_delivery_implementation.go
│   │       ├── middlewares
│   │       │   └── auth_middleware.go
│   │       ├── router
│   │       │   └── user_router.go
│   │       └── user_delivery_interface.go
│   └── http
│       └── user_delivery
│           ├── implementation
│           │   └── user_delivery_implementation.go
│           ├── middleware
│           │   └── auth_middleware.go
│           ├── router
│           │   └── user_router.go
│           └── user_delivery_interface.go
├── Dockerfile
├── domain
│   ├── slider
│   │   └── entity
│   │       └── slider_entity.go
│   └── user
│       ├── entity
│       │   └── user_entity.go
│       ├── repository
│       │   ├── implementation
│       │   │   ├── test
│       │   │   │   └── user_repository_implementation_test.go
│       │   │   └── user_repository_implementation.go
│       │   ├── model
│       │   │   └── user_model.go
│       │   └── user_repository_interface.go
│       └── usecase
│           ├── implementation
│           │   └── user_usecase_implementation.go
│           └── user_usecase_interface.go
├── go.mod
├── go.sum
├── public
│   ├── images
│   │   └── sample.jpg
│   └── index.html
├── readme.md
├── services
│   ├── bcrypt.go
│   └── jwt.go
└── utils
    ├── file_checker.go
    ├── method_handler.go
    └── request_method_check.go
```

## 🚀 Getting Started

These instructions will help you set up and run the project locally for development and testing purposes.

### ✅ Prerequisites

Make sure you have the following installed on your system:

- [Go](https://golang.org/dl/) version `1.12.x` or higher

### 📦 Installation

To get started, clone the repository and install the required dependencies:

```bash
# Clone the repository
git clone https://github.com/celpung/gocleanarch.git

# Navigate into the project directory
cd gocleanarch

# Download and tidy up the dependencies
go mod tidy

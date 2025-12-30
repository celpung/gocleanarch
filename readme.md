# GoCleanArch

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/celpung/gocleanarch)](https://goreportcard.com/report/github.com/celpung/gocleanarch)
[![GoDoc](https://godoc.org/github.com/celpung/gocleanarch?status.svg)](https://godoc.org/github.com/celpung/gocleanarch)

> **Version:** `v3.0.0`

## 📚 Introduction

**GoCleanArch** is a reference implementation of the **Clean Architecture** pattern in a Go (Golang) application. The project is structured to emphasize **separation of concerns**, **testability**, and **scalability**. By organizing the application into distinct layers, it promotes maintainable and readable codebases—especially useful for medium to large-scale systems.

## 🚀 Getting Started

These instructions will help you set up and run the project locally for development and testing purposes.

### ✅ Prerequisites

Make sure you have the following installed on your system:

- [Go](https://golang.org/dl/) version `1.24.x` or higher

### 📦 Installation

To get started, clone the repository and install the required dependencies:

```bash
# Clone the repository
git clone https://github.com/celpung/gocleanarch.git

# Navigate into the project directory
cd gocleanarch

# Download and tidy up the dependencies
go mod tidy
```

### 🗄️ Database migrations & seed

- Run migrations explicitly: `go run ./cmd/migration`
- The API also runs migrations on boot before serving traffic.
- Optional sample data: `go run ./cmd/seed` (includes migrations).
- Supported dialects: set `DB_DIALECT` to `mysql` (default) or `sqlite`. For SQLite, `DB_NAME` is used as the database file (e.g., `app.db`).

### ▶️ Run the API

```bash
go run ./cmd/api
```

### 🧪 Testing

- Unit tests use in-memory adapters for speed and isolation.

```bash
go test ./...
```

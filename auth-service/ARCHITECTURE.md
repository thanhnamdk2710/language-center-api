# Auth Service - Clean Architecture

This document describes the Clean Architecture structure of the auth-service.

## Overview

The service follows Clean Architecture principles with clear separation of concerns across layers:

```
auth-service/
├── cmd/                        # Application entry points
│   └── api/                    # HTTP API server
├── internal/
│   ├── domain/                 # Enterprise Business Rules (Layer 1)
│   │   ├── entity/             # Domain entities
│   │   ├── valueobject/        # Value objects and domain errors
│   │   ├── repository/         # Repository interfaces (ports)
│   │   └── service/            # Domain service interfaces (ports)
│   │       ├── token/
│   │       ├── password/
│   │       ├── email/
│   │       ├── otp/
│   │       └── session/
│   │
│   ├── usecase/                # Application Business Rules (Layer 2)
│   │   ├── login/
│   │   ├── register/
│   │   ├── verify_otp/
│   │   ├── forgot_password/
│   │   └── refresh_token/
│   │
│   ├── adapter/                # Interface Adapters (Layer 3)
│   │   ├── http/               # HTTP delivery (handlers, DTOs, middleware)
│   │   ├── persistence/        # Data persistence adapters
│   │   │   └── postgres/
│   │   ├── cache/              # Caching adapters
│   │   │   └── redis/
│   │   ├── mailer/             # Email delivery adapters
│   │   │   └── smtp/
│   │   ├── auth/               # Authentication adapters
│   │   │   └── jwt/
│   │   └── security/           # Security adapters
│   │       └── bcrypt/
│   │
│   ├── app/                    # Composition Root
│   │   └── container.go        # Dependency injection wiring
│   │
│   ├── config/                 # Configuration
│   └── shared/                 # Shared utilities (logger, response)
│
├── migrations/                 # Database migrations
└── tests/                      # Tests
```

## Layer Descriptions

### Layer 1: Domain (Enterprise Business Rules)

**Location**: `internal/domain/`

The innermost layer containing business entities, value objects, and interfaces (ports). This layer has **no dependencies** on outer layers.

- **`entity/`**: Core business entities (User, Session)
- **`valueobject/`**: Value objects (UserStatus, Email) and domain errors
- **`repository/`**: Repository interfaces defining data access contracts
- **`service/`**: Domain service interfaces organized by capability
  - Each subdirectory (`token/`, `password/`, etc.) contains a `Service` interface
  - Reduces naming stutter: `token.Service` instead of `TokenService`

### Layer 2: Use Cases (Application Business Rules)

**Location**: `internal/usecase/`

Application-specific business rules. Each use case is self-contained in its own package.

**Structure per use case**:

```
usecase/<feature>/
├── input.go      # Input DTO with Validate() method
├── output.go     # Output DTO
└── usecase.go    # Usecase interface and implementation
```

**Dependencies**: Only imports from `internal/domain/` (entities, repositories, service interfaces)

**Examples**:

- `login/` - User authentication
- `register/` - User registration
- `verify_otp/` - Email verification
- `forgot_password/` - Password reset flow
- `refresh_token/` - Token refresh

### Layer 3: Adapters (Interface Adapters)

**Location**: `internal/adapter/`

Concrete implementations of domain interfaces and delivery mechanisms.

#### HTTP Adapter

**Location**: `internal/adapter/http/`

- **`handler/`**: HTTP request handlers
- **`dto/`**: Data Transfer Objects for HTTP
- **`middleware/`**: HTTP middleware (auth, validation, logging)

#### Persistence Adapter

**Location**: `internal/adapter/persistence/postgres/`

Implements `repository.UserRepository` using PostgreSQL.

**Files**:

- `user_repository.go` - User data access implementation
- `db.go` - Database connection setup
- `migrate.go` - Migration runner

#### Cache Adapter

**Location**: `internal/adapter/cache/redis/`

Implements caching services using Redis.

**Files**:

- `otp_service.go` - Implements `otp.Service`
- `session_service.go` - Implements `session.Service`
- `client.go` - Redis client setup

#### Mailer Adapter

**Location**: `internal/adapter/mailer/smtp/`

Implements `email.Service` using SMTP.

**Files**:

- `email_service.go` - Email sending implementation

#### Auth Adapter

**Location**: `internal/adapter/auth/jwt/`

Implements `token.Service` using JWT.

**Files**:

- `token_service.go` - JWT token generation and validation

#### Security Adapter

**Location**: `internal/adapter/security/bcrypt/`

Implements `password.Service` using bcrypt.

**Files**:

- `password_service.go` - Password hashing and validation

### Composition Root

**Location**: `internal/app/container.go`

Dependency injection container that wires all components together. This is the only place where concrete implementations are instantiated and injected into use cases.

## Dependency Rule

Dependencies flow **inward only**:

```
Frameworks & Drivers (cmd, migrations)
        ↓
Interface Adapters (adapter/)
        ↓
Application Business Rules (usecase/)
        ↓
Enterprise Business Rules (domain/)
```

- **Domain** has no dependencies
- **Use Cases** depend only on **Domain**
- **Adapters** depend on **Domain** (implement interfaces)
- **Main/Container** depends on everything (wires it all together)

## Key Design Decisions

### 1. Capability-Based Service Packages

Domain services are organized by capability with concise names:

- `domain/service/token/` → `token.Service`
- `domain/service/password/` → `password.Service`

This reduces naming stutter and improves readability.

### 2. Technology-First Adapter Naming

Adapters are named after the technology they use:

- `adapter/persistence/postgres/` (not `adapter/repository/`)
- `adapter/cache/redis/` (not `adapter/otp/` or `adapter/session/`)
- `adapter/mailer/smtp/` (not `adapter/email/`)

This makes it easy to add alternative implementations (e.g., `adapter/persistence/mongodb/`).

### 3. Standardized Use Case Structure

Every use case follows the same pattern:

- `Input` with `Validate()` method
- `Output` struct
- `Usecase` interface with `Execute(ctx, input) (*output, error)`
- Constructor `New<Feature>Usecase(...) Usecase`

### 4. Import Aliases for Clarity

When package names conflict, use descriptive aliases:

```go
import (
    goredis "github.com/redis/go-redis/v9"
    "internal/adapter/cache/redis"
)
```

## Testing Strategy

- **Domain**: Unit tests for entities and value objects
- **Use Cases**: Unit tests with mocked repositories and services
- **Adapters**: Integration tests with real dependencies (testcontainers)
- **HTTP**: End-to-end tests

## Adding New Features

### 1. Add a new use case

```bash
mkdir -p internal/usecase/new_feature
touch internal/usecase/new_feature/{input,output,usecase}.go
```

### 2. Add a new domain service

```bash
mkdir -p internal/domain/service/newservice
touch internal/domain/service/newservice/service.go
```

### 3. Add a new adapter

```bash
mkdir -p internal/adapter/category/technology
touch internal/adapter/category/technology/implementation.go
```

### 4. Wire in container

Update `internal/app/container.go` to instantiate and inject dependencies.

## Benefits

1. **Testability**: Each layer can be tested independently with mocks
2. **Maintainability**: Clear boundaries and responsibilities
3. **Flexibility**: Easy to swap implementations (e.g., Postgres → MongoDB)
4. **Scalability**: New features follow established patterns
5. **Independence**: Business logic is independent of frameworks and tools

## References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture (Ports & Adapters)](https://alistair.cockburn.us/hexagonal-architecture/)

# Hexagonal Architecture – Database-agnostic Design

This project follows **Hexagonal Architecture (Ports & Adapters)** to ensure that **business logic is independent of database, framework, and transport layer**.

The goal is simple:

> **Change database or framework without touching business logic.**

---

## 🧠 Core Principles

* **Domain is the center** (business rules)
* **Infrastructure is replaceable**
* **Dependencies point inward**
* **Conversions happen at boundaries**

---

## 📂 Project Structure

```text
.
├── domain
│   └── user
│       ├── entity.go          # Business entities (pure domain)
│       ├── repository.go      # Port (interface)
│       └── service.go         # Use cases / business logic
│
├── adapter
│   └── http
│       ├── handler.go         # HTTP handlers
│       ├── request.go         # Request DTOs (JSON parsing)
│       └── response.go        # Response DTOs
│
├── infrastructure
│   └── mysql
│       ├── user_model.go      # Database model (GORM)
│       └── user_repository.go # Adapter (implements domain port)
│
├── main.go                    # Dependency wiring
└── README.md
```

---

## 🧩 Layer Responsibilities

### 1️⃣ Domain Layer (`domain/`)

**Purpose:** Business logic and rules
**Knows nothing about:** HTTP, JSON, DB, ORM, frameworks

#### Entity

```go
type User struct {
	ID    uint64
	Email string
	Name  string
}
```

* Represents business concepts
* No tags, no persistence concerns

#### Repository (Port)

```go
type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id uint64) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uint64) error
}
```

* Defines **what the domain needs**
* Not **how it is implemented**

#### Service (Use Case)

```go
type Service struct {
	repo Repository
}
```

* Orchestrates business rules
* Calls repository through interface
* Never changes when DB changes

---

### 2️⃣ Adapter Layer (`adapter/http/`)

**Purpose:** Handle I/O (HTTP, JSON)

#### Request DTO

```go
type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
```

* Used only for parsing input
* Not reused in domain or infrastructure

#### Handler

* Converts **Request → Domain Entity**
* Calls domain service
* Converts **Domain Entity → Response**

---

### 3️⃣ Infrastructure Layer (`infrastructure/`)

**Purpose:** Technical implementation details (DB, ORM, external services)

#### Database Model

```go
type UserModel struct {
	ID    uint64 `gorm:"primaryKey"`
	Email string `gorm:"uniqueIndex"`
	Name  string
}
```

* Exists only for persistence
* Can change freely with database decisions

#### Repository Adapter

* Implements domain repository interface
* Converts **Entity ⇄ Model**
* Uses GORM / SQL / any DB driver

This is the **only place** where entities interact with models.

---

## 🔄 Data Flow

### Incoming Request

```
HTTP JSON
   ↓
Request DTO
   ↓
HTTP Handler
   ↓
Domain Entity
   ↓
Service (Use Case)
   ↓
Repository Interface
   ↓
Repository Adapter
   ↓
Database Model
   ↓
Database
```

### Outgoing Response (reverse flow)

---

## 🔁 Where Conversions Happen

| Conversion         | Location                |
| ------------------ | ----------------------- |
| JSON → Request DTO | HTTP adapter            |
| Request → Entity   | HTTP handler / use case |
| Entity → Model     | Repository adapter      |
| Model → Entity     | Repository adapter      |
| Entity → Response  | HTTP handler            |

**Rule:**
👉 Conversions happen **only at boundaries**, never inside the domain.

---

## 🔒 What the Domain Must Never Know

* GORM
* SQL
* JSON tags
* HTTP status codes
* Database schema
* Frameworks

Violating this breaks hexagonal architecture.

---

## 🔄 Switching Database

To change database (e.g. MySQL → PostgreSQL):

1. Create a new adapter:

   ```
   infrastructure/postgres/
   ```
2. Implement the same repository interface
3. Change wiring in `main.go`

✅ Domain code stays untouched
✅ Business logic remains stable

---

## ✅ Why This Design Works

* Easy to test business logic with mocks
* Safe to refactor infrastructure
* Clear ownership of responsibilities
* Scales with team size
* Proven, boring, reliable
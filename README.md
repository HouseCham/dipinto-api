```text
    dipinto-api/
    │
    ├── cmd/
    │   └── ecommerce/           # Entry point of the application
    │       └── main.go
    │
    ├── internal/
    │   ├── domain/              # Core business logic, entities, and interfaces
    │   │   ├── model/           # Domain models/entities (e.g., User, Product)
    │   │   └── service/         # Business logic (e.g., OrderService)
    │   │
    │   ├── application/         # Application logic (use cases)
    │   │   ├── command/         # Command handlers (e.g., CreateOrder)
    │   │   ├── query/           # Query handlers (e.g., GetProduct)
    │   │   └── ports/           # Interfaces for external services (e.g., Repositories)
    │   │
    │   ├── infrastructure/      # External adapters (DB, HTTP, etc.)
    │   │   ├── persistence/     # Database adapters (e.g., GORM repository implementations)
    │   │   ├── http/            # HTTP server and handlers
    │   │   │   ├── routes/      # Define routes and link them to handlers
    │   │   │   ├── handlers/    # HTTP handlers (controller layer)
    │   │   │   └── middlewares/ # HTTP middlewares
    │   │   └── config/          # Configuration and environment setup
    │   │
    │   └── tests/               # Unit and integration tests
    │
    └── pkg/                     # Shared libraries or utilities (optional)
```

```json
{
    "server": {
      "port": 3000
    },
    "database": {
      "DNS": "postgres://user:password@host:port/db-name"
    },
    "JWT": {
      "secret_key": "YOUR-JWT-KEY"
    },
    "client": {
      "origin": "ORIGIN-THAT-WILL-CONSUME-API"
    }
  }
```
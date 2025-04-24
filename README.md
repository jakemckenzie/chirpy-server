# Chirpy Server

## Project Overview
documentation: todo
```
project-root/
├── cmd/
│   └── server/
│       └── main.go  # Entry point to start the server
├── internal/
│   ├── handlers/    # HTTP handlers
│   │   ├── healthz.go
│   │   ├── admin.go
│   │   └── chirp.go
│   ├── services/    # Business logic
│   │   ├── metrics.go
│   │   └── text.go
│   ├── middleware/  # Middleware logic
│   │   └── metrics.go
│   └── utils/       # Utility functions
│       └── response.go
├── docs/           # API documentation
│   └── api.md
├── go.mod          # Module definition
└── go.sum          # Dependency checksums
```
# Chirpy Server

## Project Overview
documentation: todo
```
chirpy-server/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── db.go
│   │   ├── models.go
│   │   └── users.sql.go
│   ├── handlers/
│   │   ├── healthz.go
│   │   ├── admin.go
│   │   ├── chirp.go
│   │   └── users.go
│   ├── services/
│   │   ├── metrics.go
│   │   └── text.go
│   ├── middleware/
│   │   └── metrics.go
│   └── utils/
│       └── response.go
├── sql/
│   ├── schema/
│   │   └── 001_users.sql
│   └── queries/
│       └── users.sql
├── docs/
│   └── api.md
├── .env
├── .gitignore
├── go.mod
└── go.sum
```
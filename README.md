# Chirpy Server

## Project Overview
documentation: todo
```
chirpy-server/
├── assets/
│   ├── logo.png
│   └── logo.pngZone.Identifier
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── admin.sql.go
│   │   ├── chirps.sql.go
│   │   ├── db.go
│   │   ├── models.go
│   │   └── users.sql.go
│   ├── handlers/
│   │   ├── admin.go
│   │   ├── chirp.go
│   │   ├── healthz.go
│   │   └── users.go
│   ├── middleware/
│   │   └── metrics.go
│   ├── services/
│   │   ├── metrics.go
│   │   └── text.go
│   └── utils/
│       └── response.go
├── sql/
│   ├── schema/
│   │   ├── 001_users.sql
│   │   └── 002_chirps.sql
│   └── queries/
│       ├── admin.sql
│       ├── chirps.sql
│       └── users.sql
├── static/
│   └── index.html
├── docs/
│   └── api.md
├── .env
├── .gitignore
├── go.mod
├── go.sum
├── out
├── README.md
└── sqlc.yaml
```
# Backend Server

Go backend server using sqlite3 db.

## Requirements

-   Go >= 1.18

## Configuration

Before building or running backend, create `config.json` in `backend` dir.

`config.json` example:
```json
{
  "server": {
    "host": "localhost",
    "port": 8080
  },
  "mail_sender": {
    "email": "<email>",
    "sendgrid_key": "<sendgrid-key>",
    "timeout": 30
  },
  "db": {
    "file": "./database.db"
  },
  "jwt": {
    "secret": "<secret>",
    "timeout": 5
  }
}
```

## Build

```bash
go build -o ./bin/server ./cmd/server.go
```

## Run

```bash
go run ./cmd/server.go
```

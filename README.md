# Notes API — DevOps Playground

A simple notes management REST API in Go built as a foundation for learning modern DevOps practices.

## Tech Stack

| Layer | Technology |
|---|---|
| **Backend** | Go 1.26, Chi router, pgx (PostgreSQL driver) |
| **Database** | PostgreSQL 16 |
| **Reverse proxy** | Nginx 1.26 |
| **Observability** | Prometheus + Grafana |
| **Logging** | slog (structured logging) |
| **CI/CD** | GitHub Actions |
| **Containerization** | Docker, Docker Compose |

## Architecture

```
                         ┌──────────┐
  http://localhost ───▶  │  Nginx   │
                         └────┬─────┘
                              │
                    ┌─────────▼──────────┐
                    │  Go API (port 8080) │
                    │  handlers → service │
                    │  → repository (pgx) │
                    └────┬─────────┬──────┘
                         │         │
              ┌──────────▼──┐  ┌───▼──────────┐
              │  PostgreSQL  │  │  Prometheus   │
              │   (port 5432)│  │   (port 9090) │
              └─────────────┘  └───┬───────────┘
                                   │
                            ┌──────▼──────┐
                            │   Grafana    │
                            │  (port 3000) │
                            └─────────────┘
```

## API Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Health check (includes DB status) |
| GET | `/api/notes` | List all notes |
| POST | `/api/notes` | Create a note |
| GET | `/api/notes/{id}` | Get a note by ID |
| PUT | `/api/notes/{id}` | Update a note |
| DELETE | `/api/notes/{id}` | Delete a note |
| GET | `/metrics` | Prometheus metrics |

## Quick start

```bash
docker compose up --build -d
```

Then open:

- **App:** http://localhost
- **Prometheus:** http://localhost:9090
- **Grafana:** http://localhost:3000

### Stop

```bash
docker compose down -v
```

## Project structure

```
cmd/api/main.go          Entry point
internal/
  config/                Env configuration
  handlers/              HTTP handlers + errors
  logger/                Structured logging (slog)
  metrics/               Prometheus metrics
  middleware/            Logging, recovery, metrics
  model/                 Data models
  repository/            PostgreSQL data access
  service/               Business logic + validation
web/                     Frontend (HTML/JS SPA)
migrations/              SQL migrations
nginx/                   Reverse proxy config
prometheus/              Prometheus scrape config
grafana/                 Dashboards (auto-provisioned)
.github/workflows/       CI/CD pipeline
```

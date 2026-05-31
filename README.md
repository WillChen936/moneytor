# Moneytor

A personal finance tracker REST API built with Go and PostgreSQL.

## Tech Stack

- **Language**: Go
- **Framework**: Gin
- **Database**: PostgreSQL
- **Query generation**: sqlc
- **Migrations**: golang-migrate
- **Auth**: JWT (access token 15m + refresh token 7d)
- **Testing**: gomock (API layer), real DB (DB layer)

## Getting Started

### Prerequisites

- Go 1.21+
- Docker (for running PostgreSQL)

### Setup

```bash
# Copy the example config
cp config.example.json config.json
# Edit config.json with your database connection and token secret
# Generate a token secret:
openssl rand -base64 32
```

### Run

```bash
# Start PostgreSQL
make postgres

# Apply migrations
make migrateup

# Start the server
make server
```

## API

Base URL: `/api/v1`

### Auth

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | `/users` | Register | — |
| POST | `/users/login` | Login, returns access + refresh token | — |
| POST | `/users/refresh` | Exchange refresh token for a new access token | — |

### Accounts

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | `/accounts` | Create an account | ✓ |
| GET | `/accounts` | List accounts | ✓ |

### Categories

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | `/categories` | Create a category | ✓ |
| GET | `/categories` | List categories | ✓ |

### Entries (Income / Expense)

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | `/entries` | Create an entry | ✓ |
| GET | `/entries` | List entries for an account | ✓ |

### Transfers

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | `/transfers` | Create a transfer between accounts | ✓ |

### Reference Data

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| GET | `/transaction-types` | List transaction types (income/expense/transfer) | — |
| GET | `/currencies` | List supported currencies | — |

## Development Commands

```bash
make test          # Run all tests
make sqlc          # Regenerate sqlc code (run after editing SQL queries)
make mockdb        # Regenerate mock (run after editing Store interface)
make migrateup     # Apply all migrations
make migratedown   # Roll back all migrations
```

## Project Structure

```
api/                  # Gin handlers, one file per resource
database/
  migrations/         # golang-migrate SQL files (up/down)
  queries/            # Raw SQL used by sqlc for code generation
  sqlc/               # Generated code (do NOT edit *.sql.go)
    store.go              # Store interface + SQLStore
    tx_create_entry.go    # Income/expense transaction
    tx_create_transfer.go # Transfer transaction (deadlock-safe)
  mocks/              # mockgen-generated Store mock (used in API tests)
token/                # JWT implementation
utils/                # Config, utilities
```

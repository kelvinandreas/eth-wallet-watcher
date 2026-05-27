# ETH Wallet Watcher — Backend

REST API for monitoring Ethereum wallets. Automatically polls transaction history via Etherscan every 5 minutes and sends notifications on new activity.

## Tech Stack

- **Go** + [Fiber](https://github.com/gofiber/fiber) (HTTP)
- **PostgreSQL** (primary database)
- **Redis** + [Asynq](https://github.com/hibiken/asynq) (background job queue & scheduler)
- **GORM** (ORM)
- **JWT** (authentication)

## Features

- Register & login with JWT
- Add / remove Ethereum wallets to monitor
- Background worker polls Etherscan every 5 minutes per wallet
- Notifications generated on new transactions
- Transaction & notification history with pagination

---

## Running with Docker (recommended)

### Prerequisites

- Docker & Docker Compose

### Setup

1. Copy the example env file and fill in the required secrets:

   ```bash
   cp .env.example .env
   ```

2. Start all services (API + PostgreSQL + Redis):
   ```bash
   docker compose up --build
   ```

The API will be available at `http://localhost:8080`.

---

## Running Locally

### Prerequisites

- Go 1.24+
- PostgreSQL 15
- Redis 7

### Setup

1. Copy and fill in the env file:

   ```bash
   cp .env.example .env
   ```

2. Start dependencies:

   ```bash
   docker compose up postgres redis
   ```

3. Run the API:
   ```bash
   go run ./cmd/main.go
   ```

---

## Environment Variables

| Variable             | Description                   | Example                           |
| -------------------- | ----------------------------- | --------------------------------- |
| `DB_HOST`            | PostgreSQL host               | `localhost`                       |
| `DB_USER`            | PostgreSQL user               | `admin`                           |
| `DB_PASSWORD`        | PostgreSQL password           | `password123`                     |
| `DB_NAME`            | PostgreSQL database name      | `watcherDB`                       |
| `DB_PORT`            | PostgreSQL port               | `5432`                            |
| `REDIS_ADDR`         | Redis address                 | `localhost:6379`                  |
| `ETHERSCAN_API_KEY`  | Etherscan API key             | `your_key_here`                   |
| `ETHERSCAN_BASE_URL` | Etherscan base URL            | `https://api.etherscan.io/v2/api` |
| `JWT_SECRET`         | Secret for signing JWT tokens | `your_secret_here`                |
| `TOKEN_TTL`          | JWT expiry in hours           | `24`                              |

> When running via Docker Compose, `DB_HOST` and `REDIS_ADDR` are automatically set to the internal service names (`postgres`, `redis`).

---

## API Endpoints

All protected routes require `Authorization: Bearer <token>` header.

### Auth

| Method | Path             | Auth | Description    |
| ------ | ---------------- | ---- | -------------- |
| POST   | `/auth/register` | No   | Register user  |
| POST   | `/auth/login`    | No   | Login, get JWT |

### Wallets

| Method | Path                              | Auth | Description                   |
| ------ | --------------------------------- | ---- | ----------------------------- |
| POST   | `/wallets`                        | Yes  | Add wallet to watch           |
| GET    | `/wallets`                        | Yes  | List monitored wallets        |
| DELETE | `/wallets/:walletID`              | Yes  | Remove wallet                 |
| GET    | `/wallets/:walletID/transactions` | Yes  | Get transactions for a wallet |

### Notifications

| Method | Path                           | Auth | Description            |
| ------ | ------------------------------ | ---- | ---------------------- |
| GET    | `/notifications`               | Yes  | List notifications     |
| PATCH  | `/notifications/:notifID/read` | Yes  | Mark notification read |

---

## Project Structure

```
cmd/
  main.go                  # Entry point, wiring, routes
internal/
  config/                  # Env config loading
  constant/                # App-wide constants
  handler/                 # HTTP handlers
  helper/                  # Utilities (pagination, cache, etc.)
  infrastructure/          # DB & Redis connections
  middleware/              # JWT middleware
  model/
    app/                   # GORM models
    request/               # Request DTOs
    response/              # Response DTOs
  repository/              # Database queries
  service/                 # Business logic
  worker/                  # Asynq background jobs & scheduler
```

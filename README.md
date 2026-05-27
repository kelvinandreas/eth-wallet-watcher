# ETH Wallet Watcher

A backend service that automatically monitors Ethereum wallet transactions. Every 5 minutes, the system polls the latest transactions from Etherscan for all registered wallets and creates notifications for the user.

## Tech Stack

- **Go** + [Fiber](https://github.com/gofiber/fiber) (HTTP framework)
- **PostgreSQL** (primary database)
- **Redis** + [Asynq](https://github.com/hibiken/asynq) (background job queue & scheduler)
- **GORM** (ORM)
- **Etherscan API** (transaction data source)
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

1. Navigate to the backend directory:

   ```bash
   cd backend
   ```

2. Copy the example env file and fill in the required secrets:

   ```bash
   cp .env.example .env
   ```

3. Start all services (API + PostgreSQL + Redis):

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

1. Navigate to the backend directory:

   ```bash
   cd backend
   ```

2. Copy and fill in the env file:

   ```bash
   cp .env.example .env
   ```

3. Start dependencies:

   ```bash
   docker compose up postgres redis
   ```

4. Install dependencies:

   ```bash
   go mod tidy
   ```

5. Run the API:

   ```bash
   go run ./cmd/main.go
   ```

The database schema is automatically migrated on startup.

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

> Get a free Etherscan API key at https://etherscan.io/myapikey
> When running via Docker Compose, `DB_HOST` and `REDIS_ADDR` are automatically set to the internal service names (`postgres`, `redis`).

---

## API Endpoints

A Postman collection is included for easy testing: `eth-wallet-watcher.postman_collection.json`

**How to import:**

1. Open Postman
2. Click **Import** → select `eth-wallet-watcher.postman_collection.json`
3. Run **Login** first — the token is automatically saved to the collection variables
4. All subsequent requests will use the token automatically

All protected routes require `Authorization: Bearer <token>` header.

### Auth

| Method | Path             | Auth | Description    |
| ------ | ---------------- | ---- | -------------- |
| POST   | `/auth/register` | No   | Register user  |
| POST   | `/auth/login`    | No   | Login, get JWT |

### Wallets

| Method | Path                              | Auth | Description                               |
| ------ | --------------------------------- | ---- | ----------------------------------------- |
| POST   | `/wallets`                        | Yes  | Add wallet to watch                       |
| GET    | `/wallets`                        | Yes  | List monitored wallets                    |
| DELETE | `/wallets/:walletID`              | Yes  | Remove wallet                             |
| GET    | `/wallets/:walletID/transactions` | Yes  | Get transactions for a wallet (paginated) |

### Notifications

| Method | Path                           | Auth | Description                    |
| ------ | ------------------------------ | ---- | ------------------------------ |
| GET    | `/notifications`               | Yes  | List notifications (paginated) |
| PATCH  | `/notifications/:notifID/read` | Yes  | Mark notification read         |

Paginated endpoints support `?page=1&limit=10` query params (default: page 1, limit 10, max 100).

---

## How the Background Worker Works

Every **5 minutes**, the system automatically:

1. Fetches all registered wallets
2. Pulls new transactions from Etherscan (starting from the last checked block)
3. Saves new transactions to the database
4. Creates a notification for each new transaction
5. Updates `last_block` on the wallet

The worker runs alongside the HTTP server — no separate process needed.

---

## Project Structure

```
backend/
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

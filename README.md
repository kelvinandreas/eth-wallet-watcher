# ETH Wallet Watcher

A backend service that automatically monitors Ethereum wallet transactions. Every 5 minutes, the system polls the latest transactions from Etherscan for all registered wallets and creates notifications for the user.

## Tech Stack

- **Go** + Fiber (HTTP framework)
- **PostgreSQL** (database)
- **Redis** + Asynq (background job queue & scheduler)
- **Etherscan API** (transaction data source)

## Prerequisites

- Go 1.21+
- Docker & Docker Compose

## Getting Started

### 1. Navigate to the backend directory

```bash
cd backend
```

### 2. Set up environment variables

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
DB_HOST=localhost
DB_USER=admin
DB_PASSWORD=password123
DB_NAME=watcherDB
DB_PORT=5432
REDIS_ADDR=localhost:6379
ETHERSCAN_API_KEY=your_etherscan_api_key
ETHERSCAN_BASE_URL=https://api.etherscan.io/v2/api
JWT_SECRET=your_jwt_secret_key
TOKEN_TTL=24
```

> Get a free Etherscan API key at https://etherscan.io/myapikey

### 3. Start PostgreSQL & Redis

```bash
docker-compose up -d
```

### 4. Install dependencies

```bash
go mod tidy
```

### 5. Run the server

```bash
go run cmd/main.go
```

The server runs at `http://localhost:8080`. The database schema is automatically migrated on startup.

---

## API Endpoints

A Postman collection is included for easy testing: `eth-wallet-watcher.postman_collection.json`

**How to import:**
1. Open Postman
2. Click **Import** → select `eth-wallet-watcher.postman_collection.json`
3. Run **Login** first — the token is automatically saved to the collection variables
4. All subsequent requests will use the token automatically

**Available endpoints:**

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/auth/register` | — | Register a new user |
| POST | `/auth/login` | — | Login and get JWT token |
| POST | `/wallets` | ✓ | Add a wallet to monitor |
| GET | `/wallets` | ✓ | Get all monitored wallets |
| DELETE | `/wallets/:walletID` | ✓ | Remove a wallet |
| GET | `/wallets/:walletID/transactions` | ✓ | Get transactions for a wallet (paginated) |
| GET | `/notifications` | ✓ | Get all notifications (paginated) |
| PATCH | `/notifications/:notifID/read` | ✓ | Mark a notification as read |

Endpoints marked **(paginated)** support `?page=1&limit=10` query params (default: page 1, limit 10, max 100).

---

## How the Background Worker Works

Every **5 minutes**, the system automatically:
1. Fetches all registered wallets
2. Pulls new transactions from Etherscan (starting from the last checked block)
3. Saves new transactions to the database
4. Creates a notification for each new transaction
5. Updates `last_block` on the wallet

The worker runs alongside the HTTP server when `go run cmd/main.go` is executed — no separate process needed.

# Sakutera

Backend API for Sakutera — a platform that helps informal/gig workers (freelancers, ride-hailing drivers, etc.) record their daily income, forecast future cash flow (via an external ML service), and issue an **Income Passport**: a shareable proof of income history (e.g. for loan applications), complete with a consent & access-log system.

---

## Key Features

- **Phone number + WhatsApp OTP authentication**, followed by a 6-digit PIN for login (`pkg/whatsapp`, `pkg/bcrypt`, `pkg/jwt`)
- **Onboarding** — select work category/platform and income source
- **Transaction recording** with hash-chain verification (`previous_hash` / `current_hash`) to prevent data tampering
- **Dashboard & ledger** income summaries
- **Forecasting** — integration with an external ML service to predict EMI (estimated monthly income), trend, and deficit risk (`pkg/mlclient`)
- **Income Passport** — issuance of an income summary document (3/6/12 months) based on transaction data and forecast results, requiring at least 30 days of data
- **Consent & Access Log** — passport owners can grant/revoke access to third-party organizations and view a log of who accessed their data
- **Token revocation** via `TokenBlacklist` so logout actually invalidates the JWT

Full endpoint documentation (with example requests/responses) is in [`docs/API.md`](docs/API.md).

---

## Tech Stack

| Package                                                                          | Purpose                                          |
| --------------------------------------------------------------------------------- | -------------------------------------------------- |
| [Gin](https://github.com/gin-gonic/gin)                                          | HTTP web framework                                |
| [GORM](https://gorm.io) + [gorm/driver/mysql](https://github.com/go-gorm/mysql)  | ORM, MariaDB/MySQL connection                     |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt)                              | JWT authentication                                |
| [google/uuid](https://github.com/google/uuid)                                    | UUID primary keys for all tables                  |
| [golang.org/x/crypto (bcrypt)](https://pkg.go.dev/golang.org/x/crypto)           | PIN hashing                                       |
| [joho/godotenv](https://github.com/joho/godotenv)                                | `.env` loading                                    |
| [gin-contrib/cors](https://github.com/gin-contrib/cors)                          | CORS middleware                                   |
| MariaDB 11.4                                                                      | Primary database (run via Docker Compose)         |
| Redis 7.4                                                                         | Provided in `docker-compose.yml` for cache/session use |

---

## Folder Structure

```
sakutera-softdev/
├── cmd/app/main.go              # Entry point: dependency wiring & server startup
├── entity/                      # GORM models (user, transaction, forecast_result,
│                                 #   income_passport, consent, access_log, organization, etc.)
├── internal/
│   ├── handler/rest/            # HTTP layer (Gin handlers), grouped by domain:
│   │                             #   auth.go, onboarding.go, transaction.go, dashboard.go,
│   │                             #   passport.go, access.go
│   ├── repository/              # Data access layer (GORM queries per entity)
│   └── service/                 # Business logic (auth, onboarding, transaction, dashboard,
│                                 #   passport, access)
├── model/                       # Request/response DTOs
├── pkg/
│   ├── bcrypt/                  # PIN hashing (cost=10)
│   ├── config/                  # .env loading & DSN builder
│   ├── database/mariadb/        # DB connection, AutoMigrate, initial data seeding
│   ├── errors/                  # Custom AppError type with standardized HTTP status codes
│   ├── jwt/                     # JWT creation & validation
│   ├── middleware/               # CORS, auth guard
│   ├── mlclient/                # HTTP client for the ML forecasting service
│   ├── response/                # Standardized JSON response envelope
│   └── whatsapp/                # Sends OTP via Fonnte (WhatsApp API)
├── docs/API.md                  # Full endpoint documentation
├── docker-compose.yml           # App + MariaDB + Redis for deployment
├── Dockerfile
├── Makefile
└── .env.example
```

---

## Architecture

3-layer, following the `handler → service → repository` pattern:

```
HTTP Request
     │
     ▼
Handler (internal/handler/rest)   → parses requests, calls service, formats response
     │
     ▼
Service (internal/service)        → business logic, orchestrates repository/pkg calls
     │
     ▼
Repository (internal/repository)  → the only layer allowed to call GORM
     │
     ▼
Database (MariaDB)
```

All dependencies are wired manually (constructor injection) in `cmd/app/main.go` — no DI framework.

---

## Running Locally

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (for MariaDB + Redis)
- (Optional) [air](https://github.com/air-verse/air) for hot-reload — config already exists at `.air.toml`

### Steps

1. **Clone the repo**

   ```bash
   git clone https://github.com/azmiagr/sakutera-softdev.git
   cd sakutera-softdev
   ```

2. **Set up environment variables**

   ```bash
   make init
   # equivalent to: cp .env.example .env && go mod tidy
   ```

   Then edit `.env` as needed:

   | Variable                     | Description                                            | Example                          |
   | ----------------------------- | -------------------------------------------------------- | ---------------------------------- |
   | `DB_HOST`                     | Database host                                            | `localhost`                       |
   | `DB_PORT`                     | Database port                                            | `3306`                            |
   | `DB_NAME`                     | Database name                                            | `sakutera`                        |
   | `DB_USER`                     | Database user                                            | `sakutera_user`                   |
   | `DB_PASSWORD`                 | Database password                                        | `sakutera_pass`                   |
   | `ADDRESS`                     | Server bind address                                      | `localhost`                       |
   | `PORT`                        | Server port                                              | `8080`                            |
   | `TIME_OUT_LIMIT`              | Request timeout (seconds)                                | `10`                               |
   | `JWT_SECRET_KEY`              | Secret used to sign JWTs (minimum 256-bit)               | a long random string               |
   | `JWT_EXP_TIME`                | JWT expiration (hours)                                   | `1`                                 |
   | `FONNTE_TOKEN`                | Fonnte API token for sending OTP via WhatsApp             | token from the Fonnte dashboard     |
   | `ALLOWED_ORIGINS`             | Allowed CORS origins (comma-separated)                    | `http://localhost:3000`            |
   | `ML_SERVICE_URL`              | Base URL of the ML forecasting service                    | `http://ml-service:8000`           |
   | `ML_SERVICE_TOKEN`            | Auth token for the ML service                              | internal token                     |
   | `ML_REQUEST_TIMEOUT_SECONDS`  | Timeout for requests to the ML service (seconds)           | `15`                                 |

3. **Start the database (MariaDB) & Redis**

   Use the existing `docker-compose.yml` (run only the `db` and `redis` services for local dev):

   ```bash
   docker compose up -d db redis
   ```

   Make sure `DB_HOST` in `.env` points to `localhost` and `DB_PORT` matches the port Compose exposes (`3307` by default — see `docker-compose.yml`), or adjust accordingly.

4. **Run the server**

   ```bash
   make run
   # equivalent to: go run cmd/app/main.go
   ```

   On startup, the app automatically runs `mariadb.Migrate()` (AutoMigrate on all entities) and `mariadb.Seed()` (seeds initial data such as work categories/platforms).

   For hot-reload during development:

   ```bash
   air
   ```

5. **The server is now available** at `http://localhost:8080/api/v1` (adjust based on `ADDRESS`/`PORT` in `.env`).

6. **Testing**

   ```bash
   make test
   # equivalent to: go test ./...
   ```

---

## Running the Full Stack via Docker Compose

To run the app + MariaDB + Redis together in containers (image pulled from GHCR):

```bash
docker compose up -d
```

See `deploy.sh` for the script used on the production server (pulls the latest image, restarts services, prunes old images).

---

## Endpoints

Base URL: `/api/v1`. See [`docs/API.md`](docs/API.md) for full request/response details.

| Method | Path                                   | Auth | Description                          |
| ------ | ---------------------------------------- | ---- | --------------------------------------- |
| POST   | `/auth/register`                        | -    | Register a new account, send OTP        |
| POST   | `/auth/verify-otp`                      | -    | Verify OTP                              |
| POST   | `/auth/check-phone`                     | -    | Check if phone number is registered (login flow) |
| POST   | `/auth/set-pin`                         | -    | Create a 6-digit PIN                    |
| POST   | `/auth/login`                           | -    | Log in with phone number + PIN          |
| POST   | `/auth/logout`                          | ✔    | Log out (revokes token via blacklist)   |
| GET    | `/onboarding/work-categories`           | ✔    | List work categories                    |
| POST   | `/onboarding/work-platform`             | ✔    | Select work platform                    |
| GET    | `/onboarding/income-sources`            | ✔    | List income sources                     |
| POST   | `/onboarding/income-source`             | ✔    | Select income source                    |
| GET    | `/dashboard`                            | ✔    | Dashboard summary                       |
| GET    | `/ledger`                               | ✔    | Transaction ledger                      |
| GET    | `/transactions/sources`                 | ✔    | List transaction sources                |
| POST   | `/transactions`                         | ✔    | Record a new transaction                |
| GET    | `/transactions`                         | ✔    | List transactions                       |
| GET    | `/passport`                             | ✔    | Get the active income passport          |
| GET    | `/passport/preview`                     | ✔    | Preview a passport before issuing it    |
| POST   | `/passport`                             | ✔    | Issue an income passport                |
| GET    | `/passport/access`                      | ✔    | List passport access consents           |
| POST   | `/passport/access`                      | ✔    | Grant access to an organization         |
| PATCH  | `/passport/access/:consent_id/revoke`   | ✔    | Revoke access                           |
| GET    | `/passport/access/logs`                 | ✔    | View who has accessed the data          |
| GET    | `/organizations`                        | ✔    | List third-party organizations          |

Endpoints marked Auth ✔ require an `Authorization: Bearer <token>` header obtained from login/OTP verification.

---

## Demo Account

To try the income passport feature without manually recording 30 days of transactions, a demo account is available with pre-seeded data (30 days of transactions + a forecast result):

- **Phone number:** `081234567890`
- **PIN:** `123456`

Log in via `POST /auth/login`:

```json
{
  "phone_number": "081234567890",
  "pin": "123456"
}
```

After logging in, use the returned token to try `GET /passport`, `GET /passport/preview`, and `POST /passport` directly.

# Chirpy: A Minimalist Twitter Clone in Go

Chirpy is a backend-focused microblogging platform (similar to Twitter) built in Go. It allows users to register, authenticate with secure JWT-based sessions, post chirps (short messages), and manage tokens with a robust refresh/revoke mechanism. The project also features webhook handling, user upgrades, and admin tools.

---

## Features

### ✅ User Authentication

* **Register/Login** with email and password
* JWT-based **access tokens** (1 hour expiry)
* **Refresh tokens** with 60-day expiry
* Secure password hashing using **bcrypt**

### ✍️ Chirps (Posts)

* Authenticated users can **create**, **delete**, and **retrieve** chirps
* Filter chirps by `author_id` and sort by creation time (asc/desc)

### 🔁 Token Management

* Auto-generate and validate JWT tokens
* Refresh access tokens using refresh token endpoint
* Revoke tokens to log out or invalidate sessions

### 🧠 Webhook Integration

* Handle third-party payment events (e.g., `user.upgraded`)
* Upgrades user status with `is_chirpy_red = true`

### 🛠 Admin Utilities

* `POST /admin/reset`: Reset hit counter
* Environment-based reset behavior (`PLATFORM=dev`)

---

## API Endpoints

### User

* `POST /api/users`: Register
* `POST /api/login`: Login & receive tokens
* `POST /api/refresh`: Refresh JWT using a refresh token
* `POST /api/revoke`: Revoke a refresh token
* `PUT /api/users`: Update password/email (planned)

### Chirps

* `POST /api/chirps`: Create a chirp (requires valid JWT)
* `GET /api/chirps`: Fetch all chirps or filter by `author_id`
* `GET /api/chirps/{chirpID}`: Get single chirp
* `DELETE /api/chirps/{chirpID}`: Delete your chirp

### Webhooks

* `POST /api/polka/webhooks`: Handle payment upgrade events

### Admin

* `POST /admin/reset`: Reset visit counter (dev-only)
* `GET /admin/metrics`: HTML metrics page

---

## Database Schema

### `users`

* `id UUID PRIMARY KEY`
* `email TEXT UNIQUE`
* `hashed_password TEXT`
* `is_chirpy_red BOOLEAN`
* `created_at`, `updated_at`

### `chirps`

* `id UUID PRIMARY KEY`
* `user_id UUID REFERENCES users(id)`
* `body TEXT`
* `created_at`, `updated_at`

### `refresh_tokens`

* `token TEXT PRIMARY KEY`
* `user_id UUID REFERENCES users(id)`
* `created_at`, `updated_at`, `expires_at`, `revoked_at`

---

## Tech Stack

* **Go** (Golang)
* **PostgreSQL** with `sqlc`
* **bcrypt** for hashing
* **JWT** for authentication
* **Goose** for migrations
* **Chi Router** for routing

---

## How to Run

```bash
# 1. Set environment variables
export JWT_SECRET="your_jwt_secret"
export DB_URL="your_postgres_connection_url"
export PLATFORM=dev

# 2. Run migrations
make migrateup

# 3. Run the server
go run .
```

---

## Author

**Asylbek Zhunusov**
Boot.dev Certified | Backend Developer | Curious Learner

Feel free to connect with me or check out more projects on my GitHub!

---

## License

This project is licensed under the MIT License.

---

## Notes

* Make sure your `.env` file contains valid secrets and database connection strings.
* Tests can be executed using CLI-based test suites.
* Consider containerizing the app for production with Docker.

---

Happy Chipping! 🐦

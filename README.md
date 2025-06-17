# Chirpy Server

A simple Go web server that powers the Chirpy frontend and serves static assets. It includes a readiness endpoint for use with load balancers or orchestration platforms.

---

## 🚀 Features

- Serves static HTML files (like `index.html`) from the `/app/` path
- Serves static assets (like images, CSS, JS) from `/assets/`
- Provides a readiness endpoint at `/healthz` for health checks

---

## 🗂 Project Structure
your-project/
├── main.go              # Main Go server file
├── index.html           # Static HTML file served at /app/index.html
└── assets/              # Static asset directory
    └── logo.png         # Logo served at /assets/logo.png



---

## 🔧 Setup and Usage

### 1. Build the server

```bash
go build -o chirpy

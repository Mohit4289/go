# Go Gin Authentication API & Testing Guide

This backend is a production-ready authentication and user management service built with **Go**, **Gin Web Framework**, **PostgreSQL** (via `pgx`), **Argon2id** password hashing, and **JWT** (`golang-jwt/jwt/v5`).

This document provides a complete specification of all API endpoints and instructions on how to test the service using **`apitest`** (CLI, HTTP API sidecar, or Browser Extension).

---

## Table of Contents

1. [Quickstart: Running the Go Backend](#quickstart-running-the-go-backend)
2. [API Specification Files](#api-specification-files)
3. [Endpoint Reference & Schemas](#endpoint-reference--schemas)
4. [Testing with `apitest` CLI](#testing-with-apitest-cli)
5. [Testing via `apitest serve` (HTTP API)](#testing-via-apitest-serve-http-api)
6. [Testing via Chrome Extension](#testing-via-chrome-extension)
7. [Manual cURL Testing Flow](#manual-curl-testing-flow)

---

## Quickstart: Running the Go Backend

### 1. Environment Configuration
Ensure your [`.env`](file:///home/guts/go/.env) file is configured in the root directory:
```env
Database_URL="postgresql://<user>:<password>@<host>:5432/<database>"
JWT_SECRET="your-256-bit-secret-key"
```

### 2. Start the Server
```bash
cd /home/guts/go
go run main.go
```
The server will start listening on **`http://localhost:8080`**.

---

## API Specification Files

We have created two standardized test specification files ready for `apitest`:

| Format | File Path | Description |
|---|---|---|
| **OpenAPI 3.0** | [`openapi.yaml`](file:///home/guts/go/openapi.yaml) | Full OpenAPI 3.0.3 YAML definition with request bodies and expected response status codes. |
| **Postman v2.1** | [`postman_collection.json`](file:///home/guts/go/postman_collection.json) | Postman Collection with configured endpoints, headers, and request samples. |

*(These files are also mirrored in [`/home/guts/apitester/testdata/`](file:///home/guts/apitester/testdata/) as `go-auth-api.yaml` and `go-auth-api.postman.json`).*

---

## Endpoint Reference & Schemas

Base URL: `http://localhost:8080`

| Method | Endpoint | Auth Required | How Backend Accepts Values | Expected Status | Description |
|---|---|---|---|---|---|
| `POST` | `/auth/user` | No | JSON Body `{"name", "email", "pass"}` | `201 Created` | Register a new user |
| `POST` | `/auth/login` | No | JSON Body `{"email", "pass"}` | `200 OK` | Authenticates user, sets `access_token` & `refresh_token` cookies, returns JWT `token` in JSON |
| `POST` | `/auth/refresh_token` | Yes | HTTP Cookie: `refresh_token` | `200 OK` | Verifies refresh token hash in DB & issues new `access_token` cookie |
| `GET` | `/auth/checkuser` | Yes | HTTP Cookie: `access_token` | `200 OK` | Validates JWT `user_id` claim & queries all users from PostgreSQL |
| `GET` | `/auth/userdata` | Yes | HTTP Cookie: `access_token` | `200 OK` | Validates JWT `user_id` claim & fetches caller profile from DB |
| `POST` | `/auth/logout` | Yes | HTTP Cookies: `access_token` & `refresh_token` | `200 OK` | Deletes token from database & expires auth cookies |

---

### How Authentication Works in this Backend

1. **Cookie-Based Auth**: The backend uses Gin's `c.Cookie("access_token")` for protected routes (`/auth/checkuser`, `/auth/userdata`, `/auth/logout`) and `c.Cookie("refresh_token")` for token refresh and logout.
2. **Seamless Testing with `apitest`**: When you supply your JWT token to `apitest` (via the `--auth "<token>"` flag or via the Chrome Extension's **Authorization** field), `apitest` automatically sets:
   - `Authorization: Bearer <token>`
   - `Cookie: access_token=<token>; refresh_token=<token>`
   This allows `apitest` to authenticate all protected routes on your backend directly without requiring any changes to your backend code!

---

### Detailed Endpoint Specifications

#### 1. Register User
- **Method / Path:** `POST /auth/user`
- **Headers:** `Content-Type: application/json`
- **Request Body:**
  ```json
  {
    "name": "Alex Mercer",
    "email": "alex.mercer@example.com",
    "pass": "SecurePass123!"
  }
  ```
- **Responses:**
  - `201 Created`: User successfully registered.
  - `400 Bad Request`: Missing or invalid fields.
  - `409 Conflict`: User with email already exists.
  - `500 Internal Server Error`: Database or hashing error.

#### 2. User Login
- **Method / Path:** `POST /auth/login`
- **Headers:** `Content-Type: application/json`
- **Request Body:**
  ```json
  {
    "email": "alex.mercer@example.com",
    "pass": "SecurePass123!"
  }
  ```
- **Responses:**
  - `200 OK`: Returns `{ "message": "login successful", "token": "..." }` and sets `access_token` and `refresh_token` HTTP-only cookies.
  - `401 Unauthorized`: Invalid email or password.

#### 3. Refresh Token
- **Method / Path:** `POST /auth/refresh_token`
- **Cookies / Headers:** `refresh_token=<token>`
- **Responses:**
  - `200 OK`: Returns new access token.
  - `400 Bad Request`: Refresh token missing.
  - `401 Unauthorized`: Invalid or expired refresh token.

#### 4. Check All Users
- **Method / Path:** `GET /auth/checkuser`
- **Auth:** `access_token` cookie or Bearer header
- **Responses:**
  - `200 OK`: Returns `{ "users": [...], "user_id": 1 }`
  - `401 Unauthorized`: Missing or invalid token.

#### 5. Get User Profile (Me)
- **Method / Path:** `GET /auth/userdata`
- **Auth:** `access_token` cookie or Bearer header
- **Responses:**
  - `200 OK`: Returns `{ "user": { "id": 1, "name": "...", "email": "..." } }`
  - `401 Unauthorized`: Missing or invalid token.
  - `404 Not Found`: User does not exist.

#### 6. Logout
- **Method / Path:** `POST /auth/logout`
- **Auth:** `access_token` cookie + `refresh_token` cookie
- **Responses:**
  - `200 OK`: `{ "message": "logged out successfully" }` (clears cookies)
  - `401 Unauthorized`: Missing authentication.

---

## Testing with `apitest` CLI

Navigate to the `apitester` directory or run the binary:

### 1. Test using OpenAPI Specification (Terminal Table)
```bash
cd /home/guts/apitester
./apitest run --spec /home/guts/go/openapi.yaml --base-url http://localhost:8080
```
**Sample Output:**
```
STATUS   METHOD   PATH                 EXPECTED   ACTUAL   LATENCY   DETAILS
------   ------   ----                 --------   ------   -------   -------
✅ PASS   POST     /auth/user           201        201      18ms      
✅ PASS   POST     /auth/login          200        200      15ms      
✅ PASS   POST     /auth/refresh_token  200        200      12ms      
✅ PASS   GET      /auth/checkuser      200        200      10ms      
✅ PASS   GET      /auth/userdata       200        200      11ms      
✅ PASS   POST     /auth/logout         200        200      14ms      

--- Test Run Summary ---
  Total:   6
  Passed:  6
  Failed:  0
  Unknown: 0
  Errors:  0
```

### 2. Test with JSON Output (for CI/CD Pipelines)
```bash
./apitest run --spec /home/guts/go/openapi.yaml --base-url http://localhost:8080 --json
```

### 3. Test with Authenticated Endpoints
If testing endpoints that require authorization headers:
```bash
./apitest run \
  --spec /home/guts/go/openapi.yaml \
  --base-url http://localhost:8080 \
  --auth "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 4. Test using Postman Collection
```bash
./apitest run \
  --spec /home/guts/go/postman_collection.json \
  --base-url http://localhost:8080
```

---

## Testing via `apitest serve` (HTTP API)

You can run `apitest` as an HTTP microservice / sidecar server to trigger tests programmatically.

### 1. Start the HTTP Sidecar
```bash
cd /home/guts/apitester
./apitest serve --port 8085
```

### 2. Run Test Suite via HTTP POST
```bash
curl -X POST http://localhost:8085/api/run \
  -H "Content-Type: application/json" \
  -d '{
    "spec": "/home/guts/go/openapi.yaml",
    "base_url": "http://localhost:8080",
    "timeout": 5,
    "concurrency": 10
  }'
```

### 3. Check Server Health
```bash
curl http://localhost:8085/api/health
```

---

## Testing via Chrome Extension

1. Start the backend: `go run main.go` on `http://localhost:8080`.
2. Start `apitest serve --port 8085`.
3. Open the **API Tester Chrome Extension** popup.
4. Select or paste the spec file path (`/home/guts/go/openapi.yaml`) or raw YAML.
5. Set Base URL to `http://localhost:8080`.
6. Click **Run Tests** to see live latency graphs, status codes, and test results!

---

## Manual cURL Testing Flow

To test the complete authentication cycle manually:

### 1. Register
```bash
curl -i -X POST http://localhost:8080/auth/user \
  -H "Content-Type: application/json" \
  -d '{"name":"Alex","email":"alex@example.com","pass":"Secret123"}'
```

### 2. Login (Saves cookies to `cookies.txt`)
```bash
curl -i -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"alex@example.com","pass":"Secret123"}'
```

### 3. Access Protected User Profile
```bash
curl -i -X GET http://localhost:8080/auth/userdata \
  -b cookies.txt
```

### 4. Refresh Token
```bash
curl -i -X POST http://localhost:8080/auth/refresh_token \
  -b cookies.txt \
  -c cookies.txt
```

### 5. Logout
```bash
curl -i -X POST http://localhost:8080/auth/logout \
  -b cookies.txt \
  -c cookies.txt
```

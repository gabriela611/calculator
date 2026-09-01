# Calculator Application

A full-stack calculator application featuring a **React 19 + TypeScript** frontend and a **Go microservice** backend. The application supports standard and advanced arithmetic operations with strict input validation, consistent JSON error responses, unit test suites, and Docker deployment support.

---

## Tech Stack

- **Frontend:** React 19, TypeScript, Vite, CSS (Custom Pink Glassmorphism UI)
- **Backend:** Go 1.24 (`net/http` standard library)
- **Testing:** Vitest, React Testing Library, Go `testing` package
- **Containerization:** Docker, Docker Compose, Nginx

---

## Quickstart (Docker)

The fastest way to run the entire full-stack application is with Docker Compose:

```bash
docker compose up -d --build
```

Once running, access the application in your browser at:
**`http://localhost:8080`**

To stop and remove containers:
```bash
docker compose down
```

---

## Local Development (Without Docker)

### 1. Backend (Go Service)

```bash
cd backend
go run ./cmd/server
```
The backend server listens on `http://localhost:8081` (or `PORT` env variable).

### 2. Frontend (React App)

```bash
cd frontend
npm install
npm run dev
```
The frontend dev server will start at `http://localhost:5173`.

---

## API Documentation

### `POST /api/v1/calculate`

Performs arithmetic calculations.

#### Request Headers
- `Content-Type: application/json`

#### Request Body
```json
{
  "operation": "add",
  "operands": [15, 27]
}
```

#### Supported Operations
| Operation | Operands Required | Description |
| :--- | :---: | :--- |
| `add` | 2 | `a + b` |
| `subtract` | 2 | `a - b` |
| `multiply` | 2 | `a * b` |
| `divide` | 2 | `a / b` |
| `power` | 2 | `a ^ b` |
| `percentage` | 2 | `(a / 100) * b` |
| `sqrt` | 1 | `√a` |

#### Successful Response (`200 OK`)
```json
{
  "operation": "add",
  "operands": [15, 27],
  "result": 42
}
```

#### Error Response (`400 Bad Request` or `422 Unprocessable Entity`)
```json
{
  "error": {
    "code": "DIVISION_BY_ZERO",
    "message": "Cannot divide by zero"
  }
}
```

#### Error Code Reference
- `MISSING_OPERATION` (400): `operation` field omitted.
- `INVALID_OPERATION` (400): Unrecognized operation string.
- `MISSING_OPERANDS` (400): `operands` array omitted.
- `INVALID_OPERAND_COUNT` (400): Operands array length mismatch.
- `INVALID_JSON` (400): Malformed JSON request body.
- `DIVISION_BY_ZERO` (422): Division by 0 attempt.
- `NEGATIVE_SQUARE_ROOT` (422): Square root of a negative number.

---

## Running Unit Tests

### Backend Tests (Go)
```bash
cd backend
go test -v -cover ./...
```

### Frontend Tests (Vitest)
```bash
cd frontend
npm test
```

---

## Design Decisions & Assumptions

1. **Layered Microservice Architecture**: The Go backend is decoupled into distinct model, service, and handler layers. Business logic (`internal/service`) is fully separated from HTTP transport (`internal/handler`), making it lightweight and independently testable without mocking HTTP servers.
2. **HTTP Error Taxonomy**: The API distinguishes between syntactic/validation errors (`400 Bad Request`) and domain logic violations (`422 Unprocessable Entity` for division by zero or negative square roots).
3. **Percentage Operation Convention**: Percentage calculation follows `(a / 100) * b` where the first operand is the percentage rate and the second is the base value.
4. **Unary vs. Binary Operation UI Flow**: Binary operations (`+`, `-`, `*`, `/`, `^`, `%`) wait for a second operand, while unary operations (`sqrt`) execute immediately on the active display value.
5. **Unified Multi-Stage Docker Container**: A single production container uses multi-stage builds to compile the Vite static assets and the Go binary, serving both via Nginx on port `8080` with an internal reverse proxy for `/api/`.

---

## AI Prompts & Workflow Guidance

During development, AI tooling was guided, directed, and corrected through iterative prompt instructions to enforce architectural patterns, resolve edge cases, and ensure clean deployment:

1. **Architecture Setup & API Contract Guidance**:
   - *"Scaffold a full-stack calculator monorepo with React TypeScript in `frontend/` and Go in `backend/`. Define an explicit API contract for `POST /api/v1/calculate` returning distinct `400 Bad Request` for validation errors and `422 Unprocessable Entity` for mathematical domain errors."*

2. **Backend Error Taxonomy Correction**:
   - *"The backend handler shouldn't return generic 500 status codes for bad inputs. Ensure missing operations/operands return `400` with structured error details, while division by zero and negative square root return `422` with clear error messages."*

3. **Frontend State & Unary/Binary UI Flow Guidance**:
   - *"Decouple the calculator UI logic into a custom hook `useCalculator.ts`. Ensure binary operations (`+`, `-`, `*`, `/`, `^`, `%`) store state and wait for the second operand, while unary operations like `sqrt` trigger calculation immediately on the current display value."*

4. **Docker Multi-Stage Build & Entrypoint Correction**:
   - *"Refine the multi-stage Dockerfile to build Vite static assets and the Go server binary in isolation. Configure Nginx to serve static files on port 8080 and proxy `/api/` requests to the Go backend running on port 8081 via `entrypoint.sh`."*

5. **Testing & Final Branding Correction**:
   - *"Write table-driven unit tests in Go covering all 7 operations and validation edge cases, alongside Vitest component tests. Update the HTML page title to 'calculator' and set the favicon to a pink addition (+) icon."*
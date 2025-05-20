# Backend Golang Coding Test

## Objective
Build a simple RESTful API in Golang that manages a list of users. Use MongoDB for persistence, JWT for authentication, and follow clean code practices.

---

## Project Setup

Prerequisites:

Docker & Docker Compose

Environment Variables (in .env file):

MONGO_URI=mongodb://appuser:apppassword@mongodb:27017/?authSource=appdb
MONGO_DB=appdb
MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=rootpassword
MONGO_INITDB_DATABASE=appdb
APP_DB_USER=appuser
APP_DB_PASS=apppassword
JWT_SECRET=backend-challenge

Run the app:

docker-compose up --build

App will be available at: http://localhost:8080

## JWT Token Usage

After login (POST /auth/login), you’ll receive:

{
  "status": "OK",
  "data": {
    "accessToken": "<jwt_token_here>"
  }
}

Use this token in the Authorization header:

Authorization: Bearer <jwt_token_here>

## Sample API Requests

Register

POST /auth/register
{
  "name": "Tee",
  "email": "tee@email.com",
  "password": "12345678"
}

Login

POST /auth/login
{
  "email": "tee@email.com",
  "password": "12345678"
}

Get All Users (Protected)

GET /users
Authorization: Bearer <token>

Get by ID

GET /users/:id
Authorization: Bearer <token>

Update User

PATCH /users/:id
{
  "name": "Updated Name",
  "email": "updated@email.com"
}

Delete User

DELETE /users/:id

## Assumptions / Notes

Password is hashed (not bcrypt in this example)

JWT is short-lived (24h) and HMAC signed

MongoDB initialized via init-mongo.js

Unit tests mock repository for speed and isolation

Logging uses zap with daily file rotation

---

## Requirements

### 1. User Model
Each user should have:
- `ID` (auto-generated)
- `Name` (string)
- `Email` (string, unique)
- `Password` (hashed)
- `CreatedAt` (timestamp)

---

### 2. Authentication

#### Functions
- Register a new user.
- Authenticate user and return a JWT.

#### JWT
- Use JWT for protecting endpoints.
- Use middleware to validate tokens.
- Use HMAC (HS256) with a secret key.

---

### 3. User Functions

- Create a new user.
- Fetch user by ID.
- List all users.
- Update a user's name or email.
- Delete a user.

---

### 4. MongoDB Integration
- Use the official Go MongoDB driver.
- Store and retrieve users from MongoDB.

---

### 5. Middleware
- Logging middleware that logs HTTP method, path, and execution time.

---

### 6. Concurrency Task
- Run a background goroutine every 10 seconds that logs the number of users in the DB.

---

### 7. Testing
Write unit tests

Use Go’s `testing` package. Mock MongoDB where possible.

---

## Bonus (Optional)

- Add Docker + `docker-compose` for API + MongoDB.
- Use Go interfaces to abstract MongoDB operations for testability.
- Add input validation (e.g., required fields, valid email).
- Implement graceful shutdown using `context.Context`.
- **gRPC Version**
  - Create a `.proto` file for `CreateUser` and `GetUser`.
  - Implement a gRPC server.
  - (Optional) Secure gRPC with token metadata.
- **Hexagonal Architecture**
  - Structure the project using hexagonal (ports & adapters) architecture:
    - Separate domain, application, and infrastructure layers.
    - Use interfaces for data access and external dependencies.
    - Keep business logic decoupled from frameworks and DB drivers.

---

## Submission Guidelines

- Submit a GitHub repo or zip file.
- Include a `README.md` with:
  - Project setup and run instructions
  - JWT token usage guide
  - Sample API requests/responses
  - Any assumptions or decisions made

---

## Evaluation Criteria

- Code quality, structure, and readability
- REST API correctness and completeness
- JWT implementation and security
- MongoDB usage and abstraction
- Bonus: gRPC, Docker, validation, shutdown
- Testing coverage and mocking
- Use of idiomatic Go
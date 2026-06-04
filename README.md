# GoTickets

GoTickets is a beginner-friendly REST API project written in Go. It shows how to build a small ticket booking backend with user authentication, events, bookings, PostgreSQL, GORM, Echo, validation, and JWT.

This project is useful for students who are learning how a Go web application is usually organized.

## What You Will Learn

- How to create a web server with Echo
- How to organize code into packages
- How to connect Go with PostgreSQL using GORM
- How to create users and login with JWT authentication
- How to protect routes with middleware
- How to use DTOs for request and response data
- How service, repository, and handler layers work together

## Tech Stack

- Go
- Echo for HTTP routing
- GORM for database access
- PostgreSQL as the database
- JWT for authentication
- Validator for request validation
- godotenv for loading `.env` files

## Project Structure

```text
gotickets/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── auth/                   # JWT token create/validate logic
│   ├── booking/                # Booking feature
│   ├── config/                 # Environment and database config
│   ├── event/                  # Event feature
│   ├── httpresponse/           # Common error response shape
│   ├── middlewares/            # Auth middleware
│   ├── server/                 # Echo server setup
│   └── user/                   # User register/login feature
├── .env.example                # Example environment variables
├── go.mod                      # Go module and dependencies
└── go.sum
```

## How The Code Flows

For most features, the code follows this pattern:

```text
Route -> Handler -> Service -> Repository -> Database
```

- **Route** decides the URL and HTTP method.
- **Handler** reads the request and returns the response.
- **Service** contains the main business logic.
- **Repository** talks to the database.
- **DTO** defines request and response data shapes.

Example:

```text
POST /api/v1/bookings
    -> booking handler
    -> booking service
    -> booking repository
    -> PostgreSQL
```

## Requirements

Before running the project, install:

- Go
- PostgreSQL
- Git

You can check your Go version with:

```bash
go version
```

## Setup

1. Clone the project:

```bash
git clone https://github.com/Apollo-Level2-Web-Dev/gotickets.git
cd gotickets
```

2. Create a PostgreSQL database:

```sql
CREATE DATABASE gotickets;
```

3. Create your `.env` file:

```bash
cp .env.example .env
```

4. Update `.env` with your own database values:

```env
DSN="host=localhost user=postgres password=postgres dbname=gotickets port=5432 sslmode=disable TimeZone=Asia/Dhaka"
PORT=8080
JWT_SECRET=change-this-secret
```

The `DSN` value tells GORM how to connect to PostgreSQL.

## Run The Project

Install dependencies:

```bash
go mod tidy
```

Start the server:

```bash
go run cmd/main.go
```

If everything is okay, the server will start on:

```text
http://localhost:8080
```

Check the health route:

```bash
curl http://localhost:8080/health
```

Expected response:

```text
ok
```

## API Routes

### Auth Routes

Register a new user:

```http
POST /api/v1/auth/register
```

Example:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mahi",
    "email": "mahi@example.com",
    "password": "secret123"
  }'
```

Login:

```http
POST /api/v1/auth/login
```

Example:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mahi@example.com",
    "password": "secret123"
  }'
```

The login response returns a JWT token. Keep that token for protected routes.

Get current user:

```http
GET /api/v1/auth/me
```

Example:

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Event Routes

Create an event:

```http
POST /api/v1/events
```

Example:

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Go Workshop",
    "description": "A beginner friendly Go learning event",
    "location": "Dhaka",
    "starts_at": "2026-07-01T10:00:00Z",
    "total_tickets": 100,
    "price": 500
  }'
```

Get all events:

```http
GET /api/v1/events
```

Example:

```bash
curl http://localhost:8080/api/v1/events
```

Get one event:

```http
GET /api/v1/events/:id
```

Example:

```bash
curl http://localhost:8080/api/v1/events/1
```

Update an event:

```http
PATCH /api/v1/events/:id
```

Example:

```bash
curl -X PATCH http://localhost:8080/api/v1/events/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Go Workshop",
    "price": 700
  }'
```

### Booking Routes

Booking routes are protected. You must login first and send the JWT token in the `Authorization` header.

Create a booking:

```http
POST /api/v1/bookings
```

Example:

```bash
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "event_id": 1,
    "quantity": 2
  }'
```

Get my bookings:

```http
GET /api/v1/bookings/me
```

Example:

```bash
curl http://localhost:8080/api/v1/bookings/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Important Concepts For Beginners

### `cmd/main.go`

This is where the application starts. It loads environment variables, connects to the database, and starts the server.

### `internal/config`

This package reads values from `.env` and creates the database connection.

### `internal/server`

This package creates the Echo app, adds middleware, registers routes, and starts listening on the selected port.

### Handler

A handler receives an HTTP request. It usually does these things:

1. Reads JSON input
2. Validates the input
3. Calls the service
4. Sends JSON response

### Service

A service contains business logic. For example, the booking service checks if enough tickets are available before creating a booking.

### Repository

A repository is responsible for database queries. This keeps database code separate from business logic.

### Middleware

Middleware runs before the handler. In this project, auth middleware checks the JWT token and adds user information to the request context.

## Database Tables

The project uses GORM `AutoMigrate`, so tables are created automatically when the server starts.

Current main tables:

- users
- events
- bookings

## Common Problems

### `.env` file not found

Make sure you created a `.env` file:

```bash
cp .env.example .env
```

### Database connection failed

Check these things:

- PostgreSQL is running
- Database name is correct
- Username and password are correct
- Port is correct, usually `5432`

### Protected route says unauthorized

Make sure the request has this header:

```text
Authorization: Bearer YOUR_TOKEN_HERE
```

Also make sure the token comes from the login route.

## Suggested Learning Path

If you are new to Go, read the project in this order:

1. `cmd/main.go`
2. `internal/config/config.go`
3. `internal/config/db.go`
4. `internal/server/http.go`
5. `internal/user/register.go`
6. `internal/user/handler.go`
7. `internal/user/service.go`
8. `internal/user/repository.go`
9. `internal/event`
10. `internal/booking`

This order helps you understand how the app starts, then how one feature works from route to database.

## Notes

This project is for learning. Some ideas you can add later:

- Delete event route
- Cancel booking route
- Admin-only event routes
- Unit tests
- Docker setup
- Better error messages
- Pagination for events

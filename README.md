# Chirpy Project

A RESTful API backend for a Twitter-like social media platform built with Go. Users can create accounts, post chirps (short messages), and manage their content with secure authentication.

## Features

- **User Management**: User registration, login, and profile updates
- **Chirp Management**: Create, read, and delete chirps
- **Authentication**: JWT-based authentication with refresh tokens
- **Security**: Argon2id password hashing
- **Database**: PostgreSQL for data persistence
- **Webhooks**: Polka webhook integration for premium user upgrades

## Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Migrations**: Goose
- **Query Layer**: SQLC

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- `.env` file with required configuration

## Installation

1. Clone the repository:
```bash
git clone https://github.com/anxhukumar/chirpy-project.git
cd chirpy-project
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your PostgreSQL database and create a `.env` file in the root directory:
```env
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your-secret-key-here
POLKA_KEY=your-polka-api-key-here
```

## Running the Application

Start the server:
```bash
go run main.go
```

The server will start on port 8080

## API Endpoints

### Readiness Endpoint
- `GET /api/healthz` - A simple health endpoint that returns `200 OK` to indicate the service is ready to receive traffic.

### User Management
- `POST /api/users` - Create a new user
- `PUT /api/users` - Update user information (requires authentication)

### Authentication
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token

### Chirps
- **POST `/api/chirps`** – Create a new chirp (requires authentication)
- **GET `/api/chirps`** – Get all chirps  
  - Optional: filter by author using `author_id` query param  
  - Optional: sort results using `sort=asc` or `sort=desc`
- **GET `/api/chirps/{chirpID}`** – Get a specific chirp
- **DELETE `/api/chirps/{chirpID}`** – Delete a chirp (requires authentication)

### Webhooks
- `POST /api/polka/webhooks` - Polka webhook for premium upgrades 

### Admin
- `GET /admin/metrics` - View metrics 
- `POST /admin/reset` - Reset metrics and database

## Project Structure

```
chirpy-project/
├── main.go              # Application entry point
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
├── index.html           # Frontend HTML file
└── internal/
    ├── api/             # API handlers
    ├── auth/            # Authentication logic
    ├── database/        # Database queries
    └── helper/          # Helper utilities
```

## Environment Variables

The application requires the following environment variables:

- `DB_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT token generation
- `POLKA_KEY`: API key for Polka webhook authentication

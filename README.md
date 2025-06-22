# Utils Application

A full-stack application built with Go (Gin), Next.js, and PostgreSQL.

## Architecture

- **Backend**: Go with Gin framework (Clean Architecture/DDD)
- **Frontend**: Next.js with TypeScript and Tailwind CSS
- **Database**: PostgreSQL
- **Containerization**: Docker & Docker Compose

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Node.js (for local development)
- Go 1.23+ (for local development)

### Running with Docker Compose

1. Clone the repository
2. Copy environment files:
   ```bash
   cp .env.example .env
   cp frontend/.env.local.example frontend/.env.local
   ```
3. Start all services:
   ```bash
   docker-compose up --build
   ```

The application will be available at:
- Frontend: http://localhost:3000
- API: http://localhost:8080
- Database: localhost:5432

### Local Development

#### Backend (API)
```bash
cd api
go mod download
go run cmd/server/main.go
```

#### Frontend
```bash
cd frontend
npm install
npm run dev
```

## Project Structure

```
utils/
├── api/                    # Go backend
│   ├── cmd/
│   ├── internal/
│   └── pkg/
├── frontend/               # Next.js frontend
│   ├── app/
│   ├── lib/
│   └── components/
├── db/                     # Database initialization
│   └── init/
├── docs/                   # Documentation
│   └── dev_diary/
└── docker-compose.yml
```

## API Endpoints

- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID
- `POST /api/users` - Create new user

## Development Notes

See `docs/dev_diary/` for detailed development logs and progress tracking.
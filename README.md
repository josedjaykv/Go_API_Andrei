# Andrei Mes Manur API - Dockerized

Elaborado por Jose David Jayk Vanegas y Diego Collazos Bermudez
Backend API for the Andrei Mes Manur application - a role-based system where a failed warlock named Andrei controls demons who capture network administrators.

## Features

- JWT Authentication
- Role-based authorization (Andrei, Demon, Network Admin)
- CRUD operations for users, posts, reports, and rewards
- Statistics and rankings
- Public resistance page
- Anonymous posting for Network Admins
- **Full Docker support with PostgreSQL and Nginx frontend**

## 🚀 Quick Start with Docker

### Prerequisites
- Docker
- Docker Compose

### Run the Complete System
```bash
# Build and run all services (backend, frontend, database)
docker-compose up --build

# Or in background mode
docker-compose up -d --build
```

### Services Available
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8085
- **Database**: localhost:5432

### Default Users
| Email | Password | Role |
|-------|----------|------|
| andrei@evil.com | password123 | andrei |
| demon1@evil.com | password123 | demon |
| admin1@network.com | password123 | network_admin |
| admin2@network.com | password123 | network_admin |
| admin3@network.com | password123 | network_admin |

## 🐳 Docker Commands

```bash
# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild specific service
docker-compose up --build backend

# Access database
docker-compose exec postgres psql -U postgres -d andrei_db
```

## 🛠 Manual Setup (Without Docker)

### Prerequisites

- Go 1.19 or higher
- PostgreSQL database

### Database Setup

```bash
# Create PostgreSQL container
docker run --name postgres_local -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres123 -e POSTGRES_DB=andrei_db -p 5432:5432 -d postgres:15
```

### Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod tidy
```

3. Set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres123
export DB_NAME=andrei_db
export JWT_SECRET=your-secret-key
export PORT=8085
```

4. Seed the database:
```bash
go run main.go -seed
```

5. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8085`

## API Endpoints

### Authentication

#### Register
- **POST** `/api/v1/register`
- Register a new user (Demon or Network Admin only)
- Body:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "role": "demon"
}
```

#### Login
- **POST** `/api/v1/login`
- Login with email and password
- Body:
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```
- Returns JWT token

### Public Endpoints

#### Resistance Page
- **GET** `/api/v1/resistance`
- Get all posts for the resistance page (public access)

### Andrei (Admin) Endpoints

**Authentication Required:** Bearer token with Andrei role

#### User Management
- **GET** `/api/v1/admin/users` - Get all users
- **GET** `/api/v1/admin/users/:id` - Get user by ID
- **DELETE** `/api/v1/admin/users/:id` - Delete user

#### Rewards & Punishments
- **POST** `/api/v1/admin/rewards` - Create reward/punishment for demon
- Body:
```json
{
  "demon_id": 1,
  "type": "reward",
  "title": "Good work",
  "description": "Captured 5 admins",
  "points": 100
}
```

#### Statistics
- **GET** `/api/v1/admin/stats` - Get platform statistics
- **GET** `/api/v1/admin/demons/ranking` - Get demon rankings

#### Posts Management
- **GET** `/api/v1/admin/posts` - Get all posts
- **POST** `/api/v1/admin/posts` - Create new post
- **DELETE** `/api/v1/admin/posts/:id` - Delete post

### Demon Endpoints

**Authentication Required:** Bearer token with Demon role

#### Victim Management
- **POST** `/api/v1/demons/victims` - Register new victim (Network Admin)
- Body:
```json
{
  "username": "victim1",
  "email": "victim@example.com",
  "password": "password123",
  "role": "network_admin"
}
```

- **GET** `/api/v1/demons/victims` - Get my victims

#### Reports
- **POST** `/api/v1/demons/reports` - Create report about victim
- Body:
```json
{
  "victim_id": 1,
  "title": "Victim Status",
  "description": "Successfully hypnotized"
}
```

- **GET** `/api/v1/demons/reports` - Get my reports
- **PUT** `/api/v1/demons/reports/:id` - Update report status

#### Statistics
- **GET** `/api/v1/demons/stats` - Get my personal statistics

#### Posts
- **POST** `/api/v1/demons/posts` - Create new post
- Body:
```json
{
  "title": "Survival Tip",
  "body": "How to avoid detection",
  "media": "optional-media-url"
}
```

### Network Admin Endpoints

**Authentication Required:** Bearer token with Network Admin role

#### Anonymous Posts
- **POST** `/api/v1/network-admins/posts/anonymous` - Create anonymous post
- Body:
```json
{
  "title": "Resistance Message",
  "body": "Fight back!",
  "media": "optional-media-url"
}
```

## User Roles

### Andrei (Supreme Leader)
- Full platform access
- CRUD operations on all entities
- Can assign rewards/punishments to demons
- View all statistics and rankings
- Create posts

### Demon
- Register new victims (Network Admins)
- Create reports about victims
- View personal statistics
- Create posts for resistance page
- Update own report statuses

### Network Admin
- Access resistance page
- Create anonymous posts
- Limited access (victims of the system)

## Default Credentials

After running the seed script:
- **Email:** andrei@evil.com
- **Password:** password123
- **Role:** Andrei

## Error Responses

All endpoints return errors in the following format:
```json
{
  "error": "Error message description"
}
```

## Security

- JWT tokens expire after 24 hours
- Passwords are hashed using bcrypt
- Role-based authorization on all protected endpoints
- CORS enabled for frontend integration

## Development

To run in development mode with hot reload:
```bash
go install github.com/githubnemo/CompileDaemon@latest
CompileDaemon -command="./andrei-api"
```

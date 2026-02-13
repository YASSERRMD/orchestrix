# FlowForge - Workflow & Approvals Engine

A full-stack workflow and approvals engine with strict state transitions, role-based access, audit logs, and real-time updates.

## Tech Stack

- **Frontend**: Vue 3 + Vite
- **Backend**: Go + Gin
- **Database**: SQLite (raw SQL)
- **Real-time**: Gorilla WebSocket
- **Containerization**: Docker

## Project Structure

```
flowforge/
├── backend/
│   ├── cmd/server/         # Main application entry
│   ├── internal/
│   │   ├── config/         # Configuration
│   │   ├── database/       # Database initialization
│   │   ├── handlers/       # HTTP handlers
│   │   ├── middleware/     # Auth & CORS middleware
│   │   ├── models/        # Data models
│   │   ├── repository/    # Raw SQL repositories
│   │   ├── services/     # Business logic
│   │   ├── websocket/     # WebSocket hub
│   │   └── worker/       # Background worker
│   ├── migrations/        # SQL migrations
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── views/        # Login, Dashboard, Templates, WorkflowDetail
│   │   ├── router/       # Vue Router
│   │   ├── services/     # API & WebSocket services
│   │   └── App.vue
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml
└── SPEC.md
```

## Getting Started

### Option 1: Docker Compose (Recommended)

```bash
# Build and start all services
docker-compose up --build

# Access the application
# Frontend: http://localhost:5173
# Backend API: http://localhost:8080
```

### Option 2: Local Development

#### Backend

```bash
cd backend

# Install dependencies
go mod download

# Run the server
go run ./cmd/server
```

#### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

## Default Users

After seeding the database:

| Email | Password | Role |
|-------|----------|------|
| admin@demo.com | password123 | Admin |
| manager@demo.com | password123 | Manager |
| operator@demo.com | password123 | Operator |
| viewer@demo.com | password123 | Viewer |

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register organization + admin
- `POST /api/auth/login` - Login
- `GET /api/auth/me` - Current user

### Templates
- `GET /api/templates` - List templates
- `POST /api/templates` - Create template
- `GET /api/templates/:id` - Get template
- `GET /api/templates/:id/versions` - List versions
- `PUT /api/templates/:id` - Update (new version)

### Workflows
- `GET /api/workflows` - List workflows
- `POST /api/workflows` - Create workflow
- `GET /api/workflows/:id` - Get workflow
- `POST /api/workflows/:id/advance` - Advance
- `POST /api/workflows/:id/reject` - Reject
- `POST /api/workflows/:id/rollback` - Rollback (admin/manager)
- `POST /api/workflows/:id/cancel` - Cancel
- `GET /api/workflows/:id/audit` - Audit logs

### WebSocket
- `WS /api/ws` - Real-time events

## Features

1. **Multi-tenant RBAC** - Organizations with Admin, Manager, Operator, Viewer roles
2. **Versioned Templates** - Updating templates creates new versions
3. **State Machine** - Strict transitions with audit logging
4. **Background Worker** - Handles timeouts and escalations
5. **Real-time Updates** - WebSocket events for workflow changes

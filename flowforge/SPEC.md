# FlowForge - Workflow & Approvals Engine

## Project Overview
- **Project Name**: FlowForge
- **Type**: Full-stack workflow approval system
- **Core Functionality**: Multi-tenant workflow engine with strict state transitions, role-based access, and real-time updates
- **Target Users**: Organizations needing approval workflows with audit trails

## Technology Stack
- **Frontend**: Vue 3 + Vite + Vue Router
- **Backend**: Go + Gin
- **Database**: SQLite (raw SQL)
- **Real-time**: Gorilla WebSocket
- **Containerization**: Docker + Docker Compose

---

## 1. Database Schema

### Tables

#### organizations
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY |
| name | TEXT | NOT NULL |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP |
| updated_at | DATETIME | DEFAULT CURRENT_TIMESTAMP |

#### users
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY |
| organization_id | INTEGER | FK -> organizations.id |
| email | TEXT | UNIQUE, NOT NULL |
| password_hash | TEXT | NOT NULL |
| name | TEXT | NOT NULL |
| role | TEXT | NOT NULL (admin/manager/operator/viewer) |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP |
| updated_at | DATETIME | DEFAULT CURRENT_TIMESTAMP |

#### workflow_templates
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY |
| organization_id | INTEGER | FK -> organizations.id |
| name | TEXT | NOT NULL |
| description | TEXT | |
| version | INTEGER | NOT NULL DEFAULT 1 |
| is_active | BOOLEAN | DEFAULT 1 |
| created_by | INTEGER | FK -> users.id |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP |
| updated_at | DATETIME | DEFAULT CURRENT_TIMESTAMP |

#### template_stages
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY |
| template_id | INTEGER | FK -> workflow_templates.id |
| name | TEXT | NOT NULL |
| order_index | INTEGER | NOT NULL |
| required_role | TEXT | NOT NULL |
| approval_type | TEXT | NOT NULL (manual/auto) |
| timeout_minutes | INTEGER | DEFAULT 0 |
| escalation_role | TEXT | |

#### workflow_instances
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY |
| organization_id | INTEGER | FK -> organizations.id |
| template_id | INTEGER | FK -> workflow_templates.id |
| template_version | INTEGER | NOT NULL |
| title | TEXT | NOT NULL |
| current_stage_id | INTEGER | FK -> template_stages.id |
| status | TEXT | NOT NULL (pending/approved/rejected/escalated/canceled/completed) |
| created_by | INTEGER | FK -> users.id |
| assigned_to | INTEGER | FK -> users.id |
| escalated_to | INTEGER | FK -> users.id |
| started_at | DATETIME | DEFAULT CURRENT_TIMESTAMP |
| completed_at | DATETIME | |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP |
| updated_at | DATETIME | DEFAULT CURRENT_TIMESTAMP |

#### audit_logs
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY |
| workflow_id | INTEGER | FK -> workflow_instances.id |
| actor_id | INTEGER | FK -> users.id |
| actor_role | TEXT | NOT NULL |
| from_stage_id | INTEGER | FK -> template_stages.id |
| to_stage_id | INTEGER | FK -> template_stages.id |
| action | TEXT | NOT NULL |
| reason | TEXT | |
| timestamp | DATETIME | DEFAULT CURRENT_TIMESTAMP |

#### locks (for worker concurrency)
| Column | Type | Constraints |
|--------|------|-------------|
| key | TEXT | PRIMARY KEY |
| owner | TEXT | NOT NULL |
| expires_at | DATETIME | NOT NULL |

---

## 2. API Endpoints

### Authentication
- `POST /api/auth/register` - Register new organization + admin
- `POST /api/auth/login` - Login and get JWT
- `GET /api/auth/me` - Get current user

### Users
- `GET /api/users` - List users (admin/manager)
- `POST /api/users` - Create user (admin)
- `GET /api/users/:id` - Get user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Templates
- `GET /api/templates` - List templates
- `POST /api/templates` - Create template
- `GET /api/templates/:id` - Get template
- `GET /api/templates/:id/versions` - List versions
- `PUT /api/templates/:id` - Update template (creates new version)

### Workflows
- `GET /api/workflows` - List workflows (filtered)
- `POST /api/workflows` - Create workflow instance
- `GET /api/workflows/:id` - Get workflow detail
- `POST /api/workflows/:id/advance` - Advance workflow
- `POST /api/workflows/:id/reject` - Reject workflow
- `POST /api/workflows/:id/rollback` - Rollback (manager/admin)
- `POST /api/workflows/:id/cancel` - Cancel workflow

### Audit
- `GET /api/workflows/:id/audit` - Get audit logs for workflow
- `GET /api/audit` - Global audit logs (admin)

---

## 3. State Machine

### Workflow Statuses
```
pending -> approved -> completed
pending -> rejected
pending -> escalated
escalated -> approved -> completed
escalated -> rejected
canceled (from any state)
```

### Stage Transitions
- **advance**: Move to next stage (requires role >= required_role)
- **reject**: Reject entire workflow (requires role >= required_role)
- **rollback**: Go back to previous stage (manager/admin only)
- **cancel**: Cancel workflow (creator/admin/manager)

### Auto-Approval
- If approval_type = "auto", automatically advance after timeout

---

## 4. WebSocket Events

### Server -> Client
- `workflow_created` - New workflow instance
- `stage_advanced` - Workflow moved forward
- `rejected` - Workflow was rejected
- `rolled_back` - Workflow was rolled back
- `escalated` - Workflow was escalated (timeout)
- `canceled` - Workflow was canceled

### Client -> Server
- Subscribe to organization channel
- Subscribe to user-specific events

---

## 5. Background Worker

### Functions
1. Check for stage timeouts every 30 seconds
2. For each timed-out stage:
   - Acquire lock (prevent duplicates)
   - If escalation_role exists: reassign and escalate
   - If auto-approval: advance automatically
   - Create audit log
   - Emit WebSocket event
3. Release lock after processing

---

## 6. Frontend Pages

### Login Page
- Email/password form
- JWT token storage
- Redirect to dashboard

### Dashboard
- Tabs: Active | Completed | Escalated | My Approvals
- Workflow cards with status badges
- Real-time updates via WebSocket

### Template Management
- List templates with versions
- Create/Edit template with stage builder
- Stage configuration form

### Workflow Detail
- Current stage display
- Transition timeline (audit log)
- Action buttons (role-based)
- Full history

---

## 7. Folder Structure

```
/flowforge
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── middleware/
│   │   ├── handlers/
│   │   ├── repository/
│   │   ├── services/
│   │   ├── websocket/
│   │   ├── worker/
│   │   └── models/
│   ├── migrations/
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   ├── router/
│   │   ├── services/
│   │   ├── stores/
│   │   └── App.vue
│   ├── Dockerfile
│   └── docker-compose.yml
└── docker-compose.yml
```

---

## 8. Git Branch Plan

| Phase | Branch | Description |
|-------|--------|-------------|
| 1 | phase-1-setup | Project structure, Docker, migrations |
| 2 | phase-2-auth | Auth, users, RBAC middleware |
| 3 | phase-3-templates | Template CRUD, versioning |
| 4 | phase-4-workflows | Instances, state machine, audit |
| 5 | phase-5-worker | Background worker, timeouts |
| 6 | phase-6-websocket | WebSocket hub, events |
| 7 | phase-7-frontend | Vue 3 frontend |
| 8 | phase-8-integration | Docker compose, testing |

---

## 9. Commit Breakdown

### Phase 1 (Setup)
1. Create project folders
2. Backend main.go + basic Gin setup
3. SQLite migration scripts
4. Docker files
5. docker-compose

### Phase 2 (Auth)
1. Database queries for users/orgs
2. Auth handlers + JWT
3. Role middleware
4. User CRUD

### Phase 3 (Templates)
1. Template repository
2. Template handlers
3. Stage management
4. Versioning logic

### Phase 4 (Workflows)
1. Instance repository
2. State machine service
3. Audit logging
4. Workflow handlers

### Phase 5 (Worker)
1. Lock mechanism
2. Timeout checker
3. Escalation logic

### Phase 6 (WebSocket)
1. Hub implementation
2. Event broadcasting
3. Client handlers

### Phase 7 (Frontend)
1. Vue project setup
2. Router + API service
3. Login page
4. Dashboard
5. Template management
6. Workflow detail

### Phase 8 (Integration)
1. Docker compose
2. Seed data
3. Testing
4. Final polish

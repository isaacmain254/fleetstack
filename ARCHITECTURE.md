# FleetStack Project Structure

> A proposed monorepo architecture for **FleetStack**, designed for long-term development as an open-source deployment platform.

---

# Project Goals

FleetStack is intended to be:

- A long-term open-source project
- Easy for new contributors to understand
- Scalable without major refactoring
- Organized around **responsibility**, not just programming language or framework

The architecture separates the system into independent components that each have a clear purpose.

---

# Repository Structure

```text
FleetStack/
│
├── frontend/                 # Vue dashboard
│   ├── src/
│   ├── public/
│   ├── Dockerfile
│   └── package.json
│
├── backend/                  # Go API (Control Plane)
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── auth/
│   │   ├── deployments/
│   │   ├── docker/
│   │   ├── workers/
│   │   ├── scheduler/
│   │   ├── projects/
│   │   ├── users/
│   │   ├── database/
│   │   └── api/
│   │
│   ├── pkg/
│   ├── migrations/
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── worker/                   # Go Worker Agent
│   ├── cmd/
│   │   └── worker/
│   ├── internal/
│   ├── Dockerfile
│   └── go.mod
│
├── cli/                      # Optional CLI (future)
│
├── docs/                     # Documentation
│
├── docker/
│   ├── traefik/
│   ├── postgres/
│   └── redis/
│
├── scripts/
│
├── docker-compose.yml
├── .env.example
├── README.md
└── LICENSE
```

---

# Architecture Overview

```text
               +----------------------+
               |      Frontend        |
               |   Vue Dashboard UI   |
               +----------+-----------+
                          |
                          |
                    REST / WebSocket
                          |
                          ▼
               +----------------------+
               |      Backend         |
               |    Control Plane     |
               +----------+-----------+
                          |
                  Job Queue / API
                          |
                          ▼
               +----------------------+
               |       Worker         |
               |  Executes Deployments|
               +----------+-----------+
                          |
                    Docker Engine
                          |
                          ▼
                   Containers / Apps
```

---

# Why Separate Backend and Worker?

Although both applications are written in Go, they have completely different responsibilities.

Separating them improves scalability, testing, and maintainability.

## Backend (Control Plane)

The backend manages the platform.

It **does not** communicate directly with Docker running on remote servers.

### Responsibilities

- Authentication
- User management
- Project management
- Deployment records
- GitHub webhooks
- Worker registration
- Scheduling
- REST API
- Dashboard communication

Think of it as the **brain** of the platform.

---

## Worker (Execution Engine)

The worker is installed on a VPS or server.

Its job is to execute deployment tasks requested by the backend.

### Responsibilities

- Build Docker images
- Pull repositories
- Start containers
- Stop containers
- Restart services
- Stream logs
- Collect CPU and memory metrics
- Perform health checks
- Report deployment status

Think of it as the **hands** of the platform.

---

# Backend Structure

Instead of organizing code into generic folders such as:

```text
controllers/
services/
models/
```

Organize by feature.

```text
internal/
├── auth/
├── projects/
├── deployments/
├── workers/
├── scheduler/
├── docker/
├── database/
└── api/
```

Each feature contains everything it needs.

Example:

```text
internal/
└── deployments/
    ├── handler.go
    ├── service.go
    ├── repository.go
    └── model.go
```

### Benefits

- Easier navigation
- Better modularity
- Reduced coupling
- Self-contained features

---

# Frontend Structure

The frontend should also be organized by feature rather than by component type.

```text
frontend/src/
├── components/
├── features/
├── pages/
├── router/
├── stores/
├── services/
├── types/
└── utils/
```

Inside the `features` directory:

```text
features/
├── projects/
├── deployments/
├── workers/
├── settings/
├── users/
└── logs/
```

This structure scales significantly better as the application grows.

---

# Docker Organization

Infrastructure should live in its own directory instead of cluttering the project root.

```text
docker/
├── postgres/
├── redis/
└── traefik/
```

### Benefits

- Cleaner repository
- Easier infrastructure management
- Better separation of concerns

---

# Documentation

Documentation is one of the most valuable parts of an open-source project.

Suggested structure:

```text
docs/
├── architecture.md
├── worker.md
├── deployment.md
├── api.md
└── security.md
```

Include diagrams wherever possible to help contributors understand the system quickly.

---

# Future Expansion

As FleetStack evolves, additional directories can be introduced without restructuring the project.

## SDKs

```text
sdk/
├── typescript/
└── go/
```

## Plugins

```text
plugins/
├── github/
├── gitlab/
└── bitbucket/
```

These are not required for the initial version and should be added only when needed.

---

# Minimal First Version

To keep development focused, start with only the essential directories:

```text
FleetStack/
├── frontend/
├── backend/
├── worker/
├── docker/
├── docs/
├── docker-compose.yml
└── README.md
```

This provides a solid foundation while avoiding unnecessary complexity.

---

# Design Philosophy

The guiding principle for the project should be:

> **Organize by responsibility, not by technology.**

This means:

- **Frontend** handles everything related to the web interface.
- **Backend** acts as the control plane, managing users, projects, deployments, and workers.
- **Worker** is the execution agent installed on servers that performs deployment tasks and interacts with Docker.

This separation closely reflects the real architecture of the platform, making the codebase easier to understand, maintain, test, and extend as the project grows.

---

# Summary

By adopting this structure from the beginning, FleetStack gains:

- Clear separation of concerns
- Better scalability
- Easier onboarding for contributors
- Feature-based organization
- Clean infrastructure management
- A maintainable architecture suitable for long-term open-source development

The initial version remains intentionally simple while leaving room for future growth without requiring major refactoring.
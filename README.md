# Task Management Full-Stack (2025) — Go + React

A restored educational task-management application with a Go backend, React-based browser UI, and an optional GORM/PostgreSQL mode.

> This is a **2026 reconstruction** of the 2025 team assignment, **not an exact copy** of a teammate's original source code. Preserve attribution for earlier group work.

## Quick start (Windows, without installing Go or Node.js)

1. Download this repository as ZIP using **Code → Download ZIP**.
2. Extract the full archive.
3. Open the folder `Task_Management_FullStack_2025`.
4. Double-click `START_WINDOWS.bat`.
5. Open **http://localhost:8082**.

The fast mode uses a bundled Windows binary and local JSON storage. To use **GORM + PostgreSQL**, install Docker Desktop and use `START_DOCKER_WINDOWS.bat`.

## Features

- Create, edit, delete and search tasks
- Statuses, priorities, deadlines and Kanban board
- Task statistics, REST API and CSV export
- Persist data between restarts

## Source code

- [Full documentation](Task_Management_FullStack_2025/README.md)
- [Russian setup](Task_Management_FullStack_2025/README_RU.md)
- [Backend](Task_Management_FullStack_2025/backend/)
- [Frontend](Task_Management_FullStack_2025/frontend/)
- [Compose configuration](Task_Management_FullStack_2025/compose.yaml)

## Build and verify

The backend module is in `Task_Management_FullStack_2025/backend`:

```bash
cd Task_Management_FullStack_2025/backend
go test ./...
go build -o TaskManagement .
```

CI runs these Go checks automatically. Tests confirm backend behavior covered by test cases; **production PostgreSQL, Docker deployment and browser E2E are not covered by those checks**.

## Safety

The bundled binary is only for Windows. Keep local databases, personal tasks and credentials out of the repository. Don't publish secrets or private task content.

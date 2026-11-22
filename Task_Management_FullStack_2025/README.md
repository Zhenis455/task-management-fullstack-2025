# Task Management Full-Stack (2025) — Restored Edition

A clean-room reconstruction of a 2025 Go + React/GORM course project. Public team references: [backend](https://github.com/Bekzhanizb/TaskManagerBackend) and [frontend](https://github.com/Bekzhanizb/TaskManagerFrontend). Original source code was not copied.

**Windows (no prerequisites):** unzip, double-click `START_WINDOWS.bat`, browse to `http://localhost:8082`. Portable mode persists data to `backend/data/tasks.json` (does not use PostgreSQL).

**PostgreSQL/GORM mode:** start Docker Desktop and run `docker compose up --build`, then open the same URL. Stop the portable process first because they share a port. Docker first-run downloads Go modules and PostgreSQL image. Data is persisted in Docker volume `pg_data`.

Read `README_RU.md` for complete setup, GitHub commands, API reference and verification status.

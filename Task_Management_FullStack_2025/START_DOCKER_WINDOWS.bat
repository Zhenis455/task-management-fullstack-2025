@echo off
cd /d "%~dp0"
where docker >nul 2>nul
if errorlevel 1 (
  echo Docker was not found. Install and start Docker Desktop first.
  echo For the portable version run START_WINDOWS.bat instead.
  pause
  exit /b 1
)
echo Starting PostgreSQL + GORM mode on http://localhost:8082
echo Close portable START_WINDOWS.bat first if it is running.
docker compose up --build
if errorlevel 1 pause

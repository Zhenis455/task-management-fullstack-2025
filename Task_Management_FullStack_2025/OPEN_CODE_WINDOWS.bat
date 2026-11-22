@echo off
cd /d "%~dp0"
where code >nul 2>nul
if errorlevel 1 (
  echo VS Code command is unavailable. Opening project folder instead.
  start "" "%~dp0"
) else (
  code .
)

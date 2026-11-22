@echo off
setlocal
cd /d "%~dp0backend"
where go >nul 2>nul
if errorlevel 1 (
  echo Go SDK not found. Please install Go 1.23 or newer to rebuild from source.
  echo Prebuilt TaskManagement.exe can run without Go.
  pause
  exit /b 1
)
echo Building portable Windows TaskManagement.exe ...
go build -o ..\TaskManagement.exe .
if errorlevel 1 (
  echo Build failed. Review the error above.
) else (
  echo Build completed. You can run START_WINDOWS.bat
)
pause

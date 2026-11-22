@echo off
cd /d "%~dp0"
echo API smoke test. Make sure START_WINDOWS.bat is running.
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0scripts\verify-api.ps1"
if errorlevel 1 (
  echo API test failed.
) else (
  echo API test passed.
)
pause

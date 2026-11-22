@echo off
setlocal
cd /d "%~dp0"
title Task Management Full-Stack 2025
color 0B
echo =========================================================
echo   TASKFLOW - GO + REACT FULL-STACK (PORTABLE)
echo =========================================================
echo.
echo   Open in Chrome: http://localhost:8082
 echo  Tasks:          http://localhost:8082/tasks
 echo  Kanban board:   http://localhost:8082/board
 echo  REST API:       http://localhost:8082/api/tasks
 echo.
echo   Stop: Ctrl+C
 echo  Tasks are saved in backend\data\tasks.json
 echo.
if not exist "TaskManagement.exe" (
  echo ERROR: TaskManagement.exe was not found in this folder.
  echo Unzip the whole archive and run this BAT inside the folder.
  pause
  exit /b 1
)
set "TASKS_FILE=%CD%\backend\data\tasks.json"
"%~dp0TaskManagement.exe"
if errorlevel 1 (
  echo.
  echo Server stopped with an error. Check whether port 8082 is already used.
  pause
)
endlocal

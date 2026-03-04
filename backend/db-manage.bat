@echo off
REM Firebase Database Management Script for Windows
REM Automates database initialization and management tasks

setlocal enabledelayedexpansion

REM Configuration
set BACKEND_DIR=%~dp0
if not defined GOOGLE_APPLICATION_CREDENTIALS (
    set GOOGLE_APPLICATION_CREDENTIALS=%BACKEND_DIR%firebas.json
)
if not defined FIREBASE_PROJECT_ID (
    set FIREBASE_PROJECT_ID=unienroll-85a57
)

REM Colors don't work in batch, use visual separators instead
set SEPARATOR============================================

if "%1"=="" (
    echo.
    echo Firebase Database Management Tool
    echo %SEPARATOR%
    echo.
    echo Usage: db-manage.bat [command]
    echo.
    echo Commands:
    echo   init        Initialize collections and seed data
    echo   seed        Only seed sample data
    echo   init-empty  Initialize collections without seeding  
    echo   clear       Clear all data (destructive^!)
    echo   help        Show this help message
    echo.
    echo Examples:
    echo   db-manage.bat init
    echo   db-manage.bat clear
    echo.
    goto :eof
)

REM Check credentials
if not exist "%GOOGLE_APPLICATION_CREDENTIALS%" (
    echo ERROR: Firebase credentials file not found
    echo Expected: %GOOGLE_APPLICATION_CREDENTIALS%
    goto :error
)

echo Credentials verified
echo Project: %FIREBASE_PROJECT_ID%
echo.

REM Change to backend directory
cd /d "%BACKEND_DIR%"

if /i "%1"=="init" (
    echo Initializing Firebase database...
    echo.
    go run cmd/init-db/main.go -init=true -seed=true
    if errorlevel 1 goto :error
    echo.
    echo SUCCESS: Database initialization complete!
    goto :eof
)

if /i "%1"=="seed" (
    echo Seeding sample data...
    echo.
    go run cmd/init-db/main.go -init=false -seed=true
    if errorlevel 1 goto :error
    echo.
    echo SUCCESS: Sample data seeded!
    goto :eof
)

if /i "%1"=="init-empty" (
    echo Initializing empty collections...
    echo.
    go run cmd/init-db/main.go -init=true -seed=false
    if errorlevel 1 goto :error
    echo.
    echo SUCCESS: Collections initialized!
    goto :eof
)

if /i "%1"=="clear" (
    echo WARNING: About to clear ALL data from database!
    echo This action CANNOT be undone!
    echo.
    set /p response="Type 'clear all' to confirm, or press Ctrl+C to cancel: "
    
    if /i not "!response!"=="clear all" (
        echo Aborted.
        goto :eof
    )
    
    echo.
    echo Clearing all data...
    echo.
    go run cmd/init-db/main.go -clear=true
    if errorlevel 1 goto :error
    echo.
    echo SUCCESS: All data cleared!
    goto :eof
)

if /i "%1"=="help" (
    go run cmd/init-db/main.go -help
    goto :eof
)

echo Unknown command: %1
call :print_help
goto :error

:error
echo.
echo ERROR: Operation failed
exit /b 1

:eof
endlocal

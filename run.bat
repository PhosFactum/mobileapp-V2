@echo off
chcp 65001 > nul

echo [1/2] Generating Swagger documentation...
swag init -g cmd/app/main.go

if %ERRORLEVEL% neq 0 (
    echo [ERROR] Swagger generation failed
    pause
    exit /b 1
)

echo [2/2] Starting Go application server...
echo ----------------------------------------
echo App will be available at http://localhost:8080
echo Press Ctrl+C to stop
echo ----------------------------------------
go run cmd/app/main.go

pause
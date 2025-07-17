@echo off
echo 🚀 Generating Swagger documentation for Go-Metro API...

REM Check if swag is installed
swag --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ swag CLI not found. Installing...
    go install github.com/swaggo/swag/cmd/swag@latest
)

REM Generate Swagger docs
echo 📝 Generating docs...
swag init

REM Check if generation was successful
if %errorlevel% equ 0 (
    echo ✅ Swagger documentation generated successfully!
    echo 📁 Docs folder created at: ./docs/
    echo.
    echo 🌐 To view the documentation:
    echo    1. Run: go run main.go
    echo    2. Open: http://localhost:8080/swagger/index.html
    echo.
    echo 📋 Available endpoints:
    echo    - Authentication: /auth/register, /auth/login
    echo    - User Management: /user/* (requires auth)
    echo    - Card Management: /card/*
    echo    - History Management: /history/*
    echo    - Admin Management: /admin/* (requires admin role)
    echo    - Health Check: /health
) else (
    echo ❌ Failed to generate Swagger documentation
    echo 💡 Check your Swagger annotations and try again
    pause
    exit /b 1
)

pause 
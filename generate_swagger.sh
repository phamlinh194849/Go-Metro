#!/bin/bash

# Swagger Documentation Generator Script
# This script generates Swagger documentation for the Go-Metro API

echo "🚀 Generating Swagger documentation for Go-Metro API..."

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "❌ swag CLI not found. Installing..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate Swagger docs
echo "📝 Generating docs..."
swag init

# Check if generation was successful
if [ $? -eq 0 ]; then
    echo "✅ Swagger documentation generated successfully!"
    echo "📁 Docs folder created at: ./docs/"
    echo ""
    echo "🌐 To view the documentation:"
    echo "   1. Run: go run main.go"
    echo "   2. Open: http://localhost:8080/swagger/index.html"
    echo ""
    echo "📋 Available endpoints:"
    echo "   - Authentication: /auth/register, /auth/login"
    echo "   - User Management: /user/* (requires auth)"
    echo "   - Card Management: /card/*"
    echo "   - History Management: /history/*"
    echo "   - Admin Management: /admin/* (requires admin role)"
    echo "   - Health Check: /health"
else
    echo "❌ Failed to generate Swagger documentation"
    echo "💡 Check your Swagger annotations and try again"
    exit 1
fi 
#!/bin/bash

# Initialize go module
echo "Initializing go module..."
go mod init spider > /dev/null 2>&1

# Install dependencies
echo "Installing dependencies..."
go get github.com/go-chi/chi/v5 > /dev/null 2>&1
go get github.com/joho/godotenv > /dev/null 2>&1
go get github.com/lib/pq > /dev/null 2>&1
go get github.com/google/uuid > /dev/null 2>&1

echo -e "\nStart the app by running: \n\033[35m$ go run .\033[0m"
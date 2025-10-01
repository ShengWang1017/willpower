#!/bin/bash
# Linux Build Script for Willpower Forge
# This script builds the frontend and backend into a single executable

set -e  # Exit on error

echo "===================================="
echo "Building Willpower Forge for Linux"
echo "===================================="

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "Error: npm is not installed"
    echo "Please install Node.js from https://nodejs.org/"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    echo "Please install Go from https://go.dev/dl/"
    exit 1
fi

echo ""
echo "Step 1/4: Installing frontend dependencies..."
cd willpower-forge-web
npm install

echo ""
echo "Step 2/4: Building frontend..."
npm run build
cd ..

echo ""
echo "Step 3/4: Copying frontend files to backend..."
mkdir -p willpower-forge-api/web
cp -r willpower-forge-web/dist willpower-forge-api/web/

echo ""
echo "Step 4/4: Building Linux executable..."
cd willpower-forge-api
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../willpower-forge-linux .
cd ..

echo ""
echo "===================================="
echo "Build completed successfully!"
echo "===================================="
echo ""
echo "Executable: willpower-forge-linux"
echo "Size: $(du -h willpower-forge-linux | cut -f1)"
echo ""
echo "To run the application:"
echo "  1. chmod +x willpower-forge-linux"
echo "  2. ./willpower-forge-linux"
echo "  3. Open browser and visit http://localhost:8080"
echo ""

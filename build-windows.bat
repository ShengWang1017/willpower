@echo off
REM Windows Build Script for Willpower Forge
REM This script builds the frontend and backend into a single executable

echo ====================================
echo Building Willpower Forge for Windows
echo ====================================

REM Check if Node.js is installed
where npm >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Error: npm is not installed or not in PATH
    echo Please install Node.js from https://nodejs.org/
    pause
    exit /b 1
)

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://go.dev/dl/
    pause
    exit /b 1
)

echo.
echo Step 1/4: Installing frontend dependencies...
cd willpower-forge-web
call npm install
if %ERRORLEVEL% neq 0 (
    echo Error: Failed to install frontend dependencies
    cd ..
    pause
    exit /b 1
)

echo.
echo Step 2/4: Building frontend...
call npm run build
if %ERRORLEVEL% neq 0 (
    echo Error: Failed to build frontend
    cd ..
    pause
    exit /b 1
)
cd ..

echo.
echo Step 3/4: Copying frontend files to backend...
if not exist "willpower-forge-api\web" mkdir willpower-forge-api\web
xcopy /E /I /Y willpower-forge-web\dist willpower-forge-api\web\dist

echo.
echo Step 4/4: Building Windows executable...
cd willpower-forge-api
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o ..\willpower-forge-windows.exe .
if %ERRORLEVEL% neq 0 (
    echo Error: Failed to build Windows executable
    cd ..
    pause
    exit /b 1
)
cd ..

echo.
echo ====================================
echo Build completed successfully!
echo ====================================
echo.
echo Executable: willpower-forge-windows.exe
echo Size:
dir willpower-forge-windows.exe | find "willpower-forge-windows.exe"
echo.
echo To run the application:
echo   1. Double-click willpower-forge-windows.exe
echo   2. Open browser and visit http://localhost:8080
echo.
pause

@echo off
chcp 65001 >nul
setlocal EnableDelayedExpansion

echo ========================================
echo   Gomoku Server Build Script
echo ========================================
echo.

cd /d "%~dp0"

set APP_NAME=gomoku-server
set VERSION=1.0.0
set BUILD_DIR=build
set MAIN_PATH=./cmd/main.go

rem »ñÈ¡ git commit hash
for /f "delims=" %%i in ('git rev-parse --short HEAD 2^>nul') do set GIT_COMMIT=%%i
if "%GIT_COMMIT%"=="" set GIT_COMMIT=unknown

set LDFLAGS=-s -w -X main.Version=%VERSION% -X main.GitCommit=%GIT_COMMIT%

echo Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%APP_NAME%-windows-amd64.exe %MAIN_PATH%
if errorlevel 1 (
    echo [ERROR] Windows build failed!
    exit /b 1
)
echo [OK] Windows amd64 build complete.

echo Building for Windows (386)...
set GOOS=windows
set GOARCH=386
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%APP_NAME%-windows-386.exe %MAIN_PATH%
if errorlevel 1 (
    echo [ERROR] Windows 386 build failed!
    exit /b 1
)
echo [OK] Windows 386 build complete.

echo Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%APP_NAME%-linux-amd64 %MAIN_PATH%
if errorlevel 1 (
    echo [ERROR] Linux amd64 build failed!
    exit /b 1
)
echo [OK] Linux amd64 build complete.

echo Building for Linux (arm64)...
set GOOS=linux
set GOARCH=arm64
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%APP_NAME%-linux-arm64 %MAIN_PATH%
if errorlevel 1 (
    echo [ERROR] Linux arm64 build failed!
    exit /b 1
)
echo [OK] Linux arm64 build complete.

echo Building for Darwin (macOS amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%APP_NAME%-darwin-amd64 %MAIN_PATH%
if errorlevel 1 (
    echo [ERROR] Darwin amd64 build failed!
    exit /b 1
)
echo [OK] Darwin amd64 build complete.

echo Building for Darwin (macOS arm64)...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%APP_NAME%-darwin-arm64 %MAIN_PATH%
if errorlevel 1 (
    echo [ERROR] Darwin arm64 build failed!
    exit /b 1
)
echo [OK] Darwin arm64 build complete.

echo.
echo ========================================
echo   Build Complete!
echo ========================================
echo.
echo Output directory: %BUILD_DIR%\
dir /b %BUILD_DIR%
echo.

pause

#!/bin/bash

echo "========================================"
echo "  Gomoku Server Build Script"
echo "========================================"
echo

cd "$(dirname "$0")"

APP_NAME="gomoku-server"
VERSION="1.0.0"
BUILD_DIR="build"
MAIN_PATH="./cmd/main.go"

# Get git commit hash
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS="-s -w -X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT}"

# Create build directory
mkdir -p ${BUILD_DIR}

echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${APP_NAME}-windows-amd64.exe ${MAIN_PATH}
[ $? -eq 0 ] && echo "[OK] Windows amd64 build complete." || { echo "[ERROR] Windows build failed!"; exit 1; }

echo "Building for Windows (386)..."
GOOS=windows GOARCH=386 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${APP_NAME}-windows-386.exe ${MAIN_PATH}
[ $? -eq 0 ] && echo "[OK] Windows 386 build complete." || { echo "[ERROR] Windows 386 build failed!"; exit 1; }

echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${APP_NAME}-linux-amd64 ${MAIN_PATH}
[ $? -eq 0 ] && echo "[OK] Linux amd64 build complete." || { echo "[ERROR] Linux amd64 build failed!"; exit 1; }

echo "Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${APP_NAME}-linux-arm64 ${MAIN_PATH}
[ $? -eq 0 ] && echo "[OK] Linux arm64 build complete." || { echo "[ERROR] Linux arm64 build failed!"; exit 1; }

echo "Building for Darwin (macOS amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${APP_NAME}-darwin-amd64 ${MAIN_PATH}
[ $? -eq 0 ] && echo "[OK] Darwin amd64 build complete." || { echo "[ERROR] Darwin amd64 build failed!"; exit 1; }

echo "Building for Darwin (macOS arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${APP_NAME}-darwin-arm64 ${MAIN_PATH}
[ $? -eq 0 ] && echo "[OK] Darwin arm64 build complete." || { echo "[ERROR] Darwin arm64 build failed!"; exit 1; }

echo
echo "========================================"
echo "  Build Complete!"
echo "========================================"
echo
echo "Output directory: ${BUILD_DIR}/"
ls -lh ${BUILD_DIR}/
echo

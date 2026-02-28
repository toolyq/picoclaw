@echo off
setlocal

:: 设置交叉编译环境变量
set GOOS=linux
set GOARCH=amd64

echo [INFO] Checking Go environment...
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Go is not installed or not in your PATH.
    exit /b 1
)

echo [INFO] Downloading dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Failed to download dependencies.
    exit /b 1
)

echo [INFO] Staging workspace for embedding...
if exist "cmd\picoclaw\internal\onboard\workspace" rmdir /s /q "cmd\picoclaw\internal\onboard\workspace"
xcopy "workspace" "cmd\picoclaw\internal\onboard\workspace\" /E /I /Y /Q >nul
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Failed to stage workspace.
    exit /b 1
)

echo [INFO] Building for Ubuntu AMD64 (Linux/amd64)...
if not exist "build" mkdir "build"

:: 构建并添加版本信息（参考 Makefile 中的 LDFLAGS）
echo [INFO] Building PicoClaw...
go build -v -tags stdjson -ldflags "-s -w" -o build/picoclaw-linux-amd64 ./cmd/picoclaw
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Failed to build PicoClaw.
    set BUILD_STATUS=%ERRORLEVEL%
)

echo [INFO] Building PicoClaw Launcher...
go build -v -tags stdjson -ldflags "-s -w" -o build/picoclaw-launcher-linux-amd64 ./cmd/picoclaw-launcher
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Failed to build PicoClaw Launcher.
    set BUILD_STATUS=%ERRORLEVEL%
)

echo [INFO] Building PicoClaw Launcher TUI...
go build -v -tags stdjson -ldflags "-s -w" -o build/picoclaw-launcher-tui-linux-amd64 ./cmd/picoclaw-launcher-tui
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Failed to build PicoClaw Launcher TUI.
    set BUILD_STATUS=%ERRORLEVEL%
)

echo [INFO] Cleaning up staged workspace...
if exist "cmd\picoclaw\internal\onboard\workspace" rmdir /s /q "cmd\picoclaw\internal\onboard\workspace"

if %BUILD_STATUS% EQU 0 (
    echo [SUCCESS] Build successful! 
    echo Output: build/picoclaw-linux-amd64
) else (
    echo [ERROR] Build failed!
    exit /b %BUILD_STATUS%
)

pause

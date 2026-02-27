setlocal

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

echo [INFO] Building PicoClaw...
go build -ldflags "-s -w" -o picoclaw.exe ./cmd/picoclaw
set BUILD_STATUS=%ERRORLEVEL%

echo [INFO] Cleaning up staged workspace...
if exist "cmd\picoclaw\internal\onboard\workspace" rmdir /s /q "cmd\picoclaw\internal\onboard\workspace"

if %BUILD_STATUS% EQU 0 (
    echo [SUCCESS] Build successful! Output: picoclaw.exe
) else (
    echo [ERROR] Build failed!
    exit /b %BUILD_STATUS%
)

cmd /k
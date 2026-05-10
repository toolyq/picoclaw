@echo off
setlocal enabledelayedexpansion

echo [INFO] Building PicoClaw Web...

:: Check for pnpm
where pnpm >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    set PNPM_CMD=pnpm
) else (
    echo [INFO] pnpm not found, falling back to npx pnpm...
    set PNPM_CMD=npx -y pnpm
)

echo [INFO] Using command: %PNPM_CMD%

echo [INFO] Building Web Frontend...
pushd web\frontend
call %PNPM_CMD% install
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Failed to install Web Frontend dependencies.
    popd
    exit /b 1
)

call %PNPM_CMD% build:backend
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Failed to build Web Frontend.
    popd
    exit /b 1
)
popd

echo [INFO] Building Web Backend...
pushd web\backend
go build -ldflags "-s -w" -o ..\..\picoclaw-web.exe .
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Failed to build Web Backend.
    popd
    exit /b 1
)
popd

echo [SUCCESS] Web build successful! Output: picoclaw-web.exe
exit /b 0

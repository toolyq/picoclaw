# Set environment variables for cross-compilation
$env:GOOS = "linux"
$env:GOARCH = "arm64"

# Get version and commit information (if Git is installed)
$VERSION = try { git describe --tags --always --dirty 2>$null } catch { "dev" }
$GIT_COMMIT = try { git rev-parse --short=8 HEAD 2>$null } catch { "dev" }
$BUILD_TIME = Get-Date -Format "yyyy-MM-ddTHH:mm:ssK"
$GO_VERSION = (go version).Split(' ')[2]

# Build variables
$BINARY_NAME = "picoclaw"
$BUILD_DIR = "build"
$CMD_DIR = "cmd/picoclaw"
$LDFLAGS = "-X main.version=$VERSION -X main.gitCommit=$GIT_COMMIT -X main.buildTime=$BUILD_TIME -X main.goVersion=$GO_VERSION -s -w"

Write-Host "[INFO] 正在为 Ubuntu ARM64 (Linux/arm64) 编译 $BINARY_NAME..."

# Create build directory
if (!(Test-Path -Path $BUILD_DIR)) {
    New-Item -ItemType Directory -Path $BUILD_DIR
}

# Run go build
go build -ldflags $LDFLAGS -tags stdjson -o "$BUILD_DIR/$BINARY_NAME-linux-arm64" "./$CMD_DIR"

if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] 编译成功！"
    Write-Host "输出文件路径: $BUILD_DIR/$BINARY_NAME-linux-arm64"
} else {
    Write-Host "[ERROR] 编译失败！"
    exit $LASTEXITCODE
}

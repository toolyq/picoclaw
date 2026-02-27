#!/bin/bash

# 设置变量
BINARY_NAME="picoclaw"
BUILD_DIR="build"
CMD_DIR="cmd/picoclaw"

# 获取版本和构建信息
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT=$(git rev-parse --short=8 HEAD 2>/dev/null || echo "dev")
BUILD_TIME=$(date +%FT%T%z)
GO_VERSION=$(go version | awk '{print $3}')

# 编译参数
LDFLAGS="-ldflags \"-X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT} -X main.buildTime=${BUILD_TIME} -X main.goVersion=${GO_VERSION} -s -w\""

echo "[INFO] 正在为 Ubuntu ARM64 (Linux/arm64) 编译 ${BINARY_NAME}..."

# 确保构建目录存在
mkdir -p ${BUILD_DIR}

# 清理并准备 workspace
echo "[INFO] Staging workspace for embedding..."
rm -rf cmd/picoclaw/internal/onboard/workspace
cp -r workspace cmd/picoclaw/internal/onboard/workspace

# 设置环境变量并进行交叉编译
GOOS=linux GOARCH=arm64 go build -v -tags stdjson ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-arm64 ./${CMD_DIR}
BUILD_STATUS=$?

# 清理
rm -rf cmd/picoclaw/internal/onboard/workspace

if [ $BUILD_STATUS -eq 0 ]; then
    echo "[SUCCESS] 编译成功！"
    echo "输出文件路径: ${BUILD_DIR}/${BINARY_NAME}-linux-arm64"
else
    echo "[ERROR] 编译失败！"
    exit 1
fi

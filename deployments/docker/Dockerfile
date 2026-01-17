# 多阶段构建 Dockerfile
#
# Game Slots Version Server Docker 镜像
#
# 使用方法：
#   docker build -f deployments/docker/Dockerfile -t game-slots-vsn:latest .
#
# 运行容器：
#   docker run -p 8080:8080 game-slots-vsn:latest

# ==================== 阶段1: 构建阶段 ====================
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /build

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖（利用 Docker 缓存）
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（静态链接）
# -ldflags="-s -w" 去除调试信息和符号表，减小二进制大小
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/server \
    ./cmd/server

# ==================== 阶段2: 运行阶段 ====================
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/server /app/server

# 复制 GeoIP 数据（如果需要）
# COPY deployments/data/geoip.mmdb /app/data/geoip.mmdb

# 修改文件权限
RUN chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health/live || exit 1

# 设置环境变量
ENV GIN_MODE=release \
    SERVER_HOST=0.0.0.0 \
    SERVER_PORT=8080 \
    LOG_LEVEL=info

# 启动应用
ENTRYPOINT ["/app/server"]
CMD ["serve"]

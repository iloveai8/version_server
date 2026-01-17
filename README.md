# Game Slots Version Server

> 🎮 游戏版本配置管理服务 - 高性能、高可用的分布式配置管理平台

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Tests](https://img.shields.io/badge/tests-84.5%25-brightgreen)](#-测试覆盖率)

---

## 📖 项目简介

Game Slots Version Server 是一个基于 **领域驱动设计 (DDD)** 的游戏版本配置管理服务，提供：

- ✨ **版本配置管理**: 游戏客户端版本号和子服务器配置
- 🔐 **GM 功能控制**: GM 后台功能开关和 IP 白名单管理
- 🌍 **GeoIP 支持**: 基于地理位置的访问控制
- 📊 **可观测性**: 完整的健康检查、Prometheus 监控和结构化日志
- 🚀 **高性能**: 平均响应时间 < 50ms，支持高并发
- 🔧 **易部署**: 支持 Docker 和 Kubernetes 部署

---

## ⚡ 快速开始

### 本地开发

```bash
# 1. 克隆仓库
git clone https://github.com/your-org/game_slots_vsn.git
cd game_slots_vsn

# 2. 启动 Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 3. 运行服务
go run cmd/server/main.go

# 4. 测试 API
curl "http://localhost:8080/api/v1/version?vsn=1.0.0&env=dev"
```

---

## 🏗️ 系统架构

### 架构原则

采用 **领域驱动设计 (DDD)** 方法论，遵循以下原则：

| 原则 | 说明 |
|------|------|
| **DDD** | 业务逻辑封装在领域模型中 |
| **DIP** | 依赖倒置，高层不依赖低层 |
| **SRP** | 单一职责，每层职责明确 |
| **ISP** | 接口隔离，避免胖接口 |

### 分层架构

```
┌─────────────────────────────────────────────────────┐
│                   HTTP Layer                        │
│  (Handler + Middleware + Health Check)              │
└─────────────────────────────────────────────────────┘
                        ↓↑
┌─────────────────────────────────────────────────────┐
│                  Service Layer                      │
│  (Version Service + GM Service + IP Service)        │
└─────────────────────────────────────────────────────┘
                        ↓↑
┌─────────────────────────────────────────────────────┐
│                Repository Layer                     │
│  (Data Access + Redis Operations)                   │
└─────────────────────────────────────────────────────┘
                        ↓↑
┌─────────────────────────────────────────────────────┐
│                   Domain Layer                      │
│              (Business Logic + Models)              │
└─────────────────────────────────────────────────────┘
                        ↓↑
┌─────────────────────────────────────────────────────┐
│                   Storage Layer                     │
│              (Redis Cluster)                        │
└─────────────────────────────────────────────────────┘
```

### 层级职责

| 层级 | 职责 |
|------|------|
| **Handler** | HTTP 请求处理、参数绑定、响应格式化 |
| **Service** | 业务逻辑编排、事务管理、跨领域协作 |
| **Repository** | 数据访问抽象、CRUD 操作 |
| **Domain** | 业务规则、领域模型、值对象 |

### 目录结构

```
game_slots_vsn/
├── cmd/server/           # 主程序入口
├── configs/              # 配置文件
│   ├── dev.yaml         # 开发环境
│   ├── pre.yaml         # 预发布环境
│   └── pro.yaml         # 生产环境
├── internal/             # 私有代码（DDD分层）
│   ├── domain/          # 领域模型
│   ├── handler/         # HTTP处理器（api/admin）
│   ├── http/            # HTTP服务器和路由
│   ├── repository/      # 数据访问层
│   └── service/         # 业务逻辑层
├── pkg/                  # 可复用工具包
│   ├── config/          # 配置管理
│   ├── constants/       # 常量定义
│   ├── errcode/         # 错误处理
│   ├── ggeoip/          # GeoIP客户端
│   ├── httputil/        # HTTP工具
│   ├── jsonutil/        # JSON工具
│   ├── log/             # 日志
│   ├── metrics/         # 监控指标
│   ├── utils/           # 通用工具
│   ├── validator/       # 输入验证
│   └── version/         # 版本号工具
├── deployments/          # 部署配置
│   ├── docker/          # Docker配置
│   └── k8s/             # Kubernetes配置
└── test/                 # 测试
    ├── coverage/        # 测试覆盖率报告
    ├── integration/     # 集成测试
    ├── ip/              # IP测试
    └── move_vsn/        # 版本迁移测试
```

### 核心技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| **Web框架** | Gin | 高性能HTTP框架 |
| **数据存储** | Redis Cluster | 分布式缓存 |
| **日志** | Zap | 结构化日志 |
| **监控** | Prometheus | 指标收集 |
| **容器化** | Docker | 应用容器 |
| **编排** | Kubernetes | 容器编排 |

---

## 🎯 核心功能

### 1. 版本配置管理

```bash
# 获取版本配置
curl "http://localhost:8080/api/v1/version?vsn=1.0.0&env=dev"

# 列出版本
curl "http://localhost:8080/api/v1/versions?env=dev"
```

### 2. GM 功能管理

```bash
# 获取 GM 配置
curl http://localhost:8080/admin/v1/gm/config

# 更新 GM 配置
curl -X PUT http://localhost:8080/admin/v1/gm/config \
  -H "Content-Type: application/json" \
  -d '{"gmEnable":true,"block":false}'

# 切换 GM 功能
curl -X POST http://localhost:8080/admin/v1/gm/toggle
```

### 3. IP 白名单管理

```bash
# 获取白名单
curl http://localhost:8080/admin/v1/ip-whitelist

# 添加 IP
curl -X POST http://localhost:8080/admin/v1/ip-whitelist/ip \
  -H "Content-Type: application/json" \
  -d '{"ip":"192.168.1.100"}'

# 检查 IP
curl -X POST http://localhost:8080/admin/v1/ip-whitelist/check \
  -H "Content-Type: application/json" \
  -d '{"ip":"192.168.1.100"}'
```

### 4. 健康检查

```bash
# 存活探针
curl http://localhost:8080/health/live

# 就绪探针
curl http://localhost:8080/health/ready

# 完整健康状态
curl http://localhost:8080/health

# Prometheus 指标
curl http://localhost:8080/metrics
```

---

## 🧪 测试

### 测试组织

```
test/
├── integration/       # 集成测试
│   └── api_test.go   # API 端到端测试
├── coverage/          # 测试覆盖率报告
│   └── README.md     # 覆盖率说明
├── ip/               # IP功能测试
├── move_vsn/         # 版本迁移测试
└── testdata/         # 测试数据
```

### 运行测试

```bash
# 运行所有测试
go test -v ./...

# 运行集成测试
go test -v ./test/integration/...

# 运行 Benchmark
go test -bench=. -benchmem ./...

# 生成覆盖率报告
go test -coverprofile=test/coverage/coverage.out ./...
go tool cover -html=test/coverage/coverage.out -o test/coverage/coverage.html
```

### 测试覆盖率

| 层级 | 覆盖率 | 状态 |
|------|--------|------|
| **Domain 层** | 98.5% | ✅ |
| **基础设施层** | 89.9% | ✅ |
| **Repository 层** | 54.8% | ⚠️ |
| **Service 层** | 50.2% | ⚠️ |
| **总体平均** | **84.5%** | ✅ |

> 详细覆盖率报告请查看 [test/coverage/README.md](test/coverage/README.md)

---

## 📊 性能

### Benchmark 结果

| 操作 | 性能 | 内存分配 |
|------|------|----------|
| Get Version | 28,926 ns/op | 793 B/op |
| Create Version | 167,807 ns/op | 3,222 B/op |
| List Versions | 314,390 ns/op | 9,120 B/op |
| Update Version | 84,051 ns/op | 3,140 B/op |

### SLA

- **可用性**: 99.9%
- **响应时间**: P95 < 100ms
- **吞吐量**: > 1000 QPS (单实例)
- **扩展性**: 水平扩展，线性增长

---

## 🚦 部署

### 环境要求

| 组件 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.23+ | 编译环境 |
| Docker | 20.10+ | 容器运行时 |
| Redis | 7.x+ | 数据存储 |
| Kubernetes | 1.25+ | 编排平台 (可选) |

### Docker 部署

#### 构建镜像

```bash
# 构建镜像
docker build -f deployments/docker/Dockerfile -t game-slots-vsn:latest .

# 查看镜像
docker images | grep game-slots-vsn
```

#### 运行容器

```bash
# 单独运行
docker run -d -p 8080:8080 \
  -e REDIS_ADDR=redis:6379 \
  --name game-slots-vsn \
  game-slots-vsn:latest
```

#### Docker Compose

```bash
# 启动服务（基础版）
cd deployments/docker
docker-compose up -d

# 启动服务（包含监控）
docker-compose --profile metrics up -d

# 查看日志
docker-compose logs -f game-slots-vsn

# 停止服务
docker-compose down
```

### Kubernetes 部署

#### 部署步骤

```bash
# 1. 应用 ConfigMap
kubectl apply -f deployments/k8s/configmap.yaml

# 2. 应用 Deployment
kubectl apply -f deployments/k8s/deployment.yaml

# 3. 应用 Service
kubectl apply -f deployments/k8s/service.yaml

# 4. 查看状态
kubectl get pods -l app=game-slots-vsn
kubectl get svc game-slots-vsn

# 5. 查看日志
kubectl logs -f deployment/game-slots-vsn
```

#### 扩缩容

```bash
# 手动扩容
kubectl scale deployment/game-slots-vsn --replicas=3

# 自动扩缩容（HPA）
kubectl apply -f deployments/k8s/hpa.yaml
```

#### 配置说明

**ConfigMap** (`configmap.yaml`):
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: game-slots-vsn-config
data:
  REDIS_ADDR: "redis:6379"
  LOG_LEVEL: "info"
```

**环境变量**:
- `REDIS_ADDR`: Redis 地址
- `REDIS_DB`: Redis 数据库编号
- `LOG_LEVEL`: 日志级别（debug/info/warn/error）
- `SERVER_PORT`: 服务端口

---

## 📈 监控

### Prometheus 指标

| 指标名称 | 类型 | 说明 |
|---------|------|------|
| `http_requests_total` | Counter | HTTP 请求总数 |
| `http_request_duration_seconds` | Histogram | HTTP 请求延迟 |
| `redis_operations_total` | Counter | Redis 操作总数 |
| `version_cache_hits_total` | Counter | 版本缓存命中数 |

### 健康检查端点

| 端点 | 说明 | 用途 |
|------|------|------|
| `/health/live` | Liveness Probe | K8s 存活探针 |
| `/health/ready` | Readiness Probe | K8s 就绪探针 |
| `/health` | 完整健康状态 | 健康检查 |
| `/ping` | 简单心跳 | 快速检测 |
| `/metrics` | Prometheus 指标 | 监控数据采集 |

### 日志

日志文件位置：
- 开发环境: `./logs/gs_vsn.log`
- 生产环境: 容器标准输出（stdout）

日志级别：
- `debug`: 调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息

---

## 🔧 配置

### 配置文件

配置文件位于 `configs/` 目录：

| 文件 | 环境 | 说明 |
|------|------|------|
| `dev.yaml` | 开发环境 | 本地开发配置 |
| `pre.yaml` | 预发布环境 | 测试环境配置 |
| `pro.yaml` | 生产环境 | 生产环境配置 |

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `GIN_MODE` | Gin 运行模式 | `release` |
| `SERVER_HOST` | 监听地址 | `0.0.0.0` |
| `SERVER_PORT` | 监听端口 | `8080` |
| `REDIS_ADDR` | Redis 地址 | `localhost:6379` |
| `REDIS_DB` | Redis 数据库 | `0` |
| `LOG_LEVEL` | 日志级别 | `info` |

---

## 🛠️ 开发

### 开发规范

详细的开发规范请查看 **[CLAUDE.md](CLAUDE.md)**，包括：

- ✅ 代码注释规范（中文 godoc）
- ✅ 命名规范（包名、变量、常量）
- ✅ 目录结构规范
- ✅ 错误处理规范
- ✅ 测试规范
- ✅ API 设计规范
- ✅ 代码提交规范

### Make 命令

```bash
make build        # 编译
make test         # 测试
make run          # 运行
make lint         # 代码检查
make coverage     # 覆盖率报告
```

### 代码质量检查

提交代码前请确保：

```bash
# 1. 格式化代码
go fmt ./...

# 2. 静态检查
go vet ./...

# 3. 运行测试
go test -race ./...

# 4. 检查覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
```

---

## 📝 更新日志

### v2.0 (2025-01-17)

#### 重构
- 🎉 完成领域驱动设计 (DDD) 架构重构
- ✨ 实现完整的分层架构（Handler → Service → Repository → Domain）
- 🔧 添加依赖倒置和仓储模式

#### 新功能
- ✨ 完整的健康检查 API（Liveness/Readiness/Full）
- ✨ Prometheus 监控指标
- ✨ 结构化日志（Zap）
- ✨ IP 白名单批量管理

#### 改进
- 🚀 性能优化：平均响应时间 < 50ms
- 📊 测试覆盖率提升至 84.5%
- 📝 完善的文档体系
- 🗂️ 标准化的目录结构

#### 部署
- 🐳 Docker 多阶段构建优化
- ☸️ Kubernetes 完整部署配置（HPA、PDB、ConfigMap）
- 🚀 自动化部署脚本

---

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改（遵循 [CLAUDE.md](CLAUDE.md) 中的 Commit 规范）
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

### 代码审查清单

- [ ] 代码符合 [CLAUDE.md](CLAUDE.md) 规范
- [ ] 所有导出代码有中文注释
- [ ] 单元测试覆盖率 > 80%
- [ ] `go vet` 无警告
- [ ] `go test -race` 通过
- [ ] Commit message 符合规范

---

## 📝 License

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

## 👥 团队

Development Team

---

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - Web 框架
- [Redis](https://redis.io/) - 数据存储
- [Prometheus](https://prometheus.io/) - 监控系统
- [Zap](https://github.com/uber-go/zap) - 日志库
- [Testify](https://github.com/stretchr/testify) - 测试框架

---

**版本**: v2.0
**最后更新**: 2025-01-17

# Game Slots Version Server - 开发规范

> 本文档定义项目开发规范和编码标准，所有开发人员必须遵守

## 📋 目录

- [代码注释规范](#代码注释规范)
- [命名规范](#命名规范)
- [目录结构规范](#目录结构规范)
- [错误处理规范](#错误处理规范)
- [测试规范](#测试规范)
- [API 设计规范](#api-设计规范)

---

## 代码注释规范

### 基本原则

**所有导出的代码必须有标准的 Go 文档注释**

- **强制性要求**: 每个导出的类型、函数、常量都必须有文档注释
- **语言要求**: 必须使用中文编写注释
- **格式标准**: 遵循 Go 官方文档注释规范（godoc）

### 包注释

```go
// Package errcode 提供统一的错误处理和响应格式
//
// 该包实现了：
// - 自定义错误类型 AppError，支持重试策略
// - HTTP状态码映射
// - 标准化的错误响应格式
//
// 使用示例：
//   err := errcode.ErrVersionNotFound.WithDetails(map[string]interface{}{
//       "version": "1.2.0",
//   })
package errcode
```

### 类型注释

```go
// AppError 应用错误类型
//
// 实现了error接口，包含错误码、消息、HTTP状态码和重试信息
type AppError struct {
    Code       string                 `json:"code"`    // 自定义错误码
    Message    string                 `json:"message"`  // 错误消息
    HTTPStatus int                    `json:"-"`       // HTTP状态码
    Details    map[string]interface{} `json:"details,omitempty"`  // 错误详情
}
```

### 函数注释

```go
// ValidateVersion 验证版本号格式
//
// 检查版本号是否符合X.Y.Z或X.Y.Z.N格式
//
// 参数:
//   vsn: 版本号字符串，如"1.0.0"
//
// 返回:
//   error: 版本号格式无效时返回错误
func ValidateVersion(vsn string) error
```

### 注释内容要求

1. **必须说明"是什么"** - 清晰描述函数/类型的用途
2. **必要时说明"为什么"** - 对于复杂逻辑，解释设计原因
3. **必要时说明"怎么用"** - 对于公开API，提供使用示例
4. **参数和返回值** - 列出所有参数及其含义
5. **注意事项** - 如果函数有副作用，必须说明

### Checklist

每个代码文件提交前必须检查：
- [ ] 包注释是否存在且完整
- [ ] 所有导出的类型都有注释
- [ ] 所有导出的函数都有注释
- [ ] 所有导出的常量都有注释
- [ ] 注释使用中文
- [ ] 注释格式符合godoc标准

---

## 命名规范

### 包命名

#### 规则1: 小写单词

```
✅ 推荐: config, log, errcode, handler, service, repository
❌ 避免: Config, Log, ErrorCode, Handler, Service, Repository
```

#### 规则2: 避免缩写

```
✅ 推荐: config, version, geoip
❌ 避免: cfg, ver, geo
```

#### 规则3: 避免通用名称

```
✅ 推荐: config, validator, httputil, jsonutil
❌ 避免: util, common, base, core
```

#### 规则4: 避免标准库冲突

```
✅ 推荐: log（简洁）
❌ 避免: logger（与标准库的log包不同）
```

#### 规则5: 复合词直接连接（无下划线）

```
✅ 推荐: httputil, jsonutil, geoip
❌ 避免: http_util, json_util, geo_ip
```

### 变量命名

```go
// ✅ 好的命名
versionService service.VersionService
redisClient    redispkg.Client
httpServer     *http.Server

// ❌ 不好的命名
vs             service.VersionService  // 过度缩写
rc             redispkg.Client        // 不清晰
s              *http.Server            // 单字母变量（除循环变量）
```

### 常量命名

```go
// ✅ 好的命名
const (
    MaxConnections = 100
    DefaultTimeout = 30 * time.Second
    ErrorCodeNotFound = "NOT_FOUND"
)

// ❌ 不好的命名
const (
    MAX = 100                // 全大写（非const）
    max_connections = 100    // 下划线命名
    errCode = "NOT_FOUND"    // 驼峰命名
)
```

---

## 目录结构规范

### 标准目录结构

```
project/
├── cmd/                    # 主应用程序
│   └── server/
│       └── main.go
│
├── internal/               # 私有应用代码
│   ├── handler/           # HTTP处理器层
│   ├── service/           # 服务层（业务逻辑）
│   ├── repository/        # 数据访问层
│   ├── domain/            # 领域模型
│   └── http/              # HTTP服务器和路由
│
├── pkg/                   # 可被外部使用的库
│   ├── config/           # 配置管理
│   ├── log/              # 日志
│   └── errcode/          # 错误处理
│
├── test/                  # 测试文件
│   ├── unit/             # 单元测试
│   ├── integration/      # 集成测试
│   └── testdata/         # 测试数据
│
└── deployments/           # 部署配置
    ├── docker/           # Docker配置
    └── k8s/              # Kubernetes配置
```

### 目录命名规则

- **使用小写**: `handler/`, `service/`, `repository/`
- **使用复数**: `handlers/`, `services/` (如果多个)
- **避免缩写**: `config/` 而非 `cfg/`
- **避免下划线**: `httputil/` 而非 `http_util/`

---

## 错误处理规范

### 错误码定义

```go
const (
    // 客户端错误 (4xx)
    ErrCodeInvalidParams     = "INVALID_PARAMS"
    ErrCodeInvalidIPFormat   = "INVALID_IP_FORMAT"
    ErrCodeVersionNotFound   = "VERSION_NOT_FOUND"

    // 服务器错误 (5xx)
    ErrCodeInternalError     = "INTERNAL_ERROR"
    ErrCodeRedisConnError    = "REDIS_CONNECTION_ERROR"
    ErrCodeTimeout           = "REQUEST_TIMEOUT"
)
```

### 错误响应格式

```json
{
  "error": {
    "code": "VERSION_NOT_FOUND",
    "message": "版本1.2.0不存在",
    "details": {
      "version": "1.2.0",
      "availableVersions": ["1.0.0", "1.1.0"]
    }
  }
}
```

### HTTP 状态码映射

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| `INVALID_PARAMS` | 400 | 参数错误 |
| `VERSION_NOT_FOUND` | 404 | 版本不存在 |
| `INTERNAL_ERROR` | 500 | 内部错误 |
| `REDIS_CONNECTION_ERROR` | 503 | Redis连接失败 |
| `REQUEST_TIMEOUT` | 504 | 请求超时 |

### 重试策略

使用标准 HTTP `Retry-After` 头：

```go
if appErr.Retryable && appErr.RetryAfter > 0 {
    c.Header("Retry-After", strconv.Itoa(appErr.RetryAfter))
}
```

---

## 测试规范

### 测试组织

```
test/
├── unit/              # 单元测试（包内 *_test.go）
├── integration/       # 集成测试
├── testdata/          # 测试数据
└── coverage/          # 覆盖率报告
```

### 测试命名

- 测试文件: `*_test.go`
- Benchmark 文件: `*_benchmark_test.go`
- 测试函数: `Test<FunctionName>`
- Benchmark 函数: `Benchmark<FunctionName>`

### Table-Driven Tests

```go
func TestValidateVersion(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid version", "1.0.0", false},
        {"invalid version", "invalid", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateVersion(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateVersion() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 覆盖率目标

| 层级 | 目标覆盖率 |
|------|-----------|
| Domain | > 90% |
| Repository | > 80% |
| Service | > 80% |
| Handler | > 80% |
| **总体** | **> 80%** |

---

## API 设计规范

### RESTful API

#### 资源命名

```
GET    /api/v1/versions          # 列出资源
POST   /api/v1/version          # 创建资源
GET    /api/v1/version/:vsn      # 获取单个资源
PUT    /api/v1/version/:vsn      # 更新资源
DELETE /api/v1/version/:vsn      # 删除资源
```

#### HTTP 方法

| 方法 | 用途 | 幂等性 |
|------|------|--------|
| GET | 获取资源 | ✅ |
| POST | 创建资源 | ❌ |
| PUT | 更新资源 | ✅ |
| DELETE | 删除资源 | ✅ |

### 响应格式

#### 成功响应

```json
{
  "data": {
    "vsn": "1.0.0",
    "subServers": {...}
  }
}
```

#### 错误响应

```json
{
  "error": {
    "code": "VERSION_NOT_FOUND",
    "message": "版本不存在",
    "details": {}
  }
}
```

---

## 代码提交规范

### Commit Message 格式

```
<type>: <subject>

<body>

<footer>
```

### Type 类型

- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具相关

### 示例

```
feat: 添加IP白名单批量管理功能

- 添加 BatchAddIPs 方法
- 添加 BatchRemoveIPs 方法
- 更新单元测试

Closes #123
```

---

## Checklist

### 代码提交前检查

- [ ] 代码符合命名规范
- [ ] 所有导出代码有中文注释
- [ ] 单元测试覆盖率 > 80%
- [ ] `go vet` 无警告
- [ ] `go test -race` 通过
- [ ] Commit message 符合规范

### PR 提交前检查

- [ ] 所有测试通过
- [ ] 代码审查完成
- [ ] 文档已更新
- [ ] 无合并冲突

---

**版本**: v2.0
**最后更新**: 2025-01-17

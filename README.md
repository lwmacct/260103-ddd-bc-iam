# Go DDD Package Library

基于领域驱动设计（DDD）和 CQRS 模式的可复用 Go 模块库，采用 **垂直切分的 Bounded Context 架构**。

## 特性

- **垂直切分架构**：按业务域组织模块（app/iam/crm），边界清晰
- **四层架构**：Domain → Application → Infrastructure → Transport
- **CQRS 分离**：Command/Query Repository 独立
- **依赖注入**：基于 Uber Fx
- **认证授权**：JWT + PAT 双重认证，URN 风格 RBAC
- **审计日志**：完整操作追踪
- **2FA 支持**：TOTP 双因素认证

## 技术栈

| 组件     | 技术           |
| -------- | -------------- |
| Web 框架 | Gin            |
| ORM      | GORM           |
| 数据库   | PostgreSQL     |
| 缓存     | Redis          |
| 依赖注入 | Uber Fx        |
| 配置管理 | cfgm           |
| API 文档 | Swagger (swag) |

## 架构概览

```
pkg/modules/                    # 业务模块（垂直切分）
├── app/                        # 核心治理域
│   ├── domain/                 #   领域层
│   ├── application/            #   应用层（UseCase Handlers）
│   ├── infrastructure/         #   基础设施层（Repository 实现）
│   └── transport/gin/          #   适配器层（HTTP Handler）
│
├── iam/                        # 身份管理域
│   ├── domain/
│   ├── application/
│   ├── infrastructure/         # IAM 专用基础设施（auth, twofa）
│   └── transport/gin/
│
├── crm/                        # CRM 域
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   └── transport/gin/
│
└── task/                       # 任务域
    ├── domain/
    ├── application/
    ├── infrastructure/
    └── transport/gin/

pkg/platform/                   # 平台层（跨模块技术能力）
├── cache/                      # Redis 客户端
├── db/                         # 数据库管理
├── eventbus/                   # 事件总线
├── health/                     # 健康检查
├── http/                       # HTTP 工具
├── queue/                      # Redis 队列
├── telemetry/                 # OpenTelemetry
└── validation/                 # JSON Logic 验证

pkg/shared/                     # 共享组件
├── cache/                      # 共享缓存
├── captcha/                    # 验证码
├── event/                      # 事件类型
└── health/                     # 健康检查

internal/
├── container/                  # ★ Fx 依赖注入组装点
├── bootstrap/                  # 应用启动引导
├── manualtest/                 # 集成测试
└── precommit/                  # 预提交钩子
```

**依赖方向**: `Transport → Application → Domain ← Infrastructure`

**垂直切分优势**：

- 按业务域组织，边界清晰
- 每个域包含完整的四层架构
- 共享技术基础设施在 `pkg/platform/`
- 便于独立演进和微服务化

## Bounded Context 划分

| Context | 说明           | 核心实体                           |
| ------- | -------------- | ---------------------------------- |
| `app`   | 核心治理域     | Setting, Audit, Org, Team, Task    |
| `iam`   | 身份认证与授权 | User, Role, Permission, PAT, TwoFA |
| `crm`   | 客户关系管理   | Lead, Opportunity, Contact         |
| `task`  | 任务管理域     | Task                               |

## 快速开始

### 运行示例服务器

```bash
# 确保依赖服务运行（PostgreSQL + Redis）
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:16
docker run -d -p 6379:6379 redis:alpine

# 初始化数据库
go run cmd/server/main.go db reset --force

# 启动服务
go run cmd/server/main.go
# 或使用热重载
air
```

**预置账号**: `admin / admin123`

### 在你的项目中使用

**步骤 1：复制 Container**

```bash
cp -r 260103-ddd-bc-iam/internal/container your-project/internal/
```

Container 文件结构：

| 文件            | 职责              | 修改方式       |
| --------------- | ----------------- | -------------- |
| `types.go`      | Model 列表、配置  | 添加你的 Model |
| `infra.go`      | DB/Redis/EventBus | 通常无需修改   |
| `cache.go`      | 缓存服务注册      | 添加缓存服务   |
| `service.go`    | JWT/TwoFA 服务    | 通常无需修改   |
| `http.go`       | Handler + 路由    | 添加你的路由   |
| `hooks.go`      | 生命周期钩子      | 通常无需修改   |
| `middleware.go` | 中间件注册        | 通常无需修改   |
| `router.go`     | 路由绑定          | 添加路由绑定   |

> **注**：`server.go`（HTTP 服务器）位于 `internal/bootstrap/` 目录

**步骤 2：添加自定义模块**

以 `Invoice` 模块为例，在项目中创建独立的域：

```go
// 1. 创建领域层 internal/domain/invoice/entity.go
type Invoice struct {
    ID      uint
    OrderID uint
    Amount  float64
    Status  string
}

// 2. 创建基础设施层 internal/infrastructure/persistence/invoice_model.go
type InvoiceModel struct {
    ID      uint    `gorm:"primaryKey"`
    OrderID uint    `gorm:"index;not null"`
    Amount  float64 `gorm:"type:decimal(10,2)"`
    Status  string  `gorm:"size:20"`
}

// 3. 创建应用层 internal/application/invoice/cmd_create.go
type CreateHandler struct {
    cmdRepo invoice.CommandRepository
}

// 4. 创建适配器层 internal/transport/gin/handler/invoice.go
type InvoiceHandler struct {
    createHandler *appInvoice.CreateHandler
}
```

**步骤 3：注册到 Container**

```go
// http.go - 添加 Handler
fx.Provide(
    // ...
    newInvoiceHandler,  // 你的 Handler
)
```

### 裁剪策略

不需要某些域时，直接从 `cmd/server/main.go` 移除对应的模块：

```go
fxOptions := []fx.Option{
    // Platform 层
    container.InfraModule,
    container.CacheModule,
    container.ServiceModule,

    // 业务模块 - 按需选择
    app.Module(),     // 核心治理域
    iam.Module(),     // 身份管理域
    // crm.Module(),  // 不需要 CRM，注释掉

    // HTTP 层
    container.HTTPModule,
    container.HooksModule,
}
```

## 开发命令

```bash
# 单元测试
go test ./...

# 编译检查
go build -o /dev/null ./...

# Lint 检查
golangci-lint run --new

# 数据库迁移
go run cmd/server/main.go db migrate

# 重置数据库
go run cmd/server/main.go db reset --force

# 手动集成测试
MANUAL=1 go test -v ./internal/manualtest/...
```

## API 文档

运行服务后访问 `/swagger/index.html`

## 参考示例

本库展示了完整的垂直切分 DDD 架构：

**业务域划分**：

- `pkg/modules/app/` - 核心域（组织、设置、任务、审计日志等）
- `pkg/modules/iam/` - 身份管理域（用户、认证、角色、PAT、2FA）
- `pkg/modules/crm/` - CRM 域（线索、商机、联系人、公司）
- `pkg/modules/task/` - 任务域（任务管理）

**依赖注入组装**：

- `internal/container/` - Fx 容器配置，展示如何组合所有域

**集成测试**：

- `internal/manualtest/` - HTTP API 集成测试，覆盖所有域

## Fx 模块结构

每个 Bounded Context 提供自包含的 Fx 模块：

```go
// pkg/modules/app/module.go
func Module() fx.Option {
    return fx.Module("app",
        infrastructure.PersistenceModule,  // 仓储注册
        application.UseCaseModule,         // 用例注册
        transport.HandlerModule,           // Handler 注册
    )
}
```

主程序组装：

```go
fx.New(
    fx.Supply(cfg),
    container.InfraModule,    // Platform: DB, Redis
    container.CacheModule,     // Cache services
    container.ServiceModule,   // JWT, TwoFA
    app.Module(),              // BC: App
    iam.Module(),              // BC: IAM
    crm.Module(),              // BC: CRM
    container.HTTPModule,      // HTTP: Routes
    container.HooksModule,     // Lifecycle
).Run()
```

## License

MIT

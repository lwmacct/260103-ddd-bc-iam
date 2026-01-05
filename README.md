# Go DDD Package Library

基于领域驱动设计（DDD）和 CQRS 模式的可复用 Go 模块库，采用 **垂直切分的 Bounded Context 架构**。

## 快速开始

### 运行示例服务器

```bash
# 依赖服务（PostgreSQL + Redis）
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

### 集成到你的项目

```bash
# 1. 复制 Container 配置
cp -r internal/container your-project/internal/

# 2. 在 main.go 中组装模块
fx.New(
    fx.Supply(cfg),
    container.InfraModule,     // Platform: DB, Redis
    container.CacheModule,      // Cache services
    container.ServiceModule,    // JWT, TwoFA
    iam.Module(),               // 你的业务模块
    container.HTTPModule,       // HTTP Routes
    container.HooksModule,      // Lifecycle
).Run()
```

> 📖 详细架构说明见 [`.claude/CLAUDE.md`](.claude/CLAUDE.md)

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
| API 文档 | Swagger (swag) |

## 架构概览

```
pkg/modules/                    # 业务模块（垂直切分）
├── app/                        # 核心治理域（设置、组织、审计）
├── iam/                        # 身份管理域（用户、认证、角色、PAT）
├── crm/                        # CRM 域（线索、商机、联系人）
└── task/                       # 任务域

pkg/platform/                   # 平台层（跨模块技术能力）
└── [db, redis, eventbus, http, ...]

internal/
└── container/                  # Fx 依赖注入组装点
```

| Bounded Context | 说明           | 核心实体                           |
| --------------- | -------------- | ---------------------------------- |
| `app`           | 核心治理域     | Setting, Audit, Org, Team, Task    |
| `iam`           | 身份认证与授权 | User, Role, Permission, PAT, TwoFA |
| `crm`           | 客户关系管理   | Lead, Opportunity, Contact         |
| `task`          | 任务管理域     | Task                               |

**依赖方向**: `Transport → Application → Domain ← Infrastructure`

> 📖 完整架构设计见 [`.claude/CLAUDE.md`](.claude/CLAUDE.md)

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

## License

MIT

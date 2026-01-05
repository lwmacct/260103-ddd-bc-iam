# IAM Bounded Context

基于领域驱动设计（DDD）和 CQRS 模式的 IAM（身份认证与授权）模块，采用 **垂直切分的 Bounded Context 架构**。

## 模块组成

本仓库包含两个关联模块：

- **IAM BC** (`pkg/modules/iam/`) - 身份认证与授权核心模块
- **Settings 封装** (`pkg/modules/settings/`) - 依赖外部 Settings BC 的用户/组织/团队配置层

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
# 1. 添加依赖
go get github.com/lwmacct/260103-ddd-bc-iam
go get github.com/lwmacct/260103-ddd-bc-settings  # Settings BC（被 IAM 依赖）
go get github.com/lwmacct/260103-ddd-shared       # Platform & Shared 层

# 2. 复制 Container 配置
cp -r internal/container your-project/internal/

# 3. 在 main.go 中组装模块
fx.New(
    fx.Supply(cfg),
    container.InfraModule,     // Platform: DB, Redis
    container.CacheModule,      // Cache services
    container.ServiceModule,    // JWT, TwoFA
    iam.Module(),               // IAM BC
    iamsettings.Module(),       // IAM 的 Settings 封装层
    container.HTTPModule,       // HTTP Routes
    container.HooksModule,      // Lifecycle
).Run()
```

> 📖 详细架构说明见 [`.claude/CLAUDE.md`](.claude/CLAUDE.md)

## 特性

- **垂直切分架构**：按业务域组织模块，边界清晰，可独立演化
- **四层架构**：Domain → Application → Infrastructure → Transport
- **依赖倒置**：Infrastructure 实现 Domain 接口，依赖方向单向可控
- **CQRS 分离**：Command/Query Repository 独立
- **依赖注入**：基于 Uber Fx
- **认证授权**：JWT + PAT 双重认证，URN 风格 RBAC
- **多租户支持**：组织/团队上下文动态注入，运行时变量解析
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
pkg/modules/
├── iam/                        # IAM Bounded Context（身份认证与授权）
│   ├── domain/                 # 领域层（实体、Repository 接口）
│   │   ├── user/               # 用户实体
│   │   ├── role/               # 角色与权限
│   │   ├── auth/               # 认证领域
│   │   ├── pat/                # 个人访问令牌
│   │   ├── twofa/              # 双因素认证
│   │   ├── org/                # 组织管理
│   │   └── audit/              # 审计日志
│   ├── app/                    # 应用层（UseCase Handler）
│   ├── infra/                  # 基础设施层（GORM、Redis、JWT）
│   └── adapters/gin/           # 适配器层（HTTP Handler + 路由）
│
└── settings/                   # Settings 封装层（跨 BC 依赖）
    ├── domain/                 # 用户/组织/团队配置实体
    │   ├── user/               # UserSetting 实体
    │   ├── org/                # OrgSetting 实体
    │   └── team/               # TeamSetting 实体
    ├── app/                    # 应用层（依赖外部 Settings BC 进行校验）
    ├── infra/                  # 基础设施层（持久化、缓存）
    └── adapters/gin/           # 适配器层

# 外部依赖 BC：
# - github.com/lwmacct/260103-ddd-bc-settings
#   └── 提供 Setting Schema 定义和校验逻辑
# - github.com/lwmacct/260103-ddd-shared
#   ├── platform/              # 纯技术基础设施（DB、Redis、EventBus、Queue、Telemetry）
#   └── shared/                # 接口定义层（Cache、Captcha、Event、Health）

internal/
└── container/                  # Fx 依赖注入组装点
```

**模块依赖关系**：

```
IAM BC
  ↓ 依赖
Settings 封装层
  ↓ 跨 BC 依赖
Settings BC (外部)
```

| 模块              | 职责                   | 核心实体                                                         |
| :---------------- | ---------------------- | ---------------------------------------------------------------- |
| **IAM BC**        | 身份认证与授权核心     | User, Role, Permission, PAT, TwoFA, Organization, Team, AuditLog |
| **Settings 封装** | 用户/组织/团队配置覆盖 | UserSetting, OrgSetting, TeamSetting                             |

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

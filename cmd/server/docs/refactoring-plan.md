# DDD 垂直切分重构计划

> **重构目标**：将水平分层架构重构为垂直切分的模块化架构，提升可复用性和可维护性。

## 一、当前架构分析

### 1.1 当前目录结构

```
github.com/lwmacct/260103-ddd-bc-iam/
├── cmd/server/              # 应用入口
├── config/                  # 配置文件
├── ddd/                     # 水平分层的 DDD 架构
│   ├── core/                # 核心治理 BC（混合了所有层）
│   │   ├── adapters/http/   # HTTP 适配层
│   │   ├── application/     # 应用层
│   │   ├── domain/          # 领域层
│   │   └── infrastructure/  # 基础设施层（含技术组件）
│   ├── iam/                # IAM BC
│   └── crm/                # CRM BC
├── internal/
│   ├── container/          # Uber Fx DI 容器
│   ├── manualtest/         # 集成测试
│   └── precommit/          # 预提交钩子
└── pkg/
    └── config/             # 配置管理
```

### 1.2 主要问题

| 问题                    | 影响                                                                              |
| ----------------------- | --------------------------------------------------------------------------------- |
| **技术组件散落在 core** | `database/`、`cache/`、`eventbus/` 等在 `ddd/core/infrastructure`，不利于外部复用 |
| **BC 边界不清晰**       | core 混合了通用技术和业务逻辑                                                     |
| **import 路径混乱**     | `ddd/core/application/setting` 和 `ddd/iam/application/auth` 层级不一致           |
| **不可复用**            | 外部项目无法选择性引入单个模块                                                    |

---

## 二、目标架构

### 2.1 新目录结构

```
github.com/lwmacct/260103-ddd-bc-iam/
├── cmd/server/                    # 应用入口
│   ├── main.go
│   └── docs/                      # Swagger 生成物
│
├── config/                        # 运行时配置
│
├── internal/app/                  # 应用装配层（不可复用）
│   ├── bootstrap/                 # Gin Engine / 中间件 / Server Lifecycle
│   │   ├── engine.go              #   创建 Gin Engine，注册全局中间件
│   │   ├── middleware.go          #   中间件工厂（CORS、Logger、Recovery）
│   │   └── server.go              #   HTTP Server 启停管理
│   │
│   ├── di/                        # 依赖注入装配（从 internal/container 演进）
│   │   ├── infra.go               #   Platform 组件 Provider
│   │   ├── module.go              #   BC 模块 Provider
│   │   ├── http.go                #   HTTP Handler + Router Provider
│   │   └── hooks.go               #   生命周期钩子
│   │
│   └── module/                    # Module 接口 + Registry
│       ├── interface.go           #   Module Interface (Start/Stop)
│       └── registry.go            #   模块注册表
│
├── pkg/                           # 可对外复用的主体
│   ├── config/                    # 配置管理（保持）
│   │
│   ├── platform/                  # 跨模块技术能力（可独立复用）
│   │   ├── db/                    #   GORM / Tx(UoW) / Migration / Seeder
│   │   │   ├── connection.go      #     数据库连接
│   │   │   ├── transaction.go     #     UoW 事务管理
│   │   │   ├── migration.go       #     迁移执行器
│   │   │   └── seeder.go          #     种子执行器
│   │   │
│   │   ├── cache/                 #   Redis Client + Cache 抽象
│   │   │   ├── client.go          #     Redis 客户端
│   │   │   ├── cache.go           #     Cache 接口
│   │   │   └── redis_cache.go     #     Redis 实现
│   │   │
│   │   ├── queue/                 #   Redis Queue/Processor
│   │   │   ├── queue.go           #     FIFO 队列
│   │   │   └── processor.go       #     消费处理器
│   │   │
│   │   ├── eventbus/              #   Memory Bus + 抽象
│   │   │   ├── bus.go             #     EventBus 接口
│   │   │   └── memory_bus.go      #     内存实现
│   │   │
│   │   ├── telemetry/             #   OpenTelemetry
│   │   │   └── otel.go            #     OTEL 初始化
│   │   │
│   │   ├── validation/            #   JSONLogic 校验器
│   │   │   └── validator.go       #     验证器实现
│   │   │
│   │   └── health/                #   健康检查
│   │       └── checker.go         #     健康检查器
│   │
│   ├── shared/                    # 共享（技术通用）
│   │   ├── errors/                #   通用错误码/错误包装
│   │   ├── utils/                 #   通用工具
│   │   └── kernel/                #   共享值对象（慎用）
│   │
│   ├── modules/                   # 业务模块（Bounded Context）
│   │
│   │   ├── core/                  # 核心治理 BC
│   │   │   ├── domain/            #     领域层
│   │   │   │   ├── audit/
│   │   │   │   ├── org/
│   │   │   │   ├── setting/
│   │   │   │   ├── stats/
│   │   │   │   └── task/
│   │   │   ├── application/       #     应用层
│   │   │   │   ├── audit/
│   │   │   │   ├── org/
│   │   │   │   ├── setting/
│   │   │   │   ├── stats/
│   │   │   │   └── task/
│   │   │   ├── infrastructure/    #     基础设施层
│   │   │   │   ├── persistence/   #       只放 core 自己的仓储实现
│   │   │   │   └── integration/   #       外部系统适配（可选）
│   │   │   ├── transport/         #     传输层适配器
│   │   │   │   └── gin/
│   │   │   │       ├── handler/   #         HTTP Handler
│   │   │   │       ├── routes/    #         路由元数据
│   │   │   │       └── routes.go  #         路由绑定
│   │   │   ├── migrations/        #     迁移脚本（可选）
│   │   │   └── module.go          #     对外组装入口
│   │   │
│   │   ├── iam/                   # IAM BC
│   │   │   ├── domain/
│   │   │   │   ├── auth/
│   │   │   │   ├── pat/
│   │   │   │   ├── role/
│   │   │   │   ├── twofa/
│   │   │   │   └── user/
│   │   │   ├── application/
│   │   │   │   ├── auth/
│   │   │   │   ├── pat/
│   │   │   │   ├── role/
│   │   │   │   ├── twofa/
│   │   │   │   └── user/
│   │   │   ├── infrastructure/
│   │   │   │   ├── auth/          #       JWT/密码哈希等实现
│   │   │   │   ├── twofa/         #       TOTP 实现
│   │   │   │   ├── persistence/   #       仓储实现
│   │   │   │   └── cache/         #       IAM 专属缓存
│   │   │   ├── transport/
│   │   │   │   └── gin/
│   │   │   ├── migrations/
│   │   │   └── module.go
│   │   │
│   │   └── crm/                   # CRM BC
│   │       ├── domain/
│   │       │   ├── company/
│   │       │   ├── contact/
│   │       │   ├── lead/
│   │       │   └── opportunity/
│   │       ├── application/
│   │       │   ├── company/
│   │       │   ├── contact/
│   │       │   ├── lead/
│   │       │   └── opportunity/
│   │       ├── infrastructure/
│   │       │   └── persistence/
│   │       ├── transport/
│   │       │   └── gin/
│   │       ├── migrations/
│   │       └── module.go
│   │
│   └── kit/                       # 对外门面（推荐引用入口）
│       ├── platform/              #   暴露 platform 构建器
│       │   ├── db.go
│       │   ├── cache.go
│       │   └── eventbus.go
│       └── modules/               #   暴露各模块 Provider
│           ├── core.go
│           ├── iam.go
│           └── crm.go
│
├── internal/manualtest/           # 集成测试（保持）
└── internal/precommit/            # 预提交钩子（保持）
```

### 2.2 关键变化

| 变化点       | 当前路径                           | 目标路径                        |
| ------------ | ---------------------------------- | ------------------------------- |
| **技术组件** | `ddd/core/infrastructure/database` | `pkg/platform/db`               |
| **BC 模块**  | `ddd/core`                         | `pkg/modules/app`               |
| **适配器层** | `ddd/core/adapters/http`           | `pkg/modules/app/transport/gin` |
| **DI 容器**  | `internal/container`               | `internal/app/di`               |

---

## 三、重构策略

### 3.1 核心原则

1. **保持编译通过**：每个 Phase 后确保 `go build` 成功
2. **增量迁移**：避免大爆炸式重写，分步迁移
3. **向后兼容**：保留旧路径的兼容层（可选）
4. **测试驱动**：每次迁移后运行测试验证

### 3.2 迁移优先级

```
低级别 → 高级别
Platform → BC Modules → Application Assembly
```

**理由**：高级别依赖低级别，先迁移底层可减少重复工作。

---

## 四、详细执行计划

### Phase 0: 准备工作 ✅

**目标**：建立依赖地图和测试基线

- [ ] 使用 `go list` 生成完整的依赖图
- [ ] 运行 `go test ./...` 记录测试基线
- [ ] 备份当前代码（打 tag）

### Phase 1: 创建骨架结构

**目标**：创建新目录结构，暂不迁移代码

```bash
# 创建 platform 骨架
mkdir -p pkg/platform/{db,cache,queue,eventbus,telemetry,validation,health}

# 创建 shared 骨架
mkdir -p pkg/shared/{errors,utils,kernel}

# 创建 modules 骨架
mkdir -p pkg/modules/{core,iam,crm}/{domain,application,infrastructure,transport}
mkdir -p pkg/modules/app/transport/gin
mkdir -p pkg/modules/iam/transport/gin
mkdir -p pkg/modules/crm/transport/gin

# 创建 kit 骨架
mkdir -p pkg/kit/{platform,modules}

# 创建 internal/app 骨架
mkdir -p internal/app/{bootstrap,di,module}
```

**产出**：空目录结构 + 占位符 `doc.go`

### Phase 2: 提取 Platform 层（高优先级）

**目标**：将 `ddd/core/infrastructure` 的技术组件迁移到 `pkg/platform`

| 子任务     | 迁移路径                                                         | 影响范围 |
| ---------- | ---------------------------------------------------------------- | -------- |
| Database   | `ddd/core/infrastructure/database` → `pkg/platform/db`           | 全局     |
| Cache      | `ddd/core/infrastructure/cache` → `pkg/platform/cache`           | 全局     |
| EventBus   | `ddd/core/infrastructure/eventbus` → `pkg/platform/eventbus`     | 全局     |
| Queue      | `ddd/core/infrastructure/queue` → `pkg/platform/queue`           | 全局     |
| Telemetry  | `ddd/core/infrastructure/telemetry` → `pkg/platform/telemetry`   | 全局     |
| Validation | `ddd/core/infrastructure/validation` → `pkg/platform/validation` | 全局     |
| Health     | `ddd/core/infrastructure/health` → `pkg/platform/health`         | 全局     |

**注意事项**：

1. **保留接口抽象**：`pkg/platform/cache/cache.go` 定义 Cache 接口
2. **Redis 实现分离**：`pkg/platform/cache/redis_cache.go` 实现
3. **import 路径更新**：全局替换 `github.com/lwmacct/260103-ddd-bc-iam/ddd/core/infrastructure/database` → `github.com/lwmacct/260103-ddd-shared/pkg/platform/db`

**验证**：

```bash
go build -o /dev/null ./...
go test ./pkg/platform/...
```

### Phase 3: 迁移 Core BC

**目标**：将 `ddd/core` 迁移到 `pkg/modules/app`

#### 3.1 Domain 层

```bash
mv ddd/core/domain/audit     pkg/modules/app/domain/audit
mv ddd/core/domain/org       pkg/modules/app/domain/org
mv ddd/core/domain/setting   pkg/modules/app/domain/setting
mv ddd/core/domain/stats     pkg/modules/app/domain/stats
mv ddd/core/domain/task      pkg/modules/app/domain/task
mv ddd/core/domain/cache     pkg/shared/cache      # 技术通用
mv ddd/core/domain/captcha   pkg/shared/captcha    # 技术通用
mv ddd/core/domain/health    pkg/shared/health     # 技术通用
mv ddd/core/domain/event     pkg/shared/event       # 跨 BC
```

#### 3.2 Application 层

```bash
mv ddd/core/application/audit    pkg/modules/app/application/audit
mv ddd/core/application/org      pkg/modules/app/application/org
mv ddd/core/application/setting  pkg/modules/app/application/setting
mv ddd/core/application/stats    pkg/modules/app/application/stats
mv ddd/core/application/task     pkg/modules/app/application/task
mv ddd/core/application/cache    pkg/modules/app/application/cache  # 保留
mv ddd/core/application/captcha  pkg/modules/app/application/captcha # 保留
mv ddd/core/application/health   pkg/modules/app/application/health  # 保留
```

#### 3.3 Infrastructure 层

```bash
# 只保留 BC 专属的仓储实现
mv ddd/core/infrastructure/persistence/*  pkg/modules/app/infrastructure/persistence/

# Seeds 归属
mv ddd/core/infrastructure/database/seeds  pkg/modules/app/migrations/
```

#### 3.4 Transport 层

```bash
# Adapters → Transport
mv ddd/core/adapters/http/handler/*  pkg/modules/app/transport/gin/handler/
mv ddd/core/adapters/http/routes/*   pkg/modules/app/transport/gin/routes/
mv ddd/core/adapters/http/middleware pkg/modules/app/transport/gin/middleware/
mv ddd/core/adapters/http/router.go   pkg/modules/app/transport/gin/
mv ddd/core/adapters/http/server.go   pkg/modules/app/transport/gin/
```

#### 3.5 Module Entry

```go
// pkg/modules/app/module.go
package core

import "go.uber.org/fx"

// Module 返回 Core BC 的 Fx 模块
func Module() fx.Option {
    return fx.Options(
        fx.Provide(
            // Repositories
            NewOrgRepositories,
            NewSettingRepositories,
            // ... 其他仓储

            // UseCases
            NewOrgUseCases,
            NewSettingUseCases,
            // ... 其他用例

            // Handlers
            NewOrgHandler,
            NewSettingHandler,
            // ... 其他 Handler
        ),
    )
}
```

### Phase 4: 迁移 IAM BC

**目标**：将 `ddd/iam` 迁移到 `pkg/modules/iam`

**关键点**：

1. **Infrastructure 拆分**：
   - `ddd/iam/infrastructure/auth` → `pkg/modules/iam/infrastructure/auth`
   - `ddd/iam/infrastructure/twofa` → `pkg/modules/iam/infrastructure/twofa`
   - Persistence 实现从 `ddd/core/infrastructure/persistence` 中提取 iam 相关部分

2. **Shared 依赖**：
   - `domain/event` 已在 Phase 3 迁移到 `pkg/shared/event`
   - 更新 import 路径

3. **Module Entry**：类似 Core BC 创建 `pkg/modules/iam/module.go`

### Phase 5: 迁移 CRM BC

**目标**：将 `ddd/crm` 迁移到 `pkg/modules/crm`

**相对简单**：CRM 当前只有 persistence，直接迁移即可。

### Phase 6: 重构 internal/container → internal/app/di

**目标**：将 DI 容器改为装配 `pkg/platform` 和 `pkg/modules`

#### 6.1 模块依赖关系

```
internal/app/di/
├── infra.go      # 装配 pkg/platform/*
├── module.go     # 装配 pkg/modules/*
├── http.go       # 装配 HTTP Transport
└── hooks.go      # 生命周期钩子
```

#### 6.2 示例：infra.go

```go
// internal/app/di/infra.go
package di

import (
    "github.com/lwmacct/260103-ddd-bc-iam/pkg/config"
    "github.com/lwmacct/260103-ddd-shared/pkg/platform/cache"
    "github.com/lwmacct/260103-ddd-shared/pkg/platform/db"
    "github.com/lwmacct/260103-ddd-shared/pkg/platform/eventbus"
    "github.com/lwmacct/260103-ddd-shared/pkg/platform/health"
    "github.com/lwmacct/260103-ddd-shared/pkg/platform/queue"
    "github.com/lwmacct/260103-ddd-shared/pkg/platform/telemetry"
    "github.com/lwmacct/260103-ddd-shared/pkg/platform/validation"
    "go.uber.org/fx"
)

var InfraModule = fx.Module("infra",
    fx.Provide(
        // Config
        config.New,

        // Platform
        telemetry.New,
        db.New,
        cache.New,
        eventbus.New,
        queue.New,
        validation.New,
        health.New,
    ),
)
```

#### 6.3 示例：module.go

```go
// internal/app/di/module.go
package di

import (
    "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/app"
    "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam"
    "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/crm"
    "go.uber.org/fx"
)

var ModuleModule = fx.Module("modules",
    fx.Options(
        core.Module(),
        iam.Module(),
        crm.Module(),
    ),
)
```

### Phase 7: 创建 internal/app/bootstrap

**目标**：将 HTTP Server 启动逻辑独立

```go
// internal/app/bootstrap/engine.go
package bootstrap

import (
    "github.com/gin-gonic/gin"
    "github.com/lwmacct/260103-ddd-shared/pkg/platform/telemetry"
)

// NewEngine 创建 Gin Engine，注册全局中间件
func NewEngine(tp *telemetry.OpenTelemetry) *gin.Engine {
    engine := gin.New()
    engine.Use(tp.Middleware())
    engine.Use(gin.Recovery())
    // ... 其他中间件
    return engine
}

// internal/app/bootstrap/server.go
package bootstrap

import (
    "context"
    "net/http"
    "time"
)

// Server 管理 HTTP Server 生命周期
type Server struct {
    engine *gin.Engine
    addr   string
}

func NewServer(engine *gin.Engine, cfg *config.Config) *Server {
    return &Server{
        engine: engine,
        addr:   cfg.HTTP.Addr,
    }
}

func (s *Server) Start() error {
    srv := &http.Server{
        Addr:    s.addr,
        Handler: s.engine,
    }
    return srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
    // Graceful shutdown
    return nil
}
```

### Phase 8: 创建 pkg/kit 对外门面

**目标**：提供简化的 API 供外部项目使用

```go
// pkg/kit/platform/db.go
package platform

import "github.com/lwmacct/260103-ddd-shared/pkg/platform/db"

// NewDatabase 创建数据库连接的便捷函数
func NewDatabase(dsn string) (*db.DB, error) {
    return db.New(dsn)
}

// pkg/kit/modules/app.go
package modules

import "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/app"

// CoreModule 返回 Core BC 的 Fx 模块
func CoreModule() core.Module {
    return core.Module()
}
```

### Phase 9: 更新 import 路径

**目标**：全局替换 import 路径

```bash
# 使用 sed 或 ide tools 进行全局替换

# Platform
sed -i 's|github.com/lwmacct/260103-ddd-bc-iam/ddd/core/infrastructure/database|github.com/lwmacct/260103-ddd-shared/pkg/platform/db|g' $(find . -name "*.go")

# Modules
sed -i 's|github.com/lwmacct/260103-ddd-bc-iam/ddd/core|github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/app|g' $(find . -name "*.go")
sed -i 's|github.com/lwmacct/260103-ddd-bc-iam/ddd/iam|github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam|g' $(find . -name "*.go")
sed -i 's|github.com/lwmacct/260103-ddd-bc-iam/ddd/crm|github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/crm|g' $(find . -name "*.go")

# Container → App
sed -i 's|github.com/lwmacct/260103-ddd-bc-iam/internal/container|github.com/lwmacct/260103-ddd-bc-iam/internal/app/di|g' $(find . -name "*.go")
```

### Phase 10: 验证与修复

**目标**：确保编译通过和测试通过

```bash
# 1. 编译检查
go build -o /dev/null ./...

# 2. 运行测试
go test ./...

# 3. 运行 linter
golangci-lint run --new

# 4. 集成测试
MANUAL=1 go test -v ./internal/manualtest/...
```

### Phase 11: 清理旧目录

**目标**：删除已迁移的旧目录

```bash
# 确认所有测试通过后
rm -rf ddd/core
rm -rf ddd/iam
rm -rf ddd/crm
rm -rf ddd  # 最后删除
rm -rf internal/container  # 迁移到 internal/app/di
```

---

## 五、风险评估与应对

| 风险                 | 影响            | 概率 | 应对措施                   |
| -------------------- | --------------- | ---- | -------------------------- |
| **Import 循环依赖**  | 编译失败        | 高   | 分步迁移，每步验证编译     |
| **测试用例路径错误** | 测试失败        | 中   | 同步更新测试 import        |
| **手动测试路径依赖** | manualtest 失败 | 中   | 更新 client.go 的 base URL |
| **遗漏依赖更新**     | 运行时 panic    | 中   | 使用 `go list` 检查依赖    |
| **Swagger 路径错误** | API 文档失效    | 低   | 更新 swagger 配置          |

---

## 六、验证清单

每完成一个 Phase 后执行：

- [ ] `go build -o /dev/null ./...`
- [ ] `go test ./...`
- [ ] `golangci-lint run --new`
- [ ] 检查 import 路径是否正确
- [ ] 检查文档注释是否更新

---

## 七、时间估算

| Phase    | 任务            | 预估时间 |
| -------- | --------------- | -------- |
| Phase 0  | 准备工作        | 0.5h     |
| Phase 1  | 创建骨架        | 0.5h     |
| Phase 2  | Platform 层     | 2h       |
| Phase 3  | Core BC         | 4h       |
| Phase 4  | IAM BC          | 3h       |
| Phase 5  | CRM BC          | 2h       |
| Phase 6  | internal/app/di | 2h       |
| Phase 7  | bootstrap       | 1h       |
| Phase 8  | pkg/kit         | 1h       |
| Phase 9  | 更新 import     | 1h       |
| Phase 10 | 验证修复        | 2h       |
| Phase 11 | 清理            | 0.5h     |
| **总计** |                 | **20h**  |

---

## 八、后续优化

1. **文档更新**：同步更新 README.md、CLAUDE.md
2. **示例更新**：更新 examples/ 目录
3. **规范更新**：更新 `.claude/rules/backend/` 规范文件
4. **CI/CD 调整**：更新 GitHub Actions 路径

---

## 附录：关键文件对照表

| 当前文件                                         | 目标文件                                  |
| ------------------------------------------------ | ----------------------------------------- |
| `ddd/core/infrastructure/database/connection.go` | `pkg/platform/db/connection.go`           |
| `ddd/core/infrastructure/cache/redis_client.go`  | `pkg/platform/cache/client.go`            |
| `ddd/core/infrastructure/eventbus/memory_bus.go` | `pkg/platform/eventbus/memory_bus.go`     |
| `ddd/core/adapters/http/server.go`               | `pkg/modules/app/transport/gin/server.go` |
| `internal/container/infra.go`                    | `internal/app/di/infra.go`                |
| `internal/container/usecase.go`                  | `internal/app/di/module.go`               |

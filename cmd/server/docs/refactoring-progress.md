# DDD 垂直切分重构进度

> **重构开始时间**：2026-01-02
> **预计完成时间**：TBD
> **当前状态**：规划阶段 ✅

---

## 总体进度

```
[████████████████████████████████████████████████████] 100% Phase 0: 规划 ✅
[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]   0% Phase 1: 骨架创建
[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]   0% Phase 2: Platform 层
[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]   0% Phase 3: Core BC
[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]   0% Phase 4: IAM BC
[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]   0% Phase 5: CRM BC
[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]   0% Phase 6: DI 容器
[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]   0% Phase 7: Bootstrap
[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]   0% Phase 8: Kit 门面
[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]   0% Phase 9: Import 更新
[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]   0% Phase 10: 验证
[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]   0% Phase 11: 清理
```

---

## Phase 0: 规划阶段 ✅

**目标**：分析当前架构，制定详细重构计划

- [x] 分析当前目录结构
- [x] 识别依赖关系
- [x] 创建重构计划文档 (`docs/refactoring-plan.md`)
- [x] 创建辅助脚本
  - [x] `scripts/migrate-scaffold.sh` - 骨架创建脚本
  - [x] `scripts/migrate-imports.sh` - Import 替换脚本
- [x] 创建进度跟踪文档 (`docs/refactoring-progress.md`)

**产出**：

- `docs/refactoring-plan.md` - 详细重构计划
- `docs/refactoring-progress.md` - 本文档
- `scripts/migrate-scaffold.sh` - 骨架创建脚本
- `scripts/migrate-imports.sh` - Import 替换脚本

**验证**：

- [x] 计划文档完整
- [x] 脚本可执行

---

## Phase 1: 创建骨架结构

**目标**：创建新目录结构骨架，暂不迁移代码

- [ ] 执行 `scripts/migrate-scaffold.sh`
- [ ] 验证目录结构正确
- [ ] 创建占位符 `doc.go`

**验证**：

- [ ] `ls -la pkg/platform/` - 7 个子目录
- [ ] `ls -la pkg/modules/{core,iam,crm}/` - 3 个 BC
- [ ] `ls -la internal/app/` - 3 个子目录

---

## Phase 2: 提取 Platform 层

**目标**：将 `ddd/core/infrastructure` 的技术组件迁移到 `pkg/platform`

### 子任务

- [ ] **Database** (`pkg/platform/db`)
  - [ ] 迁移 `ddd/core/infrastructure/database` → `pkg/platform/db`
  - [ ] 更新 import 路径
  - [ ] 验证编译

- [ ] **Cache** (`pkg/platform/cache`)
  - [ ] 迁移 `ddd/core/infrastructure/cache` → `pkg/platform/cache`
  - [ ] 更新 import 路径
  - [ ] 验证编译

- [ ] **EventBus** (`pkg/platform/eventbus`)
  - [ ] 迁移 `ddd/core/infrastructure/eventbus` → `pkg/platform/eventbus`
  - [ ] 更新 import 路径
  - [ ] 验证编译

- [ ] **Queue** (`pkg/platform/queue`)
  - [ ] 迁移 `ddd/core/infrastructure/queue` → `pkg/platform/queue`
  - [ ] 更新 import 路径
  - [ ] 验证编译

- [ ] **Telemetry** (`pkg/platform/telemetry`)
  - [ ] 迁移 `ddd/core/infrastructure/telemetry` → `pkg/platform/telemetry`
  - [ ] 更新 import 路径
  - [ ] 验证编译

- [ ] **Validation** (`pkg/platform/validation`)
  - [ ] 迁移 `ddd/core/infrastructure/validation` → `pkg/platform/validation`
  - [ ] 更新 import 路径
  - [ ] 验证编译

- [ ] **Health** (`pkg/platform/health`)
  - [ ] 迁移 `ddd/core/infrastructure/health` → `pkg/platform/health`
  - [ ] 更新 import 路径
  - [ ] 验证编译

**验证**：

- [ ] `go build -o /dev/null ./...`
- [ ] `go test ./pkg/platform/...`

---

## Phase 3: 迁移 Core BC

**目标**：将 `ddd/core` 迁移到 `pkg/modules/app`

### 3.1 Domain 层

- [ ] 迁移 `ddd/core/domain/audit` → `pkg/modules/app/domain/audit`
- [ ] 迁移 `ddd/core/domain/org` → `pkg/modules/app/domain/org`
- [ ] 迁移 `ddd/core/domain/setting` → `pkg/modules/app/domain/setting`
- [ ] 迁移 `ddd/core/domain/stats` → `pkg/modules/app/domain/stats`
- [ ] 迁移 `ddd/core/domain/task` → `pkg/modules/app/domain/task`
- [ ] 迁移 `ddd/core/domain/cache` → `pkg/shared/cache`
- [ ] 迁移 `ddd/core/domain/captcha` → `pkg/shared/captcha`
- [ ] 迁移 `ddd/core/domain/health` → `pkg/shared/health`
- [ ] 迁移 `ddd/core/domain/event` → `pkg/shared/event`

### 3.2 Application 层

- [ ] 迁移 `ddd/core/application/audit` → `pkg/modules/app/application/audit`
- [ ] 迁移 `ddd/core/application/org` → `pkg/modules/app/application/org`
- [ ] 迁移 `ddd/core/application/setting` → `pkg/modules/app/application/setting`
- [ ] 迁移 `ddd/core/application/stats` → `pkg/modules/app/application/stats`
- [ ] 迁移 `ddd/core/application/task` → `pkg/modules/app/application/task`
- [ ] 迁移 `ddd/core/application/cache` → `pkg/modules/app/application/cache`
- [ ] 迁移 `ddd/core/application/captcha` → `pkg/modules/app/application/captcha`
- [ ] 迁移 `ddd/core/application/health` → `pkg/modules/app/application/health`

### 3.3 Infrastructure 层

- [ ] 迁移 `ddd/core/infrastructure/persistence` → `pkg/modules/app/infrastructure/persistence`
- [ ] 迁移 `ddd/core/infrastructure/database/seeds` → `pkg/modules/app/migrations`

### 3.4 Transport 层

- [ ] 迁移 `ddd/core/adapters/http/handler` → `pkg/modules/app/transport/gin/handler`
- [ ] 迁移 `ddd/core/adapters/http/routes` → `pkg/modules/app/transport/gin/routes`
- [ ] 迁移 `ddd/core/adapters/http/middleware` → `pkg/modules/app/transport/gin/middleware`
- [ ] 迁移 `ddd/core/adapters/http/router.go` → `pkg/modules/app/transport/gin/router.go`
- [ ] 迁移 `ddd/core/adapters/http/server.go` → `pkg/modules/app/transport/gin/server.go`

### 3.5 Module Entry

- [ ] 创建 `pkg/modules/app/module.go`

**验证**：

- [ ] `go build -o /dev/null ./...`
- [ ] `go test ./pkg/modules/app/...`

---

## Phase 4: 迁移 IAM BC

**目标**：将 `ddd/iam` 迁移到 `pkg/modules/iam`

### 4.1 Domain 层

- [ ] 迁移 `ddd/iam/domain/auth` → `pkg/modules/iam/domain/auth`
- [ ] 迁移 `ddd/iam/domain/pat` → `pkg/modules/iam/domain/pat`
- [ ] 迁移 `ddd/iam/domain/role` → `pkg/modules/iam/domain/role`
- [ ] 迁移 `ddd/iam/domain/twofa` → `pkg/modules/iam/domain/twofa`
- [ ] 迁移 `ddd/iam/domain/user` → `pkg/modules/iam/domain/user`

### 4.2 Application 层

- [ ] 迁移 `ddd/iam/application/auth` → `pkg/modules/iam/application/auth`
- [ ] 迁移 `ddd/iam/application/pat` → `pkg/modules/iam/application/pat`
- [ ] 迁移 `ddd/iam/application/role` → `pkg/modules/iam/application/role`
- [ ] 迁移 `ddd/iam/application/twofa` → `pkg/modules/iam/application/twofa`
- [ ] 迁移 `ddd/iam/application/user` → `pkg/modules/iam/application/user`

### 4.3 Infrastructure 层

- [ ] 迁移 `ddd/iam/infrastructure/auth` → `pkg/modules/iam/infrastructure/auth`
- [ ] 迁移 `ddd/iam/infrastructure/twofa` → `pkg/modules/iam/infrastructure/twofa`
- [ ] 从 `ddd/core/infrastructure/persistence` 中提取 IAM 相关仓储 → `pkg/modules/iam/infrastructure/persistence`

### 4.4 Transport 层

- [ ] 迁移 `ddd/iam/adapters/http/handler` → `pkg/modules/iam/transport/gin/handler`

### 4.5 Module Entry

- [ ] 创建 `pkg/modules/iam/module.go`

**验证**：

- [ ] `go build -o /dev/null ./...`
- [ ] `go test ./pkg/modules/iam/...`

---

## Phase 5: 迁移 CRM BC

**目标**：将 `ddd/crm` 迁移到 `pkg/modules/crm`

### 5.1 Domain 层

- [ ] 迁移 `ddd/crm/domain/company` → `pkg/modules/crm/domain/company`
- [ ] 迁移 `ddd/crm/domain/contact` → `pkg/modules/crm/domain/contact`
- [ ] 迁移 `ddd/crm/domain/lead` → `pkg/modules/crm/domain/lead`
- [ ] 迁移 `ddd/crm/domain/opportunity` → `pkg/modules/crm/domain/opportunity`

### 5.2 Application 层

- [ ] 迁移 `ddd/crm/application/company` → `pkg/modules/crm/application/company`
- [ ] 迁移 `ddd/crm/application/contact` → `pkg/modules/crm/application/contact`
- [ ] 迁移 `ddd/crm/application/lead` → `pkg/modules/crm/application/lead`
- [ ] 迁移 `ddd/crm/application/opportunity` → `pkg/modules/crm/application/opportunity`

### 5.3 Infrastructure 层

- [ ] 迁移 `ddd/crm/infrastructure/persistence` → `pkg/modules/crm/infrastructure/persistence`

### 5.4 Module Entry

- [ ] 创建 `pkg/modules/crm/module.go`

**验证**：

- [ ] `go build -o /dev/null ./...`
- [ ] `go test ./pkg/modules/crm/...`

---

## Phase 6: 重构 DI 容器

**目标**：将 `internal/container` 改造为 `internal/app/di`

- [ ] 创建 `internal/app/di/infra.go` - 装配 `pkg/platform/*`
- [ ] 创建 `internal/app/di/module.go` - 装配 `pkg/modules/*`
- [ ] 创建 `internal/app/di/http.go` - 装配 HTTP Transport
- [ ] 迁移 `internal/container/hooks.go` → `internal/app/di/hooks.go`
- [ ] 更新 `cmd/server/main.go` 的引用

**验证**：

- [ ] `go build -o /dev/null ./...`
- [ ] `timeout 3 go run . api` - 验证 DI 成功

---

## Phase 7: 创建 Bootstrap

**目标**：将 HTTP Server 启动逻辑独立到 `internal/app/bootstrap`

- [ ] 创建 `internal/app/bootstrap/engine.go`
- [ ] 创建 `internal/app/bootstrap/middleware.go`
- [ ] 创建 `internal/app/bootstrap/server.go`

**验证**：

- [ ] `go build -o /dev/null ./...`
- [ ] `timeout 3 go run . api` - 验证服务器启动

---

## Phase 8: 创建 Kit 门面

**目标**：提供简化的 API 供外部项目使用

- [ ] 创建 `pkg/kit/platform/db.go`
- [ ] 创建 `pkg/kit/platform/cache.go`
- [ ] 创建 `pkg/kit/platform/eventbus.go`
- [ ] 创建 `pkg/kit/modules/app.go`
- [ ] 创建 `pkg/kit/modules/iam.go`
- [ ] 创建 `pkg/kit/modules/crm.go`

**验证**：

- [ ] `go doc ./pkg/kit/...`

---

## Phase 9: 更新 Import 路径

**目标**：全局批量替换 import 路径

- [ ] 执行 `scripts/migrate-imports.sh`
- [ ] 检查 `git diff` 确认替换正确
- [ ] 手动修复特殊路径

**验证**：

- [ ] `git diff --stat`
- [ ] `go build -o /dev/null ./...`

---

## Phase 10: 验证与修复

**目标**：确保编译通过和测试通过

- [ ] `go build -o /dev/null ./...`
- [ ] `go test ./...`
- [ ] `golangci-lint run --new`
- [ ] `MANUAL=1 go test -v ./internal/manualtest/...`

**修复清单**：

- [ ] 修复编译错误
- [ ] 修复测试失败
- [ ] 修复 Linter 警告
- [ ] 修复集成测试失败

---

## Phase 11: 清理旧目录

**目标**：删除已迁移的旧目录

- [ ] 删除 `ddd/core`
- [ ] 删除 `ddd/iam`
- [ ] 删除 `ddd/crm`
- [ ] 删除 `ddd` 目录
- [ ] 删除 `internal/container`

**验证**：

- [ ] `go build -o /dev/null ./...`
- [ ] `go test ./...`

---

## 问题跟踪

| ID  | 问题 | 状态 | 解决方案 |
| --- | ---- | ---- | -------- |
| -   | -    | -    | -        |

---

## 参考资源

- 📄 [重构计划详情](./refactoring-plan.md)
- 🛠️ [骨架创建脚本](../scripts/migrate-scaffold.sh)
- 🛠️ [Import 替换脚本](../scripts/migrate-imports.sh)

---

**最后更新**：2026-01-02

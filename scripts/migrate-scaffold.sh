#!/bin/bash
# DDD 垂直切分重构 - 目录结构骨架创建脚本
#
# 用途：快速创建新的目录结构骨架
# 使用：bash scripts/migrate-scaffold.sh

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

echo "🏗️  开始创建新目录结构骨架..."

# 创建 Platform 骨架
echo "📦 创建 pkg/platform/ ..."
mkdir -p pkg/platform/db
mkdir -p pkg/platform/cache
mkdir -p pkg/platform/queue
mkdir -p pkg/platform/eventbus
mkdir -p pkg/platform/telemetry
mkdir -p pkg/platform/validation
mkdir -p pkg/platform/health

# 创建 Shared 骨架
echo "📦 创建 pkg/shared/ ..."
mkdir -p pkg/shared/errors
mkdir -p pkg/shared/utils
mkdir -p pkg/shared/kernel
mkdir -p pkg/shared/cache    # 从 core/domain/cache 迁移
mkdir -p pkg/shared/captcha  # 从 core/domain/captcha 迁移
mkdir -p pkg/shared/health   # 从 core/domain/health 迁移
mkdir -p pkg/shared/event    # 从 core/domain/event 迁移

# 创建 Modules 骨架
echo "📦 创建 pkg/modules/ ..."

# Core BC
mkdir -p pkg/modules/app/domain/audit
mkdir -p pkg/modules/app/domain/org
mkdir -p pkg/modules/app/domain/setting
mkdir -p pkg/modules/app/domain/stats
mkdir -p pkg/modules/app/domain/task

mkdir -p pkg/modules/app/application/audit
mkdir -p pkg/modules/app/application/org
mkdir -p pkg/modules/app/application/setting
mkdir -p pkg/modules/app/application/stats
mkdir -p pkg/modules/app/application/task
mkdir -p pkg/modules/app/application/cache
mkdir -p pkg/modules/app/application/captcha
mkdir -p pkg/modules/app/application/health

mkdir -p pkg/modules/app/infrastructure/persistence
mkdir -p pkg/modules/app/infrastructure/integration

mkdir -p pkg/modules/app/transport/gin/handler
mkdir -p pkg/modules/app/transport/gin/routes
mkdir -p pkg/modules/app/transport/gin/middleware

mkdir -p pkg/modules/app/migrations

# IAM BC
mkdir -p pkg/modules/iam/domain/auth
mkdir -p pkg/modules/iam/domain/pat
mkdir -p pkg/modules/iam/domain/role
mkdir -p pkg/modules/iam/domain/twofa
mkdir -p pkg/modules/iam/domain/user

mkdir -p pkg/modules/iam/application/auth
mkdir -p pkg/modules/iam/application/pat
mkdir -p pkg/modules/iam/application/role
mkdir -p pkg/modules/iam/application/twofa
mkdir -p pkg/modules/iam/application/user

mkdir -p pkg/modules/iam/infrastructure/auth
mkdir -p pkg/modules/iam/infrastructure/twofa
mkdir -p pkg/modules/iam/infrastructure/persistence
mkdir -p pkg/modules/iam/infrastructure/cache

mkdir -p pkg/modules/iam/transport/gin/handler
mkdir -p pkg/modules/iam/transport/gin/routes
mkdir -p pkg/modules/iam/transport/gin/middleware

mkdir -p pkg/modules/iam/migrations

# CRM BC
mkdir -p pkg/modules/crm/domain/company
mkdir -p pkg/modules/crm/domain/contact
mkdir -p pkg/modules/crm/domain/lead
mkdir -p pkg/modules/crm/domain/opportunity

mkdir -p pkg/modules/crm/application/company
mkdir -p pkg/modules/crm/application/contact
mkdir -p pkg/modules/crm/application/lead
mkdir -p pkg/modules/crm/application/opportunity

mkdir -p pkg/modules/crm/infrastructure/persistence

mkdir -p pkg/modules/crm/transport/gin/handler
mkdir -p pkg/modules/crm/transport/gin/routes
mkdir -p pkg/modules/crm/transport/gin/middleware

mkdir -p pkg/modules/crm/migrations

# 创建 Kit 骨架
echo "📦 创建 pkg/kit/ ..."
mkdir -p pkg/kit/platform
mkdir -p pkg/kit/modules

# 创建 internal/app 骨架
echo "📦 创建 internal/app/ ..."
mkdir -p internal/app/bootstrap
mkdir -p internal/app/di
mkdir -p internal/app/module

# 创建占位符 doc.go
echo "📝 创建占位符 doc.go 文件..."

# Platform doc.go
cat > pkg/platform/db/doc.go <<EOF
// Package db 提供数据库连接、事务管理、迁移和种子数据功能。
package db
EOF

cat > pkg/platform/cache/doc.go <<EOF
// Package cache 提供 Redis 缓存抽象和实现。
package cache
EOF

cat > pkg/platform/eventbus/doc.go <<EOF
// Package eventbus 提供内存事件总线。
package eventbus
EOF

cat > pkg/platform/queue/doc.go <<EOF
// Package queue 提供 Redis FIFO 队列。
package queue
EOF

cat > pkg/platform/telemetry/doc.go <<EOF
// Package telemetry 提供 OpenTelemetry 链路追踪。
package telemetry
EOF

cat > pkg/platform/validation/doc.go <<EOF
// Package validation 提供 JSONLogic 验证器。
package validation
EOF

cat > pkg/platform/health/doc.go <<EOF
// Package health 提供健康检查功能。
package health
EOF

# Kit doc.go
cat > pkg/kit/platform/doc.go <<EOF
// Package platform 提供技术组件的便捷构建器。
package platform
EOF

cat > pkg/kit/modules/doc.go <<EOF
// Package modules 提供业务模块的便捷入口。
package modules
EOF

# internal/app doc.go
cat > internal/app/bootstrap/doc.go <<EOF
// Package bootstrap 提供 Gin Engine、中间件和 Server 生命周期管理。
package bootstrap
EOF

cat > internal/app/di/doc.go <<EOF
// Package di 提供依赖注入装配。
package di
EOF

cat > internal/app/module/doc.go <<EOF
// Package module 提供 Module 接口和注册表。
package module
EOF

echo "✅ 目录结构骨架创建完成！"
echo ""
echo "📊 统计信息："
echo "  - Platform 组件: 7 个"
echo "  - 业务模块 (BC): 3 个 (core, iam, crm)"
echo "  - Core 子模块: 5 个 (audit, org, setting, stats, task)"
echo "  - IAM 子模块: 5 个 (auth, pat, role, twofa, user)"
echo "  - CRM 子模块: 4 个 (company, contact, lead, opportunity)"
echo ""
echo "🚀 下一步：开始执行 Phase 2 - 提取 Platform 层技术组件"

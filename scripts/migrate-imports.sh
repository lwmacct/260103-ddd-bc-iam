#!/bin/bash
# DDD 垂直切分重构 - Import 路径替换脚本
#
# 用途：批量替换 import 路径
# 使用：bash scripts/migrate-imports.sh

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

echo "🔄 开始批量替换 import 路径..."
echo "⚠️  建议先提交当前更改，以便回滚"
echo ""

# 确认操作
read -p "是否继续？(y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 操作已取消"
    exit 1
fi

# 定义替换规则
declare -A replacements=(
    # Platform 层
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/core/infrastructure/database"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/platform/db"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/core/infrastructure/cache"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/platform/cache"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/core/infrastructure/eventbus"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/platform/eventbus"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/core/infrastructure/queue"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/platform/queue"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/core/infrastructure/telemetry"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/platform/telemetry"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/core/infrastructure/validation"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/platform/validation"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/core/infrastructure/health"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/platform/health"

    # Modules - Core
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/core/domain"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/app/domain"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/core/application"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/app/application"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/core/adapters"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/app/transport"

    # Modules - IAM
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/iam/domain"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/iam/application"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/application"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/iam/adapters"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/transport"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/iam/infrastructure/auth"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infrastructure/auth"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/iam/infrastructure/twofa"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infrastructure/twofa"

    # Modules - CRM
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/crm/domain"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/crm/domain"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/crm/application"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/crm/application"
    ["github.com/lwmacct/260103-ddd-bc-iam/ddd/crm/adapters"]="github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/crm/transport"

    # Container → DI
    ["github.com/lwmacct/260103-ddd-bc-iam/internal/container"]="github.com/lwmacct/260103-ddd-bc-iam/internal/app/di"
)

# 执行替换
total=0
for old_path in "${!replacements[@]}"; do
    new_path="${replacements[$old_path]}"
    echo "🔄 替换: $old_path"
    echo "   →     $new_path"

    # 使用 sed 进行替换（兼容 Linux 和 macOS）
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS (BSD sed)
        find . -name "*.go" -type f -exec sed -i '' "s|$old_path|$new_path|g" {} \;
    else
        # Linux (GNU sed)
        find . -name "*.go" -type f -exec sed -i "s|$old_path|$new_path|g" {} \;
    fi

    ((total++))
done

echo ""
echo "✅ 替换完成！共替换 $total 条路径规则"
echo ""
echo "📊 下一步："
echo "  1. 检查替换结果: git diff"
echo "  2. 编译验证: go build -o /dev/null ./..."
echo "  3. 提交更改: git add . && git commit -m 'refactor: 更新 import 路径'"

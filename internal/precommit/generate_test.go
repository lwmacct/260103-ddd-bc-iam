package precommit_test

import (
	"testing"
)

// TestGenerate_SwaggerAnnotations 生成 Swagger 注解。
//
// TODO: 此代码生成工具依赖旧的路由注册表系统，已在架构重构中移除。
//
// 旧架构依赖：
//   - routes.All() 获取所有 operation
//   - routes.Summary/Description/Tags/Path/Method() 访问元数据
//   - adapters/http/routes.go 中的路由绑定
//
// 新架构特点：
//   - 路由定义在 BC 模块的 routes/ 子包中（如 pkg/modules/iam/transport/gin/routes/）
//   - 路由元数据直接存储在 routes.Route 结构体字段中
//   - 无全局注册表，路由在 DI 时动态收集
//
// 重新设计选项：
//  1. 从 BC 模块的 routes/*.go 文件中静态分析提取路由元数据
//  2. 通过 AST 解析 routes.Route 字面量获取所有字段
//  3. 在集成测试中通过反射收集运行时路由信息
//  4. 使用 swag init 的自动生成功能，不依赖此工具
func TestGenerate_SwaggerAnnotations(t *testing.T) {
	t.Skip("TODO: 需要重新设计以适配新的路由架构（BC 模块化 + 无全局注册表）")
}

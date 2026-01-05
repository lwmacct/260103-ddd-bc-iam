package precommit_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAnnotation_MatchOperation 检查 handler @Router 注解与 operation 的一致性。
// 规则：每个 handler 的 @Router 路径必须与 operation 匹配。
//
// TODO: 此测试依赖旧的路由注册表系统，需要重新设计。
// 新架构中路由定义在 BC 模块的 routes/ 子包中，需要收集所有 BC 模块的路由后验证。
// 参考 router_test.go 中的 TODO 说明。
func TestAnnotation_MatchOperation(t *testing.T) {
	t.Skip("TODO: 需要重新设计以适配新的路由架构（BC 模块化 + 无全局注册表）")
}

// TestAnnotation_OperationCoverage 检查 operation 端点是否都有对应的 handler 注解。
// 规则：operation 中的每个端点都必须有带 @Router 注解的 handler。
//
// TODO: 此测试依赖旧的路由注册表系统，需要重新设计。
// 参考 router_test.go 中的 TODO 说明。
func TestAnnotation_OperationCoverage(t *testing.T) {
	t.Skip("TODO: 需要重新设计以适配新的路由架构（BC 模块化 + 无全局注册表）")
}

// TestAnnotation_RequiredFields 检查 Swagger 注解必填字段。
// 规则：每个 API 端点必须有 @Summary、@Tags、@Accept、@Produce。
func TestAnnotation_RequiredFields(t *testing.T) {
	annotations := parseHandlerAnnotations(t)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.NotEmpty(t, ann.Summary,
				"missing @Summary for %s %s", ann.Method, ann.Path)
			assert.NotEmpty(t, ann.Tags,
				"missing @Tags for %s %s", ann.Method, ann.Path)
			assert.NotEmpty(t, ann.Accept,
				"missing @Accept for %s %s", ann.Method, ann.Path)
			assert.NotEmpty(t, ann.Produce,
				"missing @Produce for %s %s", ann.Method, ann.Path)
		})
	}
}

// TestAnnotation_SecurityRequired 检查非公开端点的 @Security 注解。
// 规则：除公开端点外，所有 API 都必须有 @Security BearerAuth。
//
// TODO: 此测试依赖旧的路由注册表系统来获取公开端点列表，需要重新设计。
// 新架构中需要从 BC 模块的路由定义中提取 IsPublic 信息。
// 参考 router_test.go 中的 TODO 说明。
func TestAnnotation_SecurityRequired(t *testing.T) {
	t.Skip("TODO: 需要重新设计以适配新的路由架构（BC 模块化 + 无全局注册表）")
}

// TestAnnotation_TagsFormat 检查 @Tags 格式规范。
// 规则：格式为 kebab-case（小写字母+短横线），如「user-profile」「admin-role」「auth-2fa」。
func TestAnnotation_TagsFormat(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	// 匹配 kebab-case Tags 格式：小写字母开头，可包含小写字母、数字和短横线
	tagsRe := regexp.MustCompile(`^[a-z][a-z0-9\-]*$`)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") || ann.Tags == "" {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.True(t, tagsRe.MatchString(ann.Tags),
				"@Tags should be kebab-case (e.g. 'user-profile', 'admin-role'): got %q", ann.Tags)
		})
	}
}

// TestAnnotation_ContentType 检查 @Accept 和 @Produce 值。
// 规则：必须为 json。
func TestAnnotation_ContentType(t *testing.T) {
	annotations := parseHandlerAnnotations(t)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			if ann.Accept != "" {
				assert.Equal(t, "json", ann.Accept,
					"@Accept should be 'json': got %q", ann.Accept)
			}
			if ann.Produce != "" {
				assert.Equal(t, "json", ann.Produce,
					"@Produce should be 'json': got %q", ann.Produce)
			}
		})
	}
}

// TestAnnotation_SuccessDTOExists 检查 @Success 中的 DTO 类型是否存在。
// 规则：DTO 类型必须在 internal/application/{pkg}/dto.go 中定义。
// 例外：routes.* 类型（定义在 adapters/http/routes，是路由配置的一部分）。
func TestAnnotation_SuccessDTOExists(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	dtoTypes := loadDTOTypes(t)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") || ann.SuccessDTO == "" {
			continue
		}

		// 跳过 routes.* 和 registry.* 类型（路由层定义的元数据类型）
		if strings.HasPrefix(ann.SuccessDTO, "routes.") || strings.HasPrefix(ann.SuccessDTO, "registry.") {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.True(t, dtoTypes[ann.SuccessDTO],
				"@Success DTO type not found: %q\n  available types in package: check internal/application/{pkg}/dto.go",
				ann.SuccessDTO)
		})
	}
}

// TestAnnotation_ParamDTOExists 检查 @Param body 中引用 application 层 DTO 是否存在。
// 规则：带包前缀的 DTO（如 auth.LoginDTO）必须在 application 层定义。
// 注意：无包前缀的本地类型（如 CreateSettingRequest）不检查。
func TestAnnotation_ParamDTOExists(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	dtoTypes := loadDTOTypes(t)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") || ann.ParamDTO == "" {
			continue
		}

		// 只检查带包前缀的类型（如 auth.LoginDTO）
		// 跳过 handler 本地定义的类型（无 . 分隔符）
		if !strings.Contains(ann.ParamDTO, ".") {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.True(t, dtoTypes[ann.ParamDTO],
				"@Param body DTO type not found: %q\n  available types in package: check internal/application/{pkg}/dto.go",
				ann.ParamDTO)
		})
	}
}

// TestAnnotation_QueryTypeExists 检查 @Param query 中的结构体类型是否存在。
// 规则：Query 类型必须在 handler 包中定义。
func TestAnnotation_QueryTypeExists(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	queryTypes := loadHandlerQueryTypes(t)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") || ann.QueryType == "" {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.True(t, queryTypes[ann.QueryType],
				"@Param query type not found: %q\n  check handler file for type definition",
				ann.QueryType)
		})
	}
}

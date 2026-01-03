package precommit_test

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// ============================================================
// Types
// ============================================================

// structInfo 从文件中提取的结构体信息
type structInfo struct {
	File string
	Name string
}

// funcInfo 从文件中提取的函数信息
type funcInfo struct {
	File string
	Name string
}

// handlerAnnotation 从 handler 文件解析的注解信息
type handlerAnnotation struct {
	File        string
	Method      string // from @Router [method]
	Path        string // from @Router path
	Permission  string // from @x-permission {"scope":"..."}
	Summary     string // from @Summary
	Description string // from @Description
	Tags        string // from @Tags
	Security    string // from @Security
	Accept      string // from @Accept
	Produce     string // from @Produce
	SuccessDTO  string // from @Success, e.g., "user.UserWithRolesDTO"
	ParamDTO    string // from @Param body, e.g., "auth.LoginDTO"
	QueryType   string // from @Param query, e.g., "handler.ListUsersQuery"
}

// ============================================================
// Common Helpers
// ============================================================

// parseInterfaces 使用 AST 解析 Go 文件中的接口定义
func parseInterfaces(t *testing.T, filePath string) []structInfo {
	t.Helper()

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil
	}

	var interfaces []structInfo
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, isInterface := typeSpec.Type.(*ast.InterfaceType); isInterface {
				interfaces = append(interfaces, structInfo{
					File: filepath.Base(filePath),
					Name: typeSpec.Name.Name,
				})
			}
		}
	}
	return interfaces
}

// ============================================================
// Application Layer Helpers
// ============================================================

// parseStructs 使用 AST 解析 Go 文件中的结构体定义
func parseStructs(t *testing.T, filePath string) []structInfo {
	t.Helper()

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil
	}

	var structs []structInfo
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
				structs = append(structs, structInfo{
					File: filepath.Base(filePath),
					Name: typeSpec.Name.Name,
				})
			}
		}
	}
	return structs
}

// parseFuncs 使用 AST 解析 Go 文件中的顶级导出函数（不含方法）
func parseFuncs(t *testing.T, filePath string) []funcInfo {
	t.Helper()

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil
	}

	var funcs []funcInfo
	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// 跳过方法（有 receiver）
		if funcDecl.Recv != nil {
			continue
		}
		// 只收集导出函数
		if !funcDecl.Name.IsExported() {
			continue
		}
		funcs = append(funcs, funcInfo{
			File: filepath.Base(filePath),
			Name: funcDecl.Name.Name,
		})
	}
	return funcs
}

// getApplicationFiles 获取所有 BC 模块 application 目录下的 Go 文件
func getApplicationFiles(t *testing.T) []string {
	t.Helper()

	// 搜索所有 BC 模块的 application 目录
	appDirs := []string{
		"../../pkg/modules/app/application",
		"../../pkg/modules/iam/application",
		"../../pkg/modules/crm/application",
	}
	var files []string

	for _, appDir := range appDirs {
		// 跳过不存在的目录
		if _, err := os.Stat(appDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(appDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}
			// 跳过测试文件和 handler 文件
			if strings.HasSuffix(path, "_test.go") || strings.HasSuffix(path, "_handler.go") {
				return nil
			}
			files = append(files, path)
			return nil
		})

		if err != nil {
			t.Logf("warning: failed to walk application directory %s: %v", appDir, err)
		}
	}

	return files
}

// ============================================================
// Handler Annotation Helpers
// ============================================================

// parseHandlerAnnotations 解析所有 BC 模块 handler 目录下的 Go 文件的 Swagger 注解
func parseHandlerAnnotations(t *testing.T) []handlerAnnotation {
	t.Helper()

	// 搜索所有 BC 模块的 handler 目录
	handlerDirs := []string{
		"../../pkg/modules/app/transport/gin/handler",
		"../../pkg/modules/iam/transport/gin/handler",
		"../../pkg/modules/crm/transport/gin/handler",
	}
	var annotations []handlerAnnotation

	// 正则匹配
	routerRe := regexp.MustCompile(`@Router\s+(\S+)\s+\[(\w+)\]`)
	permRe := regexp.MustCompile(`@x-permission\s+\{"scope":"([^"]+)"\}`)
	summaryRe := regexp.MustCompile(`@Summary\s+(.+)$`)
	descRe := regexp.MustCompile(`@Description\s+(.+)$`)
	tagsRe := regexp.MustCompile(`@Tags\s+(.+)$`)
	securityRe := regexp.MustCompile(`@Security\s+(\S+)`)
	acceptRe := regexp.MustCompile(`@Accept\s+(\S+)`)
	produceRe := regexp.MustCompile(`@Produce\s+(\S+)`)
	// @Success 200 {object} response.DataResponse[user.UserDTO] "描述"
	// @Success 200 {object} response.DataResponse[[]menu.MenuDTO] "描述" (数组类型)
	// 提取泛型参数中的 DTO 类型，如 user.UserDTO 或 []menu.MenuDTO
	successRe := regexp.MustCompile(`@Success\s+\d+\s+\{object\}\s+response\.\w+\[(\[\])?([^\]]+)\]`)
	// @Param request body auth.LoginDTO true "登录凭证"
	// 提取 body 参数中的 DTO 类型，如 auth.LoginDTO
	paramBodyRe := regexp.MustCompile(`@Param\s+\S+\s+body\s+(\S+)\s+`)
	// @Param params query handler.ListUsersQuery false "查询参数"
	// 提取 query 参数中的结构体类型（带 handler. 前缀或本地类型）
	paramQueryRe := regexp.MustCompile(`@Param\s+\S+\s+query\s+(handler\.\w+|\w+Query)\s+`)

	for _, handlerDir := range handlerDirs {
		// 跳过不存在的目录（如 CRM 可能还没有 handlers）
		if _, err := os.Stat(handlerDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(handlerDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}
			if strings.HasSuffix(path, "_test.go") {
				return nil
			}

			file, err := os.Open(path) //nolint:gosec // 测试代码，路径来自 filepath.Walk
			if err != nil {
				return err
			}
			defer func() { _ = file.Close() }()

			var current handlerAnnotation
			current.File = filepath.Base(path)
			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				line := scanner.Text()

				// 解析各类注解
				if matches := routerRe.FindStringSubmatch(line); len(matches) == 3 {
					current.Path = matches[1]
					current.Method = strings.ToUpper(matches[2])
				}
				if matches := permRe.FindStringSubmatch(line); len(matches) == 2 {
					current.Permission = matches[1]
				}
				if matches := summaryRe.FindStringSubmatch(line); len(matches) == 2 {
					current.Summary = strings.TrimSpace(matches[1])
				}
				if matches := descRe.FindStringSubmatch(line); len(matches) == 2 {
					current.Description = strings.TrimSpace(matches[1])
				}
				if matches := tagsRe.FindStringSubmatch(line); len(matches) == 2 {
					current.Tags = strings.TrimSpace(matches[1])
				}
				if matches := securityRe.FindStringSubmatch(line); len(matches) == 2 {
					current.Security = matches[1]
				}
				if matches := acceptRe.FindStringSubmatch(line); len(matches) == 2 {
					current.Accept = matches[1]
				}
				if matches := produceRe.FindStringSubmatch(line); len(matches) == 2 {
					current.Produce = matches[1]
				}
				if matches := successRe.FindStringSubmatch(line); len(matches) == 3 {
					current.SuccessDTO = matches[2] // 第二组是实际 DTO 类型
				}
				if matches := paramBodyRe.FindStringSubmatch(line); len(matches) == 2 {
					current.ParamDTO = matches[1]
				}
				if matches := paramQueryRe.FindStringSubmatch(line); len(matches) == 2 {
					current.QueryType = matches[1]
				}

				// 遇到 func 定义，保存当前注解
				if strings.HasPrefix(strings.TrimSpace(line), "func ") && current.Path != "" {
					annotations = append(annotations, current)
					current = handlerAnnotation{File: filepath.Base(path)}
				}
			}

			return scanner.Err()
		})

		if err != nil {
			t.Logf("warning: failed to parse handler files in %s: %v", handlerDir, err)
		}
	}

	return annotations
}

// loadDTOTypes 使用 AST 解析所有 BC 模块 application 层的 DTO 类型
func loadDTOTypes(t *testing.T) map[string]bool {
	t.Helper()

	dtoTypes := make(map[string]bool)

	// 搜索所有 BC 模块的 application 目录
	appDirs := []string{
		"../../pkg/modules/app/application",
		"../../pkg/modules/iam/application",
		"../../pkg/modules/crm/application",
	}

	for _, appDir := range appDirs {
		// 跳过不存在的目录
		if _, err := os.Stat(appDir); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(appDir)
		if err != nil {
			t.Logf("warning: failed to read application directory %s: %v", appDir, err)
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			pkgName := entry.Name()
			dtoFile := filepath.Join(appDir, pkgName, "dto.go")

			// 跳过没有 dto.go 的包
			if _, err := os.Stat(dtoFile); os.IsNotExist(err) {
				continue
			}

			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, dtoFile, nil, 0)
			if err != nil {
				continue
			}

			for _, decl := range node.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
						if strings.HasSuffix(typeSpec.Name.Name, "DTO") {
							fullName := pkgName + "." + typeSpec.Name.Name
							dtoTypes[fullName] = true
						}
					}
				}
			}
		}
	}

	if len(dtoTypes) == 0 {
		t.Log("warning: no DTO types found in any BC module")
	}

	return dtoTypes
}

// loadHandlerQueryTypes 加载所有 BC 模块 handler 目录中定义的 Query 结构体类型
func loadHandlerQueryTypes(t *testing.T) map[string]bool {
	t.Helper()

	// 搜索所有 BC 模块的 handler 目录
	handlerDirs := []string{
		"../../pkg/modules/app/transport/gin/handler",
		"../../pkg/modules/iam/transport/gin/handler",
		"../../pkg/modules/crm/transport/gin/handler",
	}
	queryTypes := make(map[string]bool)

	for _, handlerDir := range handlerDirs {
		// 跳过不存在的目录
		if _, err := os.Stat(handlerDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(handlerDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
				return err
			}
			if strings.HasSuffix(path, "_test.go") {
				return nil
			}

			for _, s := range parseStructs(t, path) {
				if strings.HasSuffix(s.Name, "Query") {
					queryTypes[s.Name] = true
					queryTypes["handler."+s.Name] = true
				}
			}
			return nil
		})

		if err != nil {
			t.Logf("warning: failed to walk handler directory %s: %v", handlerDir, err)
		}
	}

	return queryTypes
}

// ============================================================
// Domain Layer Helpers
// ============================================================

// getDomainFiles 获取所有 BC 模块 domain 目录下的 Go 文件
func getDomainFiles(t *testing.T) []string {
	t.Helper()

	// 搜索所有 BC 模块的 domain 目录
	domainDirs := []string{
		"../../pkg/modules/app/domain",
		"../../pkg/modules/iam/domain",
		"../../pkg/modules/crm/domain",
	}
	var files []string

	for _, domainDir := range domainDirs {
		// 跳过不存在的目录
		if _, err := os.Stat(domainDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(domainDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}
			// 跳过测试文件
			if strings.HasSuffix(path, "_test.go") {
				return nil
			}
			files = append(files, path)
			return nil
		})

		if err != nil {
			t.Logf("warning: failed to walk domain directory %s: %v", domainDir, err)
		}
	}

	return files
}

package precommit_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestApplication_CommandSuffix 检查 commands.go 文件的结构体命名规范。
// 规则：所有结构体必须以 Command 结尾（如 CreateUserCommand）。
func TestApplication_CommandSuffix(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 commands.go 文件
		if filename != "commands.go" {
			continue
		}

		structs := parseStructs(t, file)
		for _, s := range structs {
			t.Run(s.File+"/"+s.Name, func(t *testing.T) {
				assert.True(t, strings.HasSuffix(s.Name, "Command"),
					"struct %q in %s should end with 'Command'", s.Name, s.File)
			})
		}
	}
}

// TestApplication_QuerySuffix 检查 queries.go 文件的结构体命名规范。
// 规则：所有结构体必须以 Query 结尾（如 GetUserQuery）。
func TestApplication_QuerySuffix(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 queries.go 文件
		if filename != "queries.go" {
			continue
		}

		structs := parseStructs(t, file)
		for _, s := range structs {
			t.Run(s.File+"/"+s.Name, func(t *testing.T) {
				assert.True(t, strings.HasSuffix(s.Name, "Query"),
					"struct %q in %s should end with 'Query'", s.Name, s.File)
			})
		}
	}
}

// TestApplication_HandlerSuffix 检查 cmd_*.go 和 qry_*.go 文件的结构体命名规范。
// 规则：每个文件至少包含一个以 Handler 结尾的结构体（如 CreateUserHandler）。
func TestApplication_HandlerSuffix(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 cmd_*.go 或 qry_*.go 文件
		if !strings.HasPrefix(filename, "cmd_") && !strings.HasPrefix(filename, "qry_") {
			continue
		}

		structs := parseStructs(t, file)
		hasHandler := false
		for _, s := range structs {
			if strings.HasSuffix(s.Name, "Handler") {
				hasHandler = true
				break
			}
		}

		t.Run(filename, func(t *testing.T) {
			assert.True(t, hasHandler,
				"file %s should contain at least one struct ending with 'Handler'", filename)
		})
	}
}

// TestApplication_DTOSuffix 检查 dto.go 文件的结构体命名规范。
// 规则：所有结构体必须以 DTO 结尾（如 UserDTO）。
func TestApplication_DTOSuffix(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 dto.go 文件
		if filename != "dto.go" {
			continue
		}

		structs := parseStructs(t, file)
		for _, s := range structs {
			t.Run(s.File+"/"+s.Name, func(t *testing.T) {
				assert.True(t, strings.HasSuffix(s.Name, "DTO"),
					"struct %q in %s should end with 'DTO'", s.Name, s.File)
			})
		}
	}
}

// TestApplication_MapperFuncNaming 检查 mapper.go 文件的函数命名规范。
// 规则：所有函数必须以 To 开头、DTO 或 DTOs 结尾（如 ToUserDTO、ToAuditActionDTOs）。
func TestApplication_MapperFuncNaming(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 mapper.go 文件
		if filename != "mapper.go" {
			continue
		}

		funcs := parseFuncs(t, file)
		for _, f := range funcs {
			t.Run(f.File+"/"+f.Name, func(t *testing.T) {
				assert.True(t, strings.HasPrefix(f.Name, "To"),
					"func %q in %s should start with 'To'", f.Name, f.File)
				validSuffix := strings.HasSuffix(f.Name, "DTO") || strings.HasSuffix(f.Name, "DTOs")
				assert.True(t, validSuffix,
					"func %q in %s should end with 'DTO' or 'DTOs'", f.Name, f.File)
			})
		}
	}
}

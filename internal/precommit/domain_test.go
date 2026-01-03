package precommit_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDomain_RepositoryInterfaceSuffix 检查 repository.go 文件的接口命名规范。
// 规则：所有接口必须以 Repository 结尾（如 CommandRepository, QueryRepository）。
func TestDomain_RepositoryInterfaceSuffix(t *testing.T) {
	files := getDomainFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 repository.go 文件
		if filename != "repository.go" {
			continue
		}

		interfaces := parseInterfaces(t, file)
		for _, iface := range interfaces {
			t.Run(iface.File+"/"+iface.Name, func(t *testing.T) {
				assert.True(t, strings.HasSuffix(iface.Name, "Repository"),
					"interface %q in %s should end with 'Repository'", iface.Name, iface.File)
			})
		}
	}
}

// TestDomain_CmdRepositorySuffix 检查 cmd_*.go 文件的接口命名规范。
// 规则：所有接口必须以 CommandRepository 结尾（如 SettingCategoryCommandRepository）。
func TestDomain_CmdRepositorySuffix(t *testing.T) {
	files := getDomainFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 cmd_*.go 文件
		if !strings.HasPrefix(filename, "cmd_") {
			continue
		}

		interfaces := parseInterfaces(t, file)
		for _, iface := range interfaces {
			t.Run(iface.File+"/"+iface.Name, func(t *testing.T) {
				assert.True(t, strings.HasSuffix(iface.Name, "CommandRepository"),
					"interface %q in %s should end with 'CommandRepository'", iface.Name, iface.File)
			})
		}
	}
}

// TestDomain_QryRepositorySuffix 检查 qry_*.go 文件的接口命名规范。
// 规则：所有接口必须以 QueryRepository 结尾（如 UserSettingQueryRepository）。
func TestDomain_QryRepositorySuffix(t *testing.T) {
	files := getDomainFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 qry_*.go 文件
		if !strings.HasPrefix(filename, "qry_") {
			continue
		}

		interfaces := parseInterfaces(t, file)
		for _, iface := range interfaces {
			t.Run(iface.File+"/"+iface.Name, func(t *testing.T) {
				assert.True(t, strings.HasSuffix(iface.Name, "QueryRepository"),
					"interface %q in %s should end with 'QueryRepository'", iface.Name, iface.File)
			})
		}
	}
}

package role_test

import (
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/role"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/permission"
)

// ExampleRole 演示如何创建和使用角色。
func ExampleRole() {
	// 创建一个编辑者角色
	editor := &role.Role{
		Name:        "editor",
		DisplayName: "编辑者",
		Description: "可以编辑文章的角色",
		IsSystem:    false,
	}

	// 添加权限
	editor.AddPermission(role.NewPermission("sys:posts:*", "*"))
	editor.AddPermission(role.NewPermission("sys:categories:*", "*"))

	fmt.Printf("角色: %s\n", editor.DisplayName)
	fmt.Printf("权限数量: %d\n", editor.GetPermissionCount())

	// Output:
	// 角色: 编辑者
	// 权限数量: 2
}

// ExampleRole_HasPermission 演示如何检查角色权限。
func ExampleRole_HasPermission() {
	admin := &role.Role{
		Name:        "admin",
		DisplayName: "管理员",
		IsSystem:    true,
	}

	// 添加通配符权限
	admin.AddPermission(role.NewPermission("sys:*", "*"))

	// 检查权限
	hasUserRead := admin.HasPermission(
		permission.Operation("sys:users:read"),
		permission.Resource("user.123"),
	)
	hasPostDelete := admin.HasPermission(
		permission.Operation("sys:posts:delete"),
		permission.Resource("post.456"),
	)

	fmt.Printf("用户读取权限: %v\n", hasUserRead)
	fmt.Printf("文章删除权限: %v\n", hasPostDelete)

	// Output:
	// 用户读取权限: true
	// 文章删除权限: true
}

// ExampleRole_HasOperationPermission 演示如何检查操作权限（不检查资源）。
func ExampleRole_HasOperationPermission() {
	viewer := &role.Role{
		Name:        "viewer",
		DisplayName: "查看者",
	}

	// 添加只读权限
	viewer.AddPermission(role.NewPermission("sys:posts:read", "*"))

	// 检查操作权限
	canRead := viewer.HasOperationPermission(permission.Operation("sys:posts:read"))
	canWrite := viewer.HasOperationPermission(permission.Operation("sys:posts:write"))

	fmt.Printf("可读取: %v\n", canRead)
	fmt.Printf("可写入: %v\n", canWrite)

	// Output:
	// 可读取: true
	// 可写入: false
}

// ExampleRole_AddPermission 演示如何向角色添加权限。
func ExampleRole_AddPermission() {
	editor := &role.Role{
		Name:        "editor",
		DisplayName: "编辑者",
	}

	// 添加多个权限
	editor.AddPermission(role.NewPermission("sys:posts:create", "*"))
	editor.AddPermission(role.NewPermission("sys:posts:update", "*"))
	editor.AddPermission(role.NewPermission("sys:posts:delete", "*"))

	// 尝试添加重复权限（会被忽略）
	editor.AddPermission(role.NewPermission("sys:posts:create", "*"))

	fmt.Printf("权限总数: %d\n", editor.GetPermissionCount())

	// Output:
	// 权限总数: 3
}

// ExampleRole_RemovePermission 演示如何从角色移除权限。
func ExampleRole_RemovePermission() {
	editor := &role.Role{
		Name:        "editor",
		DisplayName: "编辑者",
	}

	editor.AddPermission(role.NewPermission("sys:posts:create", "*"))
	editor.AddPermission(role.NewPermission("sys:posts:update", "*"))
	editor.AddPermission(role.NewPermission("sys:posts:delete", "*"))

	fmt.Printf("移除前权限数: %d\n", editor.GetPermissionCount())

	// 移除删除权限
	removed := editor.RemovePermission("sys:posts:delete", "*")
	fmt.Printf("移除成功: %v\n", removed)
	fmt.Printf("移除后权限数: %d\n", editor.GetPermissionCount())

	// Output:
	// 移除前权限数: 3
	// 移除成功: true
	// 移除后权限数: 2
}

// ExampleRole_ClearPermissions 演示如何清空角色权限。
func ExampleRole_ClearPermissions() {
	editor := &role.Role{
		Name:        "editor",
		DisplayName: "编辑者",
	}

	editor.AddPermission(role.NewPermission("sys:posts:*", "*"))
	editor.AddPermission(role.NewPermission("sys:categories:*", "*"))

	fmt.Printf("清空前权限数: %d\n", editor.GetPermissionCount())
	fmt.Printf("是否为空: %v\n", editor.IsEmpty())

	// 清空所有权限
	editor.ClearPermissions()

	fmt.Printf("清空后权限数: %d\n", editor.GetPermissionCount())
	fmt.Printf("是否为空: %v\n", editor.IsEmpty())

	// Output:
	// 清空前权限数: 2
	// 是否为空: false
	// 清空后权限数: 0
	// 是否为空: true
}

// ExampleRole_SetPermissions 演示如何批量设置角色权限。
func ExampleRole_SetPermissions() {
	editor := &role.Role{
		Name:        "editor",
		DisplayName: "编辑者",
	}

	// 初始权限
	editor.AddPermission(role.NewPermission("sys:posts:*", "*"))
	fmt.Printf("初始权限数: %d\n", editor.GetPermissionCount())

	// 批量设置新权限（替换现有权限）
	newPermissions := []role.Permission{
		role.NewPermission("sys:users:read", "*"),
		role.NewPermission("sys:users:update", "*"),
		role.NewPermission("sys:roles:read", "*"),
	}
	editor.SetPermissions(newPermissions)

	fmt.Printf("设置后权限数: %d\n", editor.GetPermissionCount())

	// Output:
	// 初始权限数: 1
	// 设置后权限数: 3
}

// ExampleRole_CanBeDeleted 演示如何检查角色是否可删除。
func ExampleRole_CanBeDeleted() {
	systemRole := &role.Role{
		Name:     "admin",
		IsSystem: true,
	}

	customRole := &role.Role{
		Name:     "editor",
		IsSystem: false,
	}

	fmt.Printf("系统角色可删除: %v\n", systemRole.CanBeDeleted())
	fmt.Printf("自定义角色可删除: %v\n", customRole.CanBeDeleted())

	// Output:
	// 系统角色可删除: false
	// 自定义角色可删除: true
}

// ExamplePermission 演示如何创建权限。
func ExamplePermission() {
	// 创建不同类型的权限
	fullAccess := role.NewPermission("sys:*", "*")
	userRead := role.NewPermission("sys:users:read", "*")
	postUpdate := role.NewPermission("sys:posts:update", "post.123")

	fmt.Printf("完全访问: %s / %s\n", fullAccess.OperationPattern, fullAccess.ResourcePattern)
	fmt.Printf("用户读取: %s / %s\n", userRead.OperationPattern, userRead.ResourcePattern)
	fmt.Printf("文章更新: %s / %s\n", postUpdate.OperationPattern, postUpdate.ResourcePattern)

	// Output:
	// 完全访问: sys:* / *
	// 用户读取: sys:users:read / *
	// 文章更新: sys:posts:update / post.123
}

// ExamplePermission_Matches 演示如何匹配权限。
func ExamplePermission_Matches() {
	perm := role.NewPermission("sys:users:*", "*")

	// 匹配测试
	matches1 := perm.Matches(
		permission.Operation("sys:users:read"),
		permission.Resource("user.123"),
	)
	matches2 := perm.Matches(
		permission.Operation("sys:posts:read"),
		permission.Resource("post.456"),
	)

	fmt.Printf("匹配用户操作: %v\n", matches1)
	fmt.Printf("匹配文章操作: %v\n", matches2)

	// Output:
	// 匹配用户操作: true
	// 匹配文章操作: false
}

// Package seeds 提供 IAM 模块的专属种子数据。
//
// 种子数据用于系统初始化，包括：
//   - RBACSeeder: 角色和权限数据
//   - UserSeeder: 默认管理员用户
//   - OrganizationSeeder: 默认组织
//
// 所有种子数据实现 platform/db.Seeder 接口。
package seeds

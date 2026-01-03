package role

// DefaultUserRoleName 默认用户角色名称。
//
// 所有已认证用户隐性拥有此角色，无需在数据库中显式分配。
// 权限通过数据库查询获取，由 Seeder 定义。
//
// 用途：
//   - 新注册用户自动拥有基本权限（如访问个人资料）
//   - 保持 RBAC 模型一致性（角色 → 权限）
const DefaultUserRoleName = "user"

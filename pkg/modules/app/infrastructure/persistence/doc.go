// Package persistence 提供领域仓储接口的 GORM 实现。
//
// 本包是 DDD 基础设施层的核心，实现了所有领域仓储接口：
//
// # CQRS 仓储实现
//
// 每个领域模块都有对应的 Command 和 Query Repository 实现：
//
// 用户模块：
//   - [UserCommandRepository]: 用户写操作（创建、更新、删除、角色分配）
//   - [UserQueryRepository]: 用户读操作（查询、搜索、统计）
//   - [UserRepositories]: 聚合结构，便于依赖注入
//
// 角色模块：
//   - [RoleCommandRepository]: 角色写操作
//   - [RoleQueryRepository]: 角色读操作
//
// 其他模块（PAT、Setting、TwoFA、AuditLog）遵循相同模式。
//
// # GORM Model
//
// 每个领域实体都有对应的 GORM Model（*_model.go）：
//   - [UserModel]: 用户持久化模型
//   - [RoleModel]: 角色持久化模型
//   - 其他 Model...
//
// Model 与 Entity 通过映射函数转换：
//   - newXxxModelFromEntity(): Entity -> Model
//   - (*XxxModel).ToEntity(): Model -> Entity
//
// # 泛型仓储基类
//
// [GenericCommandRepository] 提供 CRUD 通用实现，减少样板代码：
//
//	type GenericCommandRepository[E any, M Model[E]] struct {
//	    db            *gorm.DB
//	    entityToModel EntityToModel[E, M]
//	}
//
// # 依赖倒置
//
// 本包实现 domain 层定义的接口，遵循依赖倒置原则：
//   - domain/user.CommandRepository -> persistence.UserCommandRepository
//   - domain/user.QueryRepository -> persistence.UserQueryRepository
package persistence

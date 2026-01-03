// Package infrastructure 提供验证码基础设施实现。
//
// 本包提供：
//   - [MemoryRepository]: 验证码内存存储实现
//
// 线程安全性：
//   - MemoryRepository 是并发安全的，使用 sync.RWMutex 保护内部状态。
//
// 使用方式：
//
//	repo := infracaptcha.NewMemoryRepository()
//	defer repo.Close()
package infrastructure

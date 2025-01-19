package service

import (
	"app/backend/core"
)

// BaseService 是所有 service 的基础结构
type BaseService[T core.ServiceHandler] struct {
	core *core.App
	impl T // 要求必须存在 某些方法 方法
}

// NewBaseService 创建一个新的基础 service
func NewBaseService[T core.ServiceHandler](core *core.App, impl T) BaseService[T] {
	return BaseService[T]{
		core: core,
		impl: impl,
	}
}

// GetCoreApp 获取 App 实例
func (s *BaseService[T]) GetCoreApp() *core.App {
	return s.core
}

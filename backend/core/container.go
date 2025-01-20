package core

import (
	"net/http"

	"connectrpc.com/vanguard"
)

// ServiceHandler 要求所有的 RPC 服务都实现这个接口
type ServiceHandler interface {
	GetHandle() (string, http.Handler)
}

// Container 用于注册和获取所有的 RPC 服务
type Container struct {
	services []ServiceHandler
}

func NewContainer() *Container {
	return &Container{}
}

// RegisterService 注册一个 RPC 服务
func (s *Container) RegisterService(svc ServiceHandler) {
	s.services = append(s.services, svc)
}

// GetAllServices 获取所有的 RPC 服务
func (s *Container) GetAllServices() []ServiceHandler {
	return s.services
}

// BuildServices 获取所有的 Vanguard 服务
func (s *Container) BuildServices() []*vanguard.Service {
	services := make([]*vanguard.Service, len(s.services))
	for i, svc := range s.services {
		services[i] = vanguard.NewService(svc.GetHandle())
	}
	return services
}

// Vanguard Transcoder
func (s *Container) BuildTranscoder() (*vanguard.Transcoder, error) {
	return vanguard.NewTranscoder(s.BuildServices())
}

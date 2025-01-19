package service

import (
	v1 "app/backend/pkg/gen/v1"
	"app/backend/pkg/gen/v1/v1connect"
	"context"
	"net/http"
	"strconv"

	"app/backend/core"

	"connectrpc.com/connect"
)

type TestService struct {
	v1connect.UnimplementedTestServiceHandler
	BaseService[*TestService]
}

// NewTestService 此方法不是必须，可以在 backend/app.go 中初始化。写到这里是因为大多数情况下 Service 所需的依赖都会在 core 里被定义
func NewTestService(core *core.App) *TestService {
	svc := &TestService{}
	svc.BaseService = NewBaseService(core, svc)
	return svc
}

func (s *TestService) GetHandle() (string, http.Handler) {
	return v1connect.NewTestServiceHandler(s)
}

func (s *TestService) Test1(ctx context.Context, req *connect.Request[v1.Test1Request]) (*connect.Response[v1.Test1Response], error) {
	return connect.NewResponse(&v1.Test1Response{
		Message: "Hello No Args",
	}), nil
}

func (s *TestService) Test2(ctx context.Context, req *connect.Request[v1.Test2Request]) (*connect.Response[v1.Test2Response], error) {
	return connect.NewResponse(&v1.Test2Response{
		Message: "Hello " + req.Msg.Name,
	}), nil
}

func (s *TestService) Test3(ctx context.Context, req *connect.Request[v1.Test3Request], stream *connect.ServerStream[v1.Test3Response]) error {
	for i := 0; i < 10; i++ {
		err := stream.Send(&v1.Test3Response{
			Message: "Hello " + strconv.Itoa(i),
		})
		if err != nil {
			break
		}
	}
	return nil
}

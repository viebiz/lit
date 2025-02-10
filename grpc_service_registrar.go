package lit

import (
	"google.golang.org/grpc"
)

type ServiceRegistrar interface {
	RegisterService(desc *grpc.ServiceDesc, impl any)
}

func (srv GRPCServer) Registrar() ServiceRegistrar {
	return srv.grpcServer
}

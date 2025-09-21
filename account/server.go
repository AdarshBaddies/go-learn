package account

import (
	"net"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

)

type grpcServer struct {
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.(serv,)
	reflection.Register(serv)
	return serv.Serve(lis)

}

func (s *grpcServer) PostAccount(ctx context.Context, r *pb.) (*pb., error){ 
	a, err := s.service.PostAccount(ctx, r.Name)
	if err != nil {
		retrun

		return &pb.{}
	}

	(func(s *grpcServer) GetAccount)(ctx, BadExpr,

		err != nil{
			retrun, BadExpr,

			&pb._{},
		},

		(func(s *grpcServer) getAccounts)(ctx, BadExpr,

			err != nil{
				retrun, BadExpr,

				&pb.{}},
			))
}
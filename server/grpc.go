package server

import (
	"context"
	"fmt"
	"net"

	pb "github.com/u1f35c/grpc-test/testservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedTestServiceServer
}

func (s *server) GetStatus(ctx context.Context, _ *emptypb.Empty) (*pb.StatusReply, error) {
	ss := ctx.Value(statusStoreCtxKey{})
	if ss == nil {
		return nil, fmt.Errorf("no status context")
	}

	value := ss.(*StatusStore).GetStatus()

	return &pb.StatusReply{Value: int32(value)}, nil
}

func (s *server) SetStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusReply, error) {
	ss := ctx.Value(statusStoreCtxKey{})
	if ss == nil {
		return nil, fmt.Errorf("no status context")
	}

	value := ss.(*StatusStore).SetStatus(int(req.Value))

	return &pb.StatusReply{Value: int32(value)}, nil
}

// statsHandler implements the GRPC stats.Handler interface, allowing us to be
// called when an incoming connection is established and attach a suitable
// StatusStore to the context.
type statsHandler struct{}

func (st *statsHandler) TagConn(ctx context.Context, stat *stats.ConnTagInfo) context.Context {
	return context.WithValue(ctx, statusStoreCtxKey{}, &StatusStore{})
}

func (st *statsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
}

func (st *statsHandler) TagRPC(ctx context.Context, stat *stats.RPCTagInfo) context.Context {
	return ctx
}

func (st *statsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
}

// GRPCServe starts a GRPC server on the supplied port
func GRPCServe(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	h := &statsHandler{}
	s := grpc.NewServer(grpc.StatsHandler(h))
	pb.RegisterTestServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}

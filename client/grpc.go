package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/u1f35c/grpc-test/testservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCConnect connects to a GRPC server on the given port
func GRPCConnect(port int) error {
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	c := pb.NewTestServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetStatus(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Status: %d\n", r.Value)

	r, err = c.SetStatus(ctx, &pb.StatusRequest{Value: 2})
	if err != nil {
		return err
	}

	fmt.Printf("Status: %d\n", r.Value)

	r, err = c.GetStatus(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Status: %d\n", r.Value)

	return nil
}

package client

import (
	"context"

	pb "github.com/filipio/snf-knative/src/cache/grpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Call(addr string, message string) (string, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pb.NewHandlerClient(conn)
	reply, err := client.Handle(context.Background(), &pb.HandleRequest{Message: message}, grpc.EmptyCallOption{})
	if err != nil {
		return "", err
	}

	return reply.Message, nil
}

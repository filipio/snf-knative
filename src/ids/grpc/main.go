package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/filipio/snf-knative/src/ids/grpc/client"
	pb "github.com/filipio/snf-knative/src/ids/grpc/protos"
	"google.golang.org/grpc"
)

type HandleService struct {
	pb.UnimplementedHandlerServer
}

func (s *HandleService) Handle(ctx context.Context, in *pb.HandleRequest) (*pb.HandleReply, error) {
	functionName := os.Getenv("FUNCTION_NAME")
	if functionName == "" {
		return nil, fmt.Errorf("functionName variable not set")
	}

	successors := os.Getenv("SUCCESSORS")

	// split it by comma
	successorsNames := strings.Split(successors, ",")
	fmt.Println("functionName: ", functionName, " successors: ", successorsNames)

	requestBody := "ping"

	if len(successorsNames) > 0 {
		for _, successor := range successorsNames {
			if successor == "" {
				continue
			}

			successorAddr := successor + ":80"
			if successorAddr[:7] == "http://" {
				successorAddr = successorAddr[7:]
			}
			fmt.Printf("sending request to '%s'", successorAddr)

			response, err := client.Call(successorAddr, requestBody)
			if err != nil {
				return nil, err
			}

			fmt.Println("Response from ", successor, ": ", response)
		}
	}

	return &pb.HandleReply{Message: "ok!"}, nil
}

func NewHandleService() *HandleService {
	return &HandleService{}
}

func main() {
	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterHandlerServer(grpcServer, NewHandleService())
	grpcServer.Serve(lis)
}

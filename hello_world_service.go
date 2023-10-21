package main

import (
	"context"

	pb "github.com/1005281342/trpc-demo/helloworld"
)

type helloWorldServiceImpl struct {
	pb.UnimplementedHelloWorldService
}

// Hello Hello says hello.
func (s *helloWorldServiceImpl) Hello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloResponse, error) {
	rsp := &pb.HelloResponse{
		Msg: req.GetMsg(),
	}
	return rsp, nil
}

package main

import (
	"context"
	"ds-rpc/proto"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type server struct {
}

func main() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterAddServiceServer(srv, &server{})
	reflection.Register(srv)

	if err = srv.Serve(listener); err != nil {
		panic(err)
	}
}

func (s *server) PerformanceReport(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	cpu := request.GetCPU()
	usageRAM, avaliableRam := request.GetUsedRAM(), request.GetAvaliableRAM()
	usedDisk, avaliableDisk :=request.GetUsedDisk(), request.GetAvaliableDisk()

	fmt.Println(cpu, usageRAM, avaliableRam, usedDisk, avaliableDisk)

	result := true

	return &proto.Response{Success: result}, nil
}

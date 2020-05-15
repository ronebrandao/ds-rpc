package main

import (
	"context"
	"ds-rpc/proto"
	"ds-rpc/server/db"
	"ds-rpc/server/model"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type server struct {
}

var storage db.DB

func main() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterAddServiceServer(srv, &server{})
	reflection.Register(srv)

	storage = db.DB{}

	err = storage.OpenConnection(context.Background())
	if err != nil {
		fmt.Println("*** It was not possible to connect to the database! ***")
	}

	fmt.Println("Connected to database...")

	defer storage.CloseConnection()

	if err = srv.Serve(listener); err != nil {
		panic(err)
	}
}

func (s *server) PerformanceReport(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	cpu := request.GetCPU()
	usedRAM, avaliableRAM := request.GetUsedRAM(), request.GetAvaliableRAM()
	usedDisk, avaliableDisk :=request.GetUsedDisk(), request.GetAvaliableDisk()

	fmt.Printf("CPU: %.2f %% | Used RAM: %f MB | Avaliable RAM: %f MB | Used Disk: %.2f %% | Avaliable Disk: %.2f %% \n",
		cpu, usedRAM / 1000000, avaliableRAM / 1000000, usedDisk, avaliableDisk)

	if err := storage.SavePerformanceStats(model.Status{
		CPU:           cpu,
		UsedRAM:       usedRAM,
		AvaliableRAM:  avaliableRAM,
		UsedDisk:      usedDisk,
		AvaliableDisk: avaliableDisk,
	}); err != nil {
			fmt.Println("***** It was not possible to save status to database. *****", err)
			return &proto.Response{Success: false}, nil
	}

	return &proto.Response{Success: true}, nil
}

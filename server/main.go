package main

import (
	"context"
	"ds-rpc/proto"
	"ds-rpc/server/db"
	"ds-rpc/server/model"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Connection struct {
	stream     proto.Broadcast_CreateStreamServer
	id         string
	clientname string
	active     bool
	error      chan error
}

type Server struct {
	Connection []*Connection
}

var storage db.DB

func main() {
	var connections []*Connection

	srv := grpc.NewServer()
	server := Server{Connection: connections}
	proto.RegisterBroadcastServer(srv, &server)

	reflection.Register(srv)

	storage = db.DB{}

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	if err = srv.Serve(listener); err != nil {
		panic(err)
	}
}

func (s *Server) PerformanceReport(ctx context.Context, perfReport *proto.PerfReport) (*proto.Response, error) {
	cpu := perfReport.Message.GetCPU()
	usedRAM, avaliableRAM := perfReport.Message.GetUsedRAM(), perfReport.Message.GetAvaliableRAM()
	usedDisk, avaliableDisk := perfReport.Message.GetUsedDisk(), perfReport.Message.GetAvaliableDisk()

	fmt.Printf("ID: %s  |  CPU: %.2f %% | Used RAM: %f MB | Avaliable RAM: %f MB | Used Disk: %.2f %% | Avaliable Disk: %.2f %% \n",
		perfReport.Client.Id, cpu, usedRAM/1000000, avaliableRAM/1000000, usedDisk, avaliableDisk)

	err := storage.OpenConnection(context.Background())
	if err != nil {
		fmt.Println("*** It was not possible to connect to the database! ***")
	} else {
		fmt.Println("Connected to database...")
	}

	defer storage.CloseConnection()

	err = storage.SavePerformanceStats(model.Status{
		CPU:           cpu,
		UsedRAM:       usedRAM,
		AvaliableRAM:  avaliableRAM,
		UsedDisk:      usedDisk,
		AvaliableDisk: avaliableDisk,
	})
	if err != nil {
		fmt.Println("***** It was not possible to save status to database. *****", err)
		return &proto.Response{Success: false}, nil
	}

	return &proto.Response{Success: true}, nil
}

func (s *Server) CreateStream(pconn *proto.Connect, stream proto.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream:     stream,
		id:         pconn.Client.Id,
		clientname: pconn.Client.Name,
		active:     true,
		error:      make(chan error),
	}
	s.Connection = append(s.Connection, conn)

	s.RequestInfo(context.Background(), pconn.Client)

	return <-conn.error
}

func (s *Server) RequestInfo(ctx context.Context, msg *proto.Client) (*proto.Close, error) {

	for _, conn := range s.Connection {
		go func() {
			gocron.Every(15).Seconds().Do(func() {
				fmt.Println("!-- SERVER ASKING FOR INFO --!")
				fmt.Printf("Client-name %s\n", conn.clientname)

				conn.stream.Send(&proto.ServerRequest{
					SendInfo: true,
				})
			})
		}()
	}

	<-gocron.Start()

	return &proto.Close{}, nil
}

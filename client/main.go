package main

import (
	"context"
	"crypto/sha256"
	"ds-rpc/proto"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"google.golang.org/grpc"
	"sync"
	"time"
)

var client proto.BroadcastClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	fmt.Println("!-- Information gattering starting in 10 seconds --!")

	timestamp := time.Now()

	name := flag.String("N", "Anon", "The name of the user")
	flag.Parse()

	id := sha256.Sum256([]byte(timestamp.String() + *name))

	client = proto.NewBroadcastClient(conn)
	user := &proto.Client{
		Id:   hex.EncodeToString(id[:]),
		Name: *name,
	}

	err = connect(user)
	if err != nil {
		fmt.Println("Error connecting user.")
	}

	ctx := context.Background()
	// rotina para enviar informacoes para o servidor a cada 10 seg
	go func() {
		gocron.Every(10).Seconds().Do(func() {
			fmt.Println("!-- Begining information gattering --!")

			sendReport(ctx, user)
		})
	}()

	<-gocron.Start()

}

func sendReport(ctx context.Context, user *proto.Client) {

	cpuUsage, err := getCPU()
	if err != nil {
		fmt.Println("*** It was not possible to obtain the CPU usage ***")
	}

	memoryStats, err := getMemoryStats()
	if err != nil {
		fmt.Println("*** It was not possible to memory the CPU usage ***")
	}

	if memoryStats == nil {
		memoryStats = &memory.Stats{}
	}

	req := &proto.PerfReport{
		Client: user,
		Message: &proto.Report{
			MsgId:         "abc",
			CPU:           cpuUsage,
			UsedRAM:       float32(memoryStats.Used),
			AvaliableRAM:  float32(memoryStats.Free),
			UsedDisk:      0.0,
			AvaliableDisk: 0.0,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if response, err := client.PerformanceReport(ctx, req); err == nil {
		printResponse(response)
	} else {
		printResponse(&proto.Response{Success: false})
	}
}

func printResponse(response *proto.Response) {
	if response.Success {
		fmt.Println("--- Performance status sucessefully saved! ---")
	} else {
		fmt.Println("*** Error on trying to save performance status! ***")
	}
}

func getCPU() (float32, error) {
	before, err := cpu.Get()
	if err != nil {
		return 0.0, nil
	}
	time.Sleep(time.Duration(1) * time.Second)

	after, err := cpu.Get()
	if err != nil {
		return 0.0, nil
	}

	total := float64(after.Total - before.Total)

	user := float64(after.User-before.User) / total * 100
	system := float64(after.System-before.System) / total * 100

	return float32(user + system), nil
}

func getMemoryStats() (*memory.Stats, error) {
	memory, err := memory.Get()
	if err != nil {
		return nil, err
	}

	return memory, nil
}

func connect(user *proto.Client) error {
	var streamError error
	fmt.Println("--CONNECTING STREAM--")
	fmt.Println(*user)
	stream, err := client.CreateStream(context.Background(), &proto.Connect{
		Client: user,
		Active: true,
	})
	if err != nil {
		return fmt.Errorf("Connect failed: %v", err)
	}

	wait.Add(1)

	go func(str proto.Broadcast_CreateStreamClient) {
		defer wait.Done()

		for {
			msg, err := str.Recv()

			if err != nil {
				streamError = fmt.Errorf("Error reading message: %v", err)
				break
			}

			if msg.SendInfo == true {
				fmt.Println("-- Server is asking for stats again -- ")
				sendReport(context.Background(), user)
			} else {
				fmt.Println("NO REQUEST")
			}
		}
	}(stream)

	return streamError
}

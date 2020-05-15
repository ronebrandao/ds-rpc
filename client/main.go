package main

import (
	"context"
	"ds-rpc/proto"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"google.golang.org/grpc"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	fmt.Println("!-- Information gattering starting in 10 seconds --!")

	client := proto.NewAddServiceClient(conn)

	ctx := context.Background()

	go func() {
		gocron.Every(10).Seconds().Do(func() {
			fmt.Println("!-- Begining information gattering --!")
			cpuUsage, err := getCPU()
			if err != nil {
				fmt.Println("*** It was not possible to obtain the CPU usage ***")
			}

			memoryStats, err := getMemoryStats()
			if err != nil {
				fmt.Println("*** It was not possible to memory the CPU usage ***")
			}

			req := &proto.Request{
				CPU:           cpuUsage,
				UsedRAM:       float32(memoryStats.Used),
				AvaliableRAM:  float32(memoryStats.Free),
				UsedDisk:      0.0,
				AvaliableDisk: 0.0,
			}

			if response, err := client.PerformanceReport(ctx, req); err == nil {
				printResponse(response)
			} else {
				printResponse(&proto.Response{Success: false})
			}
		})
	}()

	<-gocron.Start()

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

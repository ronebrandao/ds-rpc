package main

import (
	"ds-rpc/proto"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
)

type info struct {
	cpu int64
	usedRAM int64
	avaliableRAM int64
	usedDisk int64
	avaliableDisk int64
}

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)

	g := gin.Default()
	g.POST("/info", func(ctx *gin.Context) {
		form, errors := validateForm(ctx)

		if len(errors) > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		}

		req := &proto.Request{
			CPU: form.cpu,
			UsedRAM: form.usedRAM,
			AvaliableRAM: form.avaliableRAM,
			UsedDisk: form.usedDisk,
			AvaliableDisk: form.avaliableDisk,
		}

		if response, err := client.PerformanceReport(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{"result": fmt.Sprint(response.Success)})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})



	if err = g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}

func validateForm(ctx *gin.Context) (*info, []string) {
	errors := make([]string, 0)

	cpu, err := strconv.ParseInt(ctx.PostForm("cpu"), 10, 64)
	if err != nil {
		errors = append(errors, "Invalid Parameter 'cpu'")
	}

	usedRAM, err := strconv.ParseInt(ctx.PostForm("usedRAM"), 10, 64)
	if err != nil {
		errors = append(errors, "Invalid Parameter 'usedRAM'")
	}

	avaliableRAM, err := strconv.ParseInt(ctx.PostForm("avaliableRAM"), 10, 64)
	if err != nil {
		errors = append(errors, "Invalid Parameter 'avaliableRAM'")
	}

	usedDisk , err := strconv.ParseInt(ctx.PostForm("usedDisk"), 10, 64)
	if err != nil {
		errors = append(errors, "Invalid Parameter 'usedDisk'")
	}

	avaliableDisk  , err := strconv.ParseInt(ctx.PostForm("avaliableDisk"), 10, 64)
	if err != nil {
		errors = append(errors, "Invalid Parameter 'avaliableDisk'")
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return &info{
		cpu:           cpu,
		usedRAM:       usedRAM,
		avaliableRAM:  avaliableRAM,
		usedDisk:      usedDisk,
		avaliableDisk: avaliableDisk,
	}, nil
}

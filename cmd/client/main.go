package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/xuhe2/taskp/netapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	PORT = 1234
)

func NewBaseInfo() *netapi.BaseInfo {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &netapi.BaseInfo{Wd: workingDir}
}

func main() {
	conn, err := grpc.NewClient(fmt.Sprintf(":%d", PORT), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := netapi.NewTaskServiceClient(conn)

	taskName := flag.String("name", "<empty name>", "task name")
	command := flag.String("cmd", "echo 'empty'", "command to execute")
	flag.Parse()

	res, err := c.CommitTask(context.Background(), &netapi.CommitTaskReq{
		Task: &netapi.Task{Info: NewBaseInfo(), Name: *taskName, Command: *command},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.GetMessage())
}

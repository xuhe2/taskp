package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
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

	taskID := flag.Uint64("taskid", 0, "task id")
	taskName := flag.String("name", "", "task name")
	command := flag.String("cmd", "echo 'empty'", "command to execute")
	flag.Parse()

	// get first sub command
	mainCmd := flag.Arg(0)
	switch mainCmd {
	case "commit":
		handleCommitTaskCmd(c, *taskName, *command)
	case "list":
		handleListTaskCmd(c, *taskID, *taskName)
	default:
		panic("unknown command")
	}
}

func handleCommitTaskCmd(client netapi.TaskServiceClient, name, cmd string) {
	res, err := client.CommitTask(context.Background(), &netapi.CommitTaskReq{
		Task: &netapi.Task{Info: NewBaseInfo(), Name: name, Command: cmd},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.GetMessage())
}

func handleListTaskCmd(client netapi.TaskServiceClient, taskID uint64, taskName string) {
	res, err := client.GetTask(context.Background(), &netapi.GetTaskReq{TaskId: taskID, Name: taskName})
	if err != nil {
		panic(err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Task Name", "Status"})
	for _, task := range res.Tasks {
		t.AppendRow(table.Row{"", task.Name, ""})
		t.AppendSeparator()
	}
	t.Render()
}

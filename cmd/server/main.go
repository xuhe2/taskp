package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/xuhe2/taskp/core/db"
	"github.com/xuhe2/taskp/core/gvm"
	"github.com/xuhe2/taskp/core/task"
	"github.com/xuhe2/taskp/netapi"
	"google.golang.org/grpc"
)

const (
	PORT = 1234
)

type TaskServer struct {
	netapi.UnimplementedTaskServiceServer
}

func (s *TaskServer) CommitTask(ctx context.Context, in *netapi.Task) (*netapi.TaskResponse, error) {
	// async run task
	newTask := task.NewTask(in.Name, in.Info.Wd, in.Command)

	taskRecord := newTask.ToTaskRecord()

	databse, err := gvm.GetGlobalVar[*db.Database]("db")
	if err != nil {
		return nil, err
	}
	databse.Create(taskRecord)

	taskChan, err := gvm.GetGlobalVar[chan *task.Task]("taskChan")
	if err != nil {
		return nil, err
	}
	taskChan <- newTask

	return &netapi.TaskResponse{Message: "ok"}, nil
}

func main() {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}

	log.Printf("server listening at :%d", PORT)

	grpcServer := grpc.NewServer()
	netapi.RegisterTaskServiceServer(grpcServer, &TaskServer{})
	grpcServer.Serve(list)
}

func init() {
	initDatabase()
	initWorkers(1)
}

func initDatabase() {
	database := db.NewDatabase()
	database.InitFromDSN("./test.db")
	gvm.SetGlobalVar("db", database)
}

func initWorkers(num int) {
	taskChan := make(chan *task.Task, 1_000)
	gvm.SetGlobalVar("taskChan", taskChan)
	for i := 0; i < num; i++ {
		worker := task.NewWorker(i, taskChan)
		go worker.Run()
	}
}

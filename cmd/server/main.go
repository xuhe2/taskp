package main

import (
	"context"
	"flag"
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

func (s *TaskServer) CommitTask(ctx context.Context, in *netapi.CommitTaskReq) (*netapi.CommitTaskResp, error) {
	databse, err := gvm.GetGlobalVar[*db.Database]("db")
	if err != nil {
		return nil, err
	}

	// create task record when commit
	createTaskRecordWhenCommitFunc := func(t *task.Task) {
		databse.Create(t.ToTaskRecord())
	}

	// update task record
	// when task is running or finished
	updateTaskRecordFunc := func(t *task.Task) {
		databse.Save(t.ToTaskRecord())
	}

	// create a task
	newTask := task.NewTask(in.Task.Name, in.Task.Info.Wd, in.Task.Command).
		WithBeforeRunFunc(updateTaskRecordFunc).
		WithAfterRunFunc(updateTaskRecordFunc)

	// save task record to db
	createTaskRecordWhenCommitFunc(newTask)

	taskChan, err := gvm.GetGlobalVar[chan *task.Task]("taskChan")
	if err != nil {
		return nil, err
	}
	taskChan <- newTask

	return &netapi.CommitTaskResp{Message: fmt.Sprintf("task %s is commited", in.Task.Name)}, nil
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
	numOfWorkers := flag.Int("num_workers", 1, "number of workers")
	flag.Parse()

	initDatabase()
	initWorkers(*numOfWorkers)
}

func initDatabase() {
	database := db.NewDatabase()
	database.InitFromDSN("./test.db")
	gvm.SetGlobalVar("db", database)
}

func initWorkers(num int) {
	taskChan := make(chan *task.Task, 1_000)
	gvm.SetGlobalVar("taskChan", taskChan)
	for i := range num {
		worker := task.NewWorker(i, taskChan)
		go worker.Run()
	}
}

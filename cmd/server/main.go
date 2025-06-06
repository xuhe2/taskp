package main

import (
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

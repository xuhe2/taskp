package main

import (
	"context"
	"fmt"

	"github.com/xuhe2/taskp/core/db"
	"github.com/xuhe2/taskp/core/gvm"
	"github.com/xuhe2/taskp/core/task"
	"github.com/xuhe2/taskp/netapi"
	"gorm.io/gorm"
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

func (s *TaskServer) GetTask(ctx context.Context, in *netapi.GetTaskReq) (*netapi.GetTaskResp, error) {
	taskName := in.Name
	taskID := in.TaskId

	databse, err := gvm.GetGlobalVar[*db.Database]("db")
	if err != nil {
		return nil, err
	}

	var tx *gorm.DB
	// task id is the most important condition
	if taskID != 0 {
		tx = databse.Where("id = ?", taskID)
	}
	if taskName != "" && tx == nil {
		tx = databse.Where("name LIKE ?", taskName)
	}

	// default query all tasks
	// limit 10 tasks, order by `createdAt` field
	if tx == nil {
		tx = databse.Limit(5).Order("created_at desc")
	}

	var taskRecords []*db.TaskRecord
	tx.Find(&taskRecords) //find all task records

	var tasks []*netapi.Task
	for _, taskRecord := range taskRecords {
		tasks = append(tasks, &netapi.Task{
			Id:         uint64(taskRecord.ID),
			Name:       taskRecord.Name,
			Command:    taskRecord.Command,
			Status:     taskRecord.Status,
			CommitTime: taskRecord.CreatedAt.Format("2006-01-02 15:04:05"),
			StartTime:  taskRecord.StartTime.Format("2006-01-02 15:04:05"),
			StopTime:   taskRecord.StopTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &netapi.GetTaskResp{Tasks: tasks}, nil
}

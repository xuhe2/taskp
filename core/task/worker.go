package task

import (
	"fmt"
	"log/slog"

	"github.com/xuhe2/taskp/core/db"
	"github.com/xuhe2/taskp/core/gvm"
)

type Worker struct {
	ID                int
	TaskChannel       <-chan *Task
	QuitSignalChannel chan struct{}
}

func NewWorker(id int, taskChan <-chan *Task) *Worker {
	return &Worker{
		ID:                id,
		TaskChannel:       taskChan,
		QuitSignalChannel: make(chan struct{}),
	}
}

// this function is for goroutine, so it will block until receive a signal from QuitSignalChannel
func (w *Worker) Run() {
	for {
		select {
		case task := <-w.TaskChannel:
			// get database
			database, err := gvm.GetGlobalVar[*db.Database]("db")
			if err != nil {
				slog.Error("Get global var db failed", "error", err)
				continue
			}
			// exec task
			slog.Info(fmt.Sprintf("Worker `%d` start task: `%s`", w.ID, task.Name))
			task.Run()
			// mark task as finish
			slog.Info(fmt.Sprintf("Worker `%d` finished task: `%s`", w.ID, task.Name))
			if err := database.Save(task.ToTaskRecord()).Error; err != nil {
				slog.Error("Save task record failed", "error", err)
			}
		case <-w.QuitSignalChannel:
			// break loop when receive quit signal
			return
		}
	}
}

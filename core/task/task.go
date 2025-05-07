package task

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/xuhe2/taskp/core/db"
)

const (
	Status_Pending  = "pending"
	Status_Running  = "running"
	Status_Finished = "finished"
	Status_Failed   = "failed"
)

type Task struct {
	Name       string
	WorkingDir string
	Command    string

	Priority int // more high more prior; 0 is default
	Status   string

	StartTime time.Time
	StopTime  time.Time

	LogFile io.Writer

	Record *db.TaskRecord
}

func NewTask(name, wd, command string) *Task {
	return &Task{
		Name:       name,
		WorkingDir: wd,
		Command:    command,
		Priority:   0,
		Status:     Status_Pending,
		Record:     nil,
	}
}

func (t *Task) Run() {
	t.Status = Status_Running
	t.StartTime = time.Now()

	cmd := exec.Command("bash", "-c", t.Command)
	cmd.Dir = t.WorkingDir // set the working directory

	if t.LogFile == nil {
		// redirect stdout and stderr to os.Stdout and os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = t.LogFile
		cmd.Stderr = t.LogFile
	}

	slog.Info(fmt.Sprintf("Task `%s` started", t.Name))
	t.StartTime = time.Now() // record the start time
	err := cmd.Run()         // run the command
	defer func() {           // record the stop time
		t.StopTime = time.Now()
	}()

	if err != nil {
		slog.Error(fmt.Sprintf("Task `%s` failed: %v", t.Name, err))
		t.Status = Status_Failed
	} else {
		slog.Info(fmt.Sprintf("Task `%s` finished", t.Name))
		t.Status = Status_Finished
	}
}

func (t *Task) ToTaskRecord() *db.TaskRecord {
	var taskRecord *db.TaskRecord
	if t.Record != nil {
		taskRecord = t.Record
	} else {
		taskRecord = db.NewTaskRecord(t.Name, t.WorkingDir, t.Command)
	}

	taskRecord.Priority = t.Priority

	taskRecord.Status = t.Status

	taskRecord.StartTime = t.StartTime
	taskRecord.StopTime = t.StopTime

	// store task record in field
	t.Record = taskRecord

	return t.Record
}

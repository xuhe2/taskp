package db

import (
	"time"

	"gorm.io/gorm"
)

const (
	TaskRecordTable = "task_records"
)

type TaskRecord struct {
	gorm.Model
	Name       string
	WorkingDir string
	Command    string

	Priority int // more high, more priority; 0 is the default

	Status string

	StartTime time.Time
	StopTime  time.Time

	LogFile string
}

func NewTaskRecord(name, workingDir, command string) *TaskRecord {
	return &TaskRecord{
		Name:       name,
		WorkingDir: workingDir,
		Command:    command,
		Priority:   0,
	}
}

func (TaskRecord) TableName() string {
	return TaskRecordTable
}

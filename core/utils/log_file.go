package utils

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// generate log file from task id and return file pointer
// log file is store in `/var/log/taskp/<taskID>.log`
func GenLogFile(taskID uint) *os.File {
	// log file is stored in `~/.local/share/taskp`
	LOG_DIR_PATH := filepath.Join(os.Getenv("HOME"), ".local/share/taskp", "task_logs")

	// if log directory not exist, create it
	if _, err := os.Stat(LOG_DIR_PATH); os.IsNotExist(err) {
		if err := os.MkdirAll(LOG_DIR_PATH, 0755); err != nil {
			slog.Error(fmt.Sprintf("Error creating log directory: %s", err))
			return nil
		}
	}
	// if log file not exist, create it
	logFile, err := os.OpenFile(filepath.Join(LOG_DIR_PATH, fmt.Sprintf("%d.log", taskID)),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) // create log file
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating log file: %s", err))
		return nil
	}
	return logFile
}

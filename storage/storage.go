package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"task-cli/tasks"
)

type FileStorage struct {
	filename string
}

func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{filename}
}

func (fs *FileStorage) ReadTasks() ([]tasks.Task, error) {
	file, err := os.OpenFile(fs.filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var tasks []tasks.Task
	if len(data) > 0 {
		if err := json.Unmarshal(data, &tasks); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
		}
	}

	return tasks, nil
}

func (fs *FileStorage) WriteTasks(tasks []tasks.Task) error {
	file, err := os.OpenFile(fs.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write tasks to file: %w", err)
	}

	return nil
}

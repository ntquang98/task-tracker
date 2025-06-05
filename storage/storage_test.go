package storage

import (
	"os"
	"testing"
	"time"

	"github.com/ntquang98/task-tracker/tasks"
)

func TestFileStorage(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "tasks_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	fs := NewFileStorage(tmpFile.Name())

	now := time.Now()
	tasks := []tasks.Task{
		tasks.Task{ID: 1, Description: "Test task 1", Status: "todo", CreatedAt: &now, UpdatedAt: &now},
	}

	err = fs.WriteTasks(tasks)
	if err != nil {
		t.Fatalf("Error writing tasks to file: %v", err)
	}

	readTasks, err := fs.ReadTasks()
	if err != nil {
		t.Fatalf("Error reading tasks from file: %v", err)
	}

	if len(readTasks) != 1 || readTasks[0].Description != "Test task 1" {
		t.Errorf("Expected read tasks to match written tasks. Got: %+v", readTasks)
	}
}

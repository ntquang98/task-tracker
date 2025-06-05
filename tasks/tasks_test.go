package tasks

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// Mock storage for testing purposes.

type mockStorage struct {
	tasks []Task
	err   error
}

func (m *mockStorage) ReadTasks() ([]Task, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.tasks, nil
}

func (m *mockStorage) WriteTasks(tasks []Task) error {
	if m.err != nil {
		return m.err
	}

	m.tasks = tasks

	return nil
}

// End mock

func TestAddTask(t *testing.T) {
	store := &mockStorage{}
	tm := NewTaskManager(store)

	err := tm.AddTask("Buy groceries")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	tasks, _ := store.ReadTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Description != "Buy groceries" || tasks[0].Status != "todo" || tasks[0].ID != 1 {
		t.Errorf("Unexpected task data: %+v", tasks[0])
	}
}

func TestUpdateTask(t *testing.T) {
	now := time.Now()
	store := &mockStorage{
		tasks: []Task{
			{ID: 1, Description: "Buy groceries", Status: "todo", CreatedAt: &now, UpdatedAt: &now},
		},
	}
	tm := NewTaskManager(store)

	expectedDesc := "Buy groceries and cook dinner"
	if err := tm.UpdateTask(1, "Buy groceries and cook dinner"); err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	tasks, _ := store.ReadTasks()
	if tasks[0].Description != "Buy groceries and cook dinner" {
		t.Errorf("Unexpected task description '%s' got %s", expectedDesc, tasks[0].Description)
	}

	err := tm.UpdateTask(2, "Non-existent check")
	if err == nil || err.Error() != "task with ID 2 not found" {
		t.Errorf("Expected error for non-existing task, got: %v", err)
	}
}

func TestDeleteTask(t *testing.T) {
	now := time.Now()
	store := &mockStorage{
		tasks: []Task{
			{ID: 1, Description: "Buy groceries", Status: "todo", CreatedAt: &now, UpdatedAt: &now},
		},
	}
	tm := NewTaskManager(store)

	err := tm.DeleteTask(1)
	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	tasks, _ := store.ReadTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected no tasks left after deletion, got %d", len(tasks))
	}

	err = tm.DeleteTask(1)
	if err == nil || err.Error() != "task with ID 1 not found" {
		t.Errorf("Expected error for non-existing task, got: %v", err)
	}
}

func TestMarkTaskInProgress(t *testing.T) {
	now := time.Now()
	store := &mockStorage{
		tasks: []Task{
			{ID: 1, Description: "Buy groceries", Status: "todo", CreatedAt: &now, UpdatedAt: &now},
		},
	}
	tm := NewTaskManager(store)

	err := tm.MarkTaskInProgress(1)
	if err != nil {
		t.Fatalf("MarkTaskInProgress failed: %v", err)
	}

	tasks, _ := store.ReadTasks()
	if tasks[0].Status != "in-progress" {
		t.Errorf("Unexpected task status 'in-progress', got %s", tasks[0].Status)
	}
}

func TestMarkDone(t *testing.T) {
	now := time.Now()
	store := &mockStorage{
		tasks: []Task{
			{ID: 1, Description: "Buy groceries", Status: "todo", CreatedAt: &now, UpdatedAt: &now},
		},
	}
	tm := NewTaskManager(store)

	err := tm.MarkTaskDone(1)
	if err != nil {
		t.Fatalf("MarkTaskDone failed: %v", err)
	}

	tasks, _ := store.ReadTasks()
	if tasks[0].Status != "done" {
		t.Errorf("Unexpected task status 'done', got %s", tasks[0].Status)
	}
}

func TestMarkTodo(t *testing.T) {
	now := time.Now()
	store := &mockStorage{
		tasks: []Task{
			{ID: 1, Description: "Buy groceries", Status: "in-progress", CreatedAt: &now, UpdatedAt: &now},
		},
	}
	tm := NewTaskManager(store)

	err := tm.MarkTaskTodo(1)
	if err != nil {
		t.Fatalf("MarkTaskTodo failed: %v", err)
	}

	tasks, _ := store.ReadTasks()
	if tasks[0].Status != "todo" {
		t.Errorf("Unexpected task status 'todo', got %s", tasks[0].Status)
	}
}

func TestListTasks(t *testing.T) {
	now := time.Now()
	store := &mockStorage{
		tasks: []Task{
			{ID: 1, Description: "Buy groceries", Status: "in-progress", CreatedAt: &now, UpdatedAt: &now},
			{ID: 2, Description: "Cook dinner", Status: "todo", CreatedAt: &now, UpdatedAt: &now},
		},
	}
	tm := NewTaskManager(store)

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := tm.ListTasks("todo")
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}

	var buf strings.Builder
	_, _ = io.Copy(&buf, r)
	output := buf.String()
	if !strings.Contains(output, "Cook dinner") || strings.Contains(output, "Buy groceries") {
		t.Errorf("Unexpected output:\n%s", buf.String())
	}
}

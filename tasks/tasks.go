package tasks

import (
	"fmt"
	"time"
)

type Task struct {
	ID          int64      `json:"id"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

func (t Task) String() string {
	statusIcon := map[string]string{
		"done":        "✅",
		"in-progress": "🔄",
		"todo":        "📋",
	}[t.Status]
	return fmt.Sprintf("%s ID: %-4d %s", statusIcon, t.ID, t.Description)
}

type Storage interface {
	ReadTasks() ([]Task, error)
	WriteTasks([]Task) error
}

type TaskManager struct {
	store Storage
}

func NewTaskManager(store Storage) *TaskManager {
	return &TaskManager{store}
}

func (tm *TaskManager) AddTask(description string) error {
	tasks, err := tm.store.ReadTasks()
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	now := time.Now()
	newID := int64(1)
	if len(tasks) > 0 {
		newID = tasks[len(tasks)-1].ID + 1
	}

	newTask := Task{
		ID:          newID,
		Description: description,
		Status:      "todo",
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	tasks = append(tasks, newTask)
	if err := tm.store.WriteTasks(tasks); err != nil {
		return fmt.Errorf("failed to write tasks: %w", err)
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)

	return nil
}

func (tm *TaskManager) UpdateTask(id int64, description string) error {
	tasks, err := tm.store.ReadTasks()
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	for i, task := range tasks {
		if task.ID == id {
			now := time.Now()
			tasks[i].Description = description
			tasks[i].UpdatedAt = &now

			if err := tm.store.WriteTasks(tasks); err != nil {
				return fmt.Errorf("failed to write tasks: %w", err)
			}

			fmt.Printf("Task updated successfully (ID: %d)\n", id)
			return nil
		}
	}

	return fmt.Errorf("task with ID %d not found", id)
}

func (tm *TaskManager) DeleteTask(id int64) error {
	tasks, err := tm.store.ReadTasks()
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	newTasks := make([]Task, 0, len(tasks))
	found := false

	for _, task := range tasks {
		if task.ID == id {
			found = true
			continue
		}
		newTasks = append(newTasks, task)
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	if err := tm.store.WriteTasks(newTasks); err != nil {
		return fmt.Errorf("failed to write tasks: %w", err)
	}

	return nil
}

func (tm *TaskManager) ListTasks(filter string) error {
	tasks, err := tm.store.ReadTasks()
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Println("Looking good, no pending tasks 😄")
		return nil
	}

	for _, task := range tasks {
		if filter != "" && task.Status != filter {
			continue
		}
		fmt.Println(task)
	}

	return nil
}

func (tm *TaskManager) MarkTaskTodo(id int64) error {
	return tm.updateTaskStatus(id, "todo")
}

func (tm *TaskManager) MarkTaskInProgress(id int64) error {
	return tm.updateTaskStatus(id, "in-progress")
}

func (tm *TaskManager) MarkTaskDone(id int64) error {
	return tm.updateTaskStatus(id, "done")
}

func (tm *TaskManager) updateTaskStatus(id int64, status string) error {
	tasks, err := tm.store.ReadTasks()
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	for i, task := range tasks {
		if task.ID == id {
			now := time.Now()
			tasks[i].Status = status
			tasks[i].UpdatedAt = &now

			if err := tm.store.WriteTasks(tasks); err != nil {
				return fmt.Errorf("failed to write tasks: %w", err)
			}

			fmt.Printf("Task status updated successfully (ID: %d)\n", id)
			return nil
		}
	}

	return fmt.Errorf("task with ID %d not found", id)
}

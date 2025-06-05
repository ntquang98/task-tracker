package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int64      `json:"id"`
	Description string     `json:"decription"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updateAt"`
}

func (t Task) String() string {
	return fmt.Sprintf("ID(%d) %s - Status: %s", t.ID, t.Description, t.Status)
}

const helpText = `
NAME:
   task-cli - A command-line tool for managing tasks

USAGE:
   task-cli [command] [arguments]

COMMANDS:
   add <description>            Add a new task with the given description
                                Example: task-cli add "Buy groceries"
   
   update <id> <description>   Update the description of a task with the given ID
                                Example: task-cli update 1 "Buy groceries and cook dinner"
   
   delete <id>                 Delete a task with the given ID
                                Example: task-cli delete 1
   
   mark-in-progress <id>       Mark a task with the given ID as in-progress
                                Example: task-cli mark-in-progress 1
   
   mark-done <id>              Mark a task with the given ID as done
                                Example: task-cli mark-done 1
   
   list                        List all tasks
                                Example: task-cli list
   
   list <status>               List tasks filtered by status (todo, in-progress, done)
                                Example: task-cli list done

DESCRIPTION:
   task-cli is a simple tool to manage your tasks from the command line. You can add, update,
   delete, and mark tasks as in-progress or done. Use the list command to view all tasks or
   filter them by status.

EXAMPLES:
   Add a task:              task-cli add "Finish project report"
   Update a task:           task-cli update 1 "Finish project report and submit"
   Mark task as done:       task-cli mark-done 1
   List in-progress tasks:  task-cli list in-progress

Use "task-cli help [command]" for more information about a specific command.
`

var fileName string = "task.json"

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) <= 0 {
		log.Panicln("invalid features. using 'help' to get list of features.")
	}

	feat := argsWithoutProg[0]

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Error opening/creating file: ", err)
		return
	}

	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panicln("Error while reading file: ", err.Error())
	}

	var tasks []Task
	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, &tasks)
		if err != nil {
			log.Panicln("Error while unmarshal: ", err.Error())
		}
	}

	switch feat {
	case "list":
		listTask(argsWithoutProg[1:], tasks)
	case "add":
		addTask(argsWithoutProg[1:], file, tasks)
	case "update":
		updateTask(argsWithoutProg[1:], file, tasks)
	case "delete":
		deleteTask(argsWithoutProg[1:], file, tasks)
	case "mark-todo":
		updateTaskToDo(argsWithoutProg[1:], file, tasks)
	case "mark-in-progress":
		updateTaskInProgress(argsWithoutProg[1:], file, tasks)
	case "mark-done":
		updateTaskDone(argsWithoutProg[1:], file, tasks)
	case "help":
		help()
	default:
		fmt.Println("invalid features. using 'help' to get list of features.")
	}
}

func help() {
	fmt.Println(strings.TrimSpace(helpText))
	return
}

func listTask(args []string, tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("Tasks not found")
		return
	}

	var filterBy string
	if len(args) > 0 {
		filterBy = args[0]
	}

	for _, task := range tasks {
		if filterBy != "" && task.Status != filterBy {
			continue
		}

		fmt.Println(task)
	}
}

func addTask(args []string, file *os.File, tasks []Task) {
	if len(args) == 0 {
		log.Panicf("Invalid tasks")
		return
	}

	now := time.Now()
	newID := int64(len(tasks) + 1)

	newTask := Task{
		ID:          newID,
		Status:      "todo",
		Description: args[0],
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	tasks = append(tasks, newTask)

	if err := writeToFile(file, tasks); err != nil {
		log.Panicln("Error while write file", err.Error())
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)
}

func updateTask(args []string, file *os.File, tasks []Task) {
	if len(args) < 2 {
		log.Panic("Invalid argument")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.Panicf("Invalid ID %v", args[0])
	}

	newTasks := make([]Task, 0, len(tasks))

	now := time.Now()
	for i := range tasks {
		task := tasks[i]
		if task.ID == id {
			task.Description = args[1]
			task.UpdatedAt = &now
		}

		newTasks = append(newTasks, task)
	}

	if err := writeToFile(file, newTasks); err != nil {
		log.Panicln("Error while write file", err.Error())
	}

	fmt.Printf("Task updated successfully (ID: %d)\n", id)
}

func deleteTask(args []string, file *os.File, tasks []Task) {
	if len(args) < 1 {
		log.Panic("Invalid argument")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.Panicf("Invalid ID %v", args[0])
	}

	newTasks := make([]Task, 0, len(tasks)-1)
	for i := range tasks {
		task := tasks[i]
		if task.ID == id {
			continue
		}

		newTasks = append(newTasks, task)
	}

	if err := writeToFile(file, newTasks); err != nil {
		log.Panicln(err.Error())
	}

	fmt.Printf("Task deleted successfully (ID: %d)\n", id)
}

func writeToFile(file *os.File, tasks []Task) error {
	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	data, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("Error while marshal tasks", err.Error())
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("Error while write to file", err.Error())
	}

	return nil
}

func updateTaskToDo(args []string, file *os.File, tasks []Task) {
	if len(args) < 1 {
		log.Panic("Invalid argument")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.Panicf("Invalid ID %v", args[0])
	}

	newTasks := make([]Task, 0, len(tasks))

	now := time.Now()
	for _, task := range tasks {
		if task.ID == id {
			task.Status = "todo"
			task.UpdatedAt = &now
		}

		newTasks = append(newTasks, task)
	}

	if err := writeToFile(file, newTasks); err != nil {
		log.Panicln("Error while write file", err.Error())
	}

	fmt.Printf("Task updated successfully (ID: %d)\n", id)
}

func updateTaskInProgress(args []string, file *os.File, tasks []Task) {
	if len(args) < 1 {
		log.Panic("Invalid argument")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.Panicf("Invalid ID %v", args[0])
	}

	newTasks := make([]Task, 0, len(tasks))

	now := time.Now()
	for _, task := range tasks {
		if task.ID == id {
			task.Status = "in-progress"
			task.UpdatedAt = &now
		}

		newTasks = append(newTasks, task)
	}

	if err := writeToFile(file, newTasks); err != nil {
		log.Panicln("Error while write file", err.Error())
	}

	fmt.Printf("Task updated successfully (ID: %d)\n", id)
}

func updateTaskDone(args []string, file *os.File, tasks []Task) {
	if len(args) < 1 {
		log.Panic("Invalid argument")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.Panicf("Invalid ID %v", args[0])
	}

	newTasks := make([]Task, 0, len(tasks))

	now := time.Now()
	for _, task := range tasks {
		if task.ID == id {
			task.Status = "done"
			task.UpdatedAt = &now
		}

		newTasks = append(newTasks, task)
	}

	if err := writeToFile(file, newTasks); err != nil {
		log.Panicln("Error while write file", err.Error())
	}

	fmt.Printf("Task updated successfully (ID: %d)\n", id)
}

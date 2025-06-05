package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ntquang98/task-tracker/storage"
	"github.com/ntquang98/task-tracker/tasks"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(strings.TrimSpace(generalHelpText))
		os.Exit(1)
	}

	command := strings.ToLower(args[0])
	store := storage.NewFileStorage("tasks.json")
	taskManager := tasks.NewTaskManager(store)

	switch command {
	case "help":
		if len(args) == 1 {
			fmt.Println(strings.TrimSpace(generalHelpText))
			os.Exit(0)
		}

		subCommand := strings.ToLower(args[1])
		if helpText, exists := commandHelpTexts[subCommand]; exists {
			fmt.Println(strings.TrimSpace(helpText))
			os.Exit(0)
		}

		fmt.Printf("Error: no help available for unknown command '%s'\n", subCommand)
		fmt.Println("Run 'task-cli help' for a list of valid command")
		os.Exit(1)
	case "add":
		if len(args) < 2 {
			fmt.Println("Error: add command requires a description")
			fmt.Println(strings.TrimSpace(commandHelpTexts["add"]))
			os.Exit(1)
		}

		if err := taskManager.AddTask(args[1]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "update":
		if len(args) < 3 {
			fmt.Println("Error: update command requires an ID and description")
			fmt.Println(strings.TrimSpace(commandHelpTexts["update"]))
			os.Exit(1)
		}

		id, err := parseID(args[1])
		if err != nil {
			fmt.Printf("Error: invalid ID format: %v\n", err)
			os.Exit(1)
		}

		if err := taskManager.UpdateTask(id, args[2]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "delete":
		if len(args) < 2 {
			fmt.Println("Error: delete command requires an ID")
			fmt.Println(strings.TrimSpace(commandHelpTexts["delete"]))
			os.Exit(1)
		}

		id, err := parseID(args[1])
		if err != nil {
			fmt.Printf("Error: invalid ID format: %v\n", err)
			os.Exit(1)
		}

		if err := taskManager.DeleteTask(id); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "list":
		var filter string
		if len(args) > 1 {
			filter = args[1]
			if !isValidStatus(filter) {
				fmt.Println("Error: invalid status. Use 'todo', 'in-progress', or 'done'.")
				fmt.Println(strings.TrimSpace(commandHelpTexts["list"]))
				os.Exit(1)
			}
		}

		if err := taskManager.ListTasks(filter); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "mark-todo":
		if len(args) < 2 {
			fmt.Println("Error: mark-todo command requires an ID")
			fmt.Println(strings.TrimSpace(commandHelpTexts["mark-todo"]))
			os.Exit(1)
		}

		id, err := parseID(args[1])
		if err != nil {
			fmt.Printf("Error: invalid ID format: %v\n", err)
			os.Exit(1)
		}

		if err := taskManager.MarkTaskTodo(id); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "mark-in-progress":
		if len(args) < 2 {
			fmt.Println("Error: mark-in-progress command requires an ID")
			fmt.Println(strings.TrimSpace(commandHelpTexts["mark-in-progress"]))
			os.Exit(1)
		}

		id, err := parseID(args[1])
		if err != nil {
			fmt.Printf("Error: invalid ID format: %v\n", err)
			os.Exit(1)
		}

		if err := taskManager.MarkTaskInProgress(id); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "mark-done":
		if len(args) < 2 {
			fmt.Println("Error: mark-done command requires an ID")
			fmt.Println(strings.TrimSpace(commandHelpTexts["mark-done"]))
			os.Exit(1)
		}

		id, err := parseID(args[1])
		if err != nil {
			fmt.Printf("Error: invalid ID format: %v\n", err)
			os.Exit(1)
		}

		if err := taskManager.MarkTaskDone(id); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Error: unknown command '%s'\n", command)
		fmt.Println("Run 'task-cli help' for usage instructions.")
		os.Exit(1)
	}
}

func parseID(idStr string) (int64, error) {
	return strconv.ParseInt(idStr, 10, 64)
}

func isValidStatus(status string) bool {
	return status == "todo" || status == "in-progress" || status == "done"
}

const generalHelpText = `
NAME:
   task-cli - A command-line tool for managing tasks

USAGE:
   task-cli [command] [arguments]

COMMANDS:
   add <description>            Add a new task with the given description
   update <id> <description>   Update the description of a task with the given ID
   delete <id>                 Delete a task with the given ID
   mark-in-progress <id>       Mark a task with the given ID as in-progress
   mark-done <id>              Mark a task with the given ID as done
   list                        List all tasks
   list <status>               List tasks filtered by status (todo, in-progress, done)
   help [command]              Display help for a specific command

DESCRIPTION:
   task-cli is a simple tool to manage your tasks from the command line. You can add, update,
   delete, and mark tasks as in-progress or done. Use the list command to view all tasks or
   filter them by status.

EXAMPLES:
   Add a task:              task-cli add "Finish project report"
   Update a task:           task-cli update 1 "Finish project report and submit"
   Mark task as done:       task-cli mark-done 1
   List in-progress tasks:  task-cli list in-progress
`

var commandHelpTexts = map[string]string{
	"add": `
NAME:
   task-cli add - Add a new task

USAGE:
   task-cli add <description>

DESCRIPTION:
   Adds a new task with the provided description to the task list.

EXAMPLES:
   task-cli add "Buy groceries"
   task-cli add "Finish project report"

OUTPUT:
   Task added successfully (ID: <id>)
`,
	"update": `
NAME:
   task-cli update - Update an existing task

USAGE:
   task-cli update <id> <description>

DESCRIPTION:
   Updates the description of the task with the specified ID.

EXAMPLES:
   task-cli update 1 "Buy groceries and cook dinner"
`,
	"delete": `
NAME:
   task-cli delete - Delete a task

USAGE:
   task-cli delete <id>

DESCRIPTION:
   Deletes the task with the specified ID from the task list.

EXAMPLES:
   task-cli delete 1
`,
	"mark-todo": `
NAME:
   task-cli mark-todo - Mark a task as todo

USAGE:
   task-cli mark-todo <id>

DESCRIPTION:
   Marks the task with the specified ID as todo.

EXAMPLES:
   task-cli mark-todo 1
`,
	"mark-in-progress": `
NAME:
   task-cli mark-in-progress - Mark a task as in-progress

USAGE:
   task-cli mark-in-progress <id>

DESCRIPTION:
   Marks the task with the specified ID as in-progress.

EXAMPLES:
   task-cli mark-in-progress 1
`,
	"mark-done": `
NAME:
   task-cli mark-done - Mark a task as done

USAGE:
   task-cli mark-done <id>

DESCRIPTION:
   Marks the task with the specified ID as done.

EXAMPLES:
   task-cli mark-done 1
`,
	"list": `
NAME:
   task-cli list - List tasks

USAGE:
   task-cli list [todo|in-progress|done]

DESCRIPTION:
   Lists all tasks or tasks filtered by status (todo, in-progress, or done).
   If no status is provided, all tasks are listed.

EXAMPLES:
   task-cli list
   task-cli list done
`,
	"help": `
NAME:
   task-cli help - Display help information

USAGE:
   task-cli help [command]

DESCRIPTION:
   Displays general help for task-cli or detailed help for a specific command.

EXAMPLES:
   task-cli help
   task-cli help add
`,
}

# Task Tracker

[task-tracker](https://roadmap.sh/projects/task-tracker) Challenge from [roadmap.sh](https://roadmap.sh)

## How to run

### Build the project

```bash
# clone project
https://github.com/ntquang98/task-tracker
cd task-tracker

# build
go build -o task-cli ./cmd/task-cli
./task-cli help # list of available commands

# to add a task
./task-cli add "Buy groceries"

# To update a task
./task-cli update 1 "Buy groceries and cook dinner"

# To list task
./task-cli list # query all task
./task-cli list done # query task is done
./task-cli list todo # query task is todo
./task-cli list in-progress # query task is in-progress

# To mark a task as in progress/done/todo
./task-cli mark-in-progress 1
./task-cli mark-done 1
./task-cli mark-todo 1

# To Delete a task
./task-cli delete 1
```

### Install with go install

```bash
go install github.com/ntquang98/task-tracker/cmd/task-cli@latest
```
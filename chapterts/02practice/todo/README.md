# 📝 Go Todo CLI

A simple command-line TODO application written in Go, part of a learning/practice project.

## 🚀 Features

- Add new tasks with name and description
- List all tasks
- Get details of a task by ID
- Mark tasks as done
- Delete tasks
- Support for daemon operations (via directory input)

## 🛠 Usage

Run the app using:

```bash
go run .
```

You can use the following flags:

```
  -daemon string
        Operations source dir for daemon
  -del string
        Delete task by ID
  -done string
        Mark task as done by ID
  -get string
        Get task details by ID
  -list
        List all tasks
  -new string
        Create a new task by "<name>|<description>"
```

## 📌 Examples

### Create a new task

```bash
go run . -new "Buy groceries|Milk, Eggs, Bread"
```

### List all tasks

```bash
go run . -list
```

### Get a task by ID

```bash
go run . -get 1
```

### Mark a task as done

```bash
go run . -done 1
```

### Delete a task

```bash
go run . -del 1
```

## 📁 Project Structure

This is part of a Go learning series, located under:

```
chapterts/02practice/todo
```

## 🧠 Notes

- Tasks are stored locally.
- The `-daemon` flag is for optional directory-based task operations—implementation-dependent.


## TODO
[ ] impl daemon operations from files
[ ] tests
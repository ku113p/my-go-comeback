# 📝 Go Todo CLI

A simple command-line TODO application written in Go, part of a learning/practice project.

## 🚀 Features

- Add new tasks with name and description
- List all tasks
- Get details of a task by ID
- Mark tasks as done
- Delete tasks
- Support for **file-based automation via a daemon process**

## 🛠 Usage

Run the app using:

```bash
go run .
```

You can use the following flags:

```
  -daemon string
        Path to a directory to watch for file-based task operations
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

## 📦 Daemon Mode (`-daemon`)

Daemon mode allows automation of task operations using files. When the `-daemon` flag is used with a directory path, the CLI will:

1. **Watch the specified directory** every 10 seconds.
2. **Process files** in the directory with filenames starting with one of these prefixes:
   - `new_`: Create a new task.
   - `mark_`: Mark a task as done.
   - `delete_`: Delete a task.
3. **Expect each file** to contain JSON payloads describing the operation.
4. **Perform the operation**, save changes, and then delete the file.

### 📂 Example Usage

Run the daemon:

```bash
go run . -daemon ./ops
```

Then create files in the `./ops` directory:

#### ✅ `new_<any>.json`

```json
{
  "time": "2025-04-18T10:30:00Z",
  "name": "Read book",
  "desc": "Read 20 pages of 'Go in Action'"
}
```

#### ✅ `mark_<any>.json`

```json
{
  "id": "task-id-here"
}
```

#### ✅ `delete_<any>.json`

```json
{
  "id": "task-id-here"
}
```

The daemon will detect the files, process them accordingly, and delete them afterward.

> This feature is great for scripting, automation, or integration with other tools.

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
- Daemon mode is optional but useful for automating task input using external processes or integrations.

## TODO
- [ ] tests  
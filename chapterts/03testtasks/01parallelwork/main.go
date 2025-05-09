package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var filePath string
var maxWorkers int

func init() {
	fileFlag := flag.String("file", "", "Путь к файлу")
	workersFlag := flag.Int("workers", 1, "Максимальное количество воркеров")

	flag.Parse()

	if *fileFlag == "" {
		fmt.Println("Укажите путь к файлу с помощью -file")
		os.Exit(1)
	}

	filePath = *fileFlag
	maxWorkers = *workersFlag
}

type Task struct {
	payload string
	wg      *sync.WaitGroup
}

func newTask(payload string, wg *sync.WaitGroup) *Task {
	return &Task{payload, wg}
}

func main() {
	inputLines := readLinesFromFile(filePath)

	taskQueue := make(chan *Task)
	defer close(taskQueue) // Чтобы горутины закончивших воркеров не ждали еще активных
	workerInitDone := startWorkerPool(maxWorkers, taskQueue)

	var wg sync.WaitGroup
	defer wg.Wait()

	for line := range inputLines {
		wg.Add(1)
		task := newTask(line, &wg)

		// Либо задача сразу уходит в канал для воркеров,
		// либо мы ждем инициализации воркеров по возможности и задача всё равно уходит в канал
		select {
		case taskQueue <- task:
		case <-workerInitDone:
			taskQueue <- task
		}
	}
}

// readLinesFromFile читает строки из файла и отправляет их в канал
func readLinesFromFile(filePath string) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Ошибка при открытии файла: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()

	return lines
}

// startWorkerPool запускает пул из maxWorkers горутин-воркеров
// и возвращает канал, для получения сигналов, когда создавать
func startWorkerPool(maxWorkers int, taskQueue <-chan *Task) <-chan struct{} {
	needCreate := make(chan struct{})

	go func() {
		for range maxWorkers {
			needCreate <- struct{}{} // заполняем канал; в следующий раз создадим воркера только когда он освободится
			go worker(taskQueue)
		}
		close(needCreate) // достигли максимума воркеров
	}()

	return needCreate
}

// worker выполняет задачи, поступающие из taskQueue
func worker(taskQueue <-chan *Task) {
	for task := range taskQueue {
		process(task)
	}
}

// process обрабатывает одну задачу: если число — ждем и выводим
func process(task *Task) {
	defer task.wg.Done()

	ms, err := strconv.Atoi(task.payload)
	if err != nil {
		fmt.Printf("Строка содержит не число: %v\n", err)
		return
	}

	time.Sleep(time.Millisecond * time.Duration(ms))
	fmt.Println(task.payload)
}

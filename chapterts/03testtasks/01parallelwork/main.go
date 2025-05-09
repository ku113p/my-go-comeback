package main

import (
	"bufio"
	"context"
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
	fileFlag := flag.String("file", "", "Path to the input file")
	workersFlag := flag.Int("workers", 1, "Maximum number of workers")

	flag.Parse()

	if *fileFlag == "" {
		fmt.Println("Specify the path to the file using -file")
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputLines := readLinesFromFile(filePath)

	done := make(chan any)
	go func() {
		startWorkerPool(ctx, maxWorkers, inputLines)
		done <- struct{}{}
	}()

	<-done
}

// readLinesFromFile reads lines from a file and sends them to a channel
func readLinesFromFile(filePath string) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
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

// startWorkerPool starts a pool of maxWorkers goroutines (workers)
func startWorkerPool(ctx context.Context, maxWorkers int, inputLines <-chan string) {
	var wg sync.WaitGroup
	defer wg.Wait()

	taskQueue := make(chan *Task)
	produceWorker := make(chan any)

	go func() {
		// Fill the channel to signal available worker slots
		for range maxWorkers {
			produceWorker <- struct{}{}
			go worker(ctx, taskQueue)
		}
		close(produceWorker) // reached the maximum number of workers
	}()

	for line := range inputLines {
		wg.Add(1)
		task := newTask(line, &wg)

		select {
		case taskQueue <- task:
		case <-produceWorker:
			taskQueue <- task
		case <-ctx.Done():
			return
		}
	}

	close(taskQueue)
}

// worker processes tasks coming from taskQueue
func worker(ctx context.Context, taskQueue <-chan *Task) {
	for {
		select {
		case task, ok := <-taskQueue:
			if !ok {
				return
			}
			process(ctx, task)
		case <-ctx.Done():
			return
		}
	}
}

// process handles a single task: if the payload is a number, wait and print it
func process(_ context.Context, task *Task) {
	defer task.wg.Done()

	ms, err := strconv.Atoi(task.payload)
	if err != nil {
		fmt.Printf("Line is not a number: %v\n", err)
		return
	}

	time.Sleep(time.Millisecond * time.Duration(ms))
	fmt.Println(task.payload)
}

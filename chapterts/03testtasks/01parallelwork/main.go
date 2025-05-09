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

type toDo struct {
	s  string
	wg *sync.WaitGroup
}

func newToDo(s string, wg *sync.WaitGroup) *toDo {
	return &toDo{s, wg}
}

func main() {
	distributorTasks := tasksProducer(filePath)

	tasks := make(chan *toDo)
	wProducer := workerProducer(maxWorkers, tasks)

	var wg sync.WaitGroup
	defer wg.Wait()
	for t := range distributorTasks {
		wg.Add(1)
		td := newToDo(t, &wg)

		select {
		case tasks <- td:
		case <-wProducer:
			tasks <- td
		}
	}
}

func tasksProducer(filePath string) chan string {
	preTasks := make(chan string)
	go func() {
		defer close(preTasks)

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Ошибка при открытии файла: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			preTasks <- scanner.Text()
		}
	}()

	return preTasks
}

func workerProducer(maxWorkers int, tasks chan *toDo) <-chan any {
	waiter := make(chan any)

	go func() {
		for range maxWorkers {
			waiter <- struct{}{}
			go startNewWorker(tasks)
		}

		close(waiter)
	}()

	return waiter
}

func startNewWorker(tasks chan *toDo) {
	for t := range tasks {
		doWork(t)
	}
}

func doWork(t *toDo) {
	if toSleep, err := strconv.Atoi(t.s); err != nil {
		fmt.Printf("Строка содержит не число: %v\n", err)
	} else {
		time.Sleep(time.Millisecond * time.Duration(toSleep))
		fmt.Println(t.s)
	}
	t.wg.Done()
}

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func worker(id int, tasks <-chan int) {
	defer wg.Done()
	for t := range tasks {
		fmt.Printf("worker %d: processing %d\n", id, t)
		time.Sleep(400 * time.Millisecond)
	}
	fmt.Printf("worker %d: exiting\n", id)
}

func main() {
	tasks := make(chan int, 5)
	wg.Add(1)
	go worker(1, tasks)

	workerCount := 1
	taskID := 0

	// producer
	go func() {
		for {
			taskID++
			select {
			case tasks <- taskID:
				// normal case, channel not at capacity
				fmt.Printf("producer: queued %d\n", taskID)
			default:
				// channel full, add another worker
				workerCount++
				fmt.Printf("producer: channel full, adding worker %d\n", workerCount)
				wg.Add(1)
				go worker(workerCount, tasks)
				tasks <- taskID
			}

			// vary production rate
			time.Sleep(time.Duration(rand.Intn(200)+50) * time.Millisecond)
		}
	}()

	time.Sleep(5 * time.Second)
	close(tasks)
	wg.Wait()
	fmt.Println("done")
}

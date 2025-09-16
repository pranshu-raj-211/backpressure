package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func worker(id int, tasks <-chan int, quit <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case t, ok := <-tasks:
			if !ok {
				fmt.Printf("worker %d: exiting (channel closed)\n", id)
				return
			}
			fmt.Printf("worker %d: processing %d\n", id, t)
			time.Sleep(400 * time.Millisecond)
		case <-quit:
			// consumer sends a quit signal (analogy - master process of Nginx)
			fmt.Printf("worker %d: exiting (quit signal)\n", id)
			return
		}
	}
}

func main() {
	tasks := make(chan int, 10)
	stopProducer := make(chan struct{})

	workerQuitChans := make(map[int]chan struct{})
	workerCount := 0
	var mu sync.Mutex

	startWorker := func() {
		mu.Lock()
		defer mu.Unlock()
		workerCount++
		id := workerCount
		quit := make(chan struct{})
		workerQuitChans[id] = quit
		wg.Add(1)
		go worker(id, tasks, quit)
		fmt.Printf("controller: started worker %d (total=%d)\n", id, workerCount)
	}

	stopWorker := func() {
		mu.Lock()
		defer mu.Unlock()
		if workerCount == 0 {
			return
		}
		// stop a worker - scale down
		quit := workerQuitChans[workerCount]
		close(quit)
		delete(workerQuitChans, workerCount)
		workerCount--
		fmt.Printf("controller: stopped a worker (total=%d)\n", workerCount)
	}

	// init onw worker
	startWorker()

	// producer
	go func() {
		taskID := 0
		for {
			select {
			case <-stopProducer:
				close(tasks)
				return
			default:
				taskID++
				tasks <- taskID
				fmt.Printf("producer: queued %d\n", taskID)
				time.Sleep(time.Duration(rand.Intn(200)+50) * time.Millisecond)
			}
		}
	}()

	// consumer (controls number of workers according to load)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			queueLen := len(tasks)
			queueCap := cap(tasks)

			if queueLen > queueCap/2 {
				startWorker() // scale up
			} else if queueLen == 0 && workerCount > 1 {
				stopWorker() // scale down
			}
		}
	}()

	time.Sleep(10 * time.Second)
	close(stopProducer)

	mu.Lock()
	for _, quit := range workerQuitChans {
		close(quit)
	}
	mu.Unlock()

	wg.Wait()
	fmt.Println("done")
}

package main

import (
	"fmt"
	"time"
)


func main() {
	ch := make(chan int, 3)

	go func() {
		for i := 1; i <= 10; i++ {
			select {
			case ch <- i:
				fmt.Printf("producer: sent %d\n", i)
			default:
				// channel full, drop incoming
				fmt.Printf("producer: dropped %d (channel full)\n", i)
			}
			time.Sleep(50 * time.Millisecond)
		}
		close(ch)
	}()

	for v := range ch {
		fmt.Printf("consumer: processing %d\n", v)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println("done")
}

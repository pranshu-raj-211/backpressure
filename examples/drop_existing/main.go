package main

import (
	"fmt"
	"time"
)

func main() {
	var msg int
	ch := make(chan int, 3)
	drain := false // set true to drain whole channel when full, false to remove oldest message instead

	// producer
	go func() {
		for i := 1; i <= 40; i++ {
			select {
			case ch <- i:
				fmt.Printf("producer: sent %d\n", i)
			default:
				// drop existing - drain channel
				fmt.Println("producer: channel full")
				if drain {
					for len(ch) > 0 {
						<-ch
					}
					fmt.Println("drained channel")
					ch <- i
				} else {
					msg = <-ch // drop oldest
					fmt.Printf("dropped oldest %d\n", msg)
					ch <- i
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
		close(ch)
	}()

	// consumer
	for v := range ch {
		fmt.Printf("consumer: got %d\n", v)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println("done")
}

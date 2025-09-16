package main

import (
	"fmt"
	"time"
)

// using buffered channels - which block if the channel is full (till some messages are read)
// producer is faster than consumer
func main() {
	buf := 3
	ch := make(chan int, buf)

	// producer is much faster than consumer, but will be blocked until consumer
	// has consumed some message from the channel
	go func() {
		for i := 1; i <= 10; i++ {
			fmt.Printf("producer created %d (channel=%d/%d)\n", i, len(ch), cap(ch))
			ch <- i // blocking when channel is filled
			fmt.Printf("producer sent %d (channel=%d/%d)\n", i, len(ch), cap(ch))
			time.Sleep(100 * time.Millisecond)
		}
		close(ch)
	}()

	// Consumer
	for v := range ch {
		fmt.Printf("consumer received %d (buffer=%d/%d)\n", v, len(ch), cap(ch))
		time.Sleep(400 * time.Millisecond)
	}
	fmt.Println("done")
}

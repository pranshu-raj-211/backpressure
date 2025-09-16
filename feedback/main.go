package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 5)
	signal := make(chan int) // feedback channel

	// producer
	go func() {
		delay := 50 * time.Millisecond
		for i := 1; i <= 40; i++ {
			select {
			case ch <- i:
				fmt.Printf("producer: sent %d (delay=%v)\n", i, delay)
			default:
				fmt.Println("producer: channel full, backing off")
				delay *= 2 // multiplicative decrease
				i--        // retry same message
			}

			// adjust rate of messages
			select {
			case free := <-signal:
				if free > 2 && delay > 20*time.Millisecond {
					delay -= 10 * time.Millisecond // additive increase
				}
			default:
				// no feedback
			}

			time.Sleep(delay)
		}
		close(ch)
	}()

	// consumer
	go func() {
		for v := range ch {
			fmt.Printf("consumer: got %d\n", v)
			time.Sleep(200 * time.Millisecond)
			free := cap(ch) - len(ch)

			select {
			case signal <- free: // send feedback
			default:
			}
		}
		close(signal)
	}()

	time.Sleep(15 * time.Second)
	fmt.Println("done")
}

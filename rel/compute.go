package rel

import (
	"log"
	"time"
)

func compute() int {
	timerOneMs := time.NewTimer(2 * time.Millisecond)
	res := make(chan int, 1)
	go func() {
		count := 0
		for {
			select {
			case <-timerOneMs.C:
				res <- count
				return
			default:
				count++
			}
		}
	}()

	count := <-res
	log.Printf("counted in 1ms: %d", count)
	return count
}

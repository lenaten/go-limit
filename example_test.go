package golimit_test

import (
	"fmt"
	"time"

	"github.com/lenaten/go-limit"
)

// This example does several tasks concurrently,
// using a GoLimit to limit the number of concurrent tasks and block until all the tasks are complete.
func Example() {
	max := 5
	l := golimit.New(max)

	for i := 1; i <= max*2; i++ {
		// Increment the GoLimit counter and wait for their turn.
		l.Add(1)
		go func(i int) {
			// Decrement the counter when the goroutine completes.
			defer l.Done()
			// Do some work.
			time.Sleep(time.Millisecond * time.Duration(i) * 200)
			fmt.Println(i)
		}(i)
	}
	// Wait for all functions to complete.
	l.Wait()
}

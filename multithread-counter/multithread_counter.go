package main

import (
	"fmt"
	"sync"
)

// Counter is a thread-safe counter.
type Counter struct {
	mu    sync.Mutex
	value int
}

// NewCounter creates a new Counter.
func NewCounter() *Counter {
	return &Counter{}
}

// Increment safely increments the counter.
func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// Value returns the current value of the counter.
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	// Example usage
	c := NewCounter()
	var wg sync.WaitGroup
	numRoutines := 1000;
	numIncrements := 100;
	for routine := 0; routine < numRoutines; routine++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < numIncrements; i++ {
				c.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Println("Final counter value:", c.Value())
}
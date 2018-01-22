package golimit

// A GoLimit waits for a collection of goroutines to finish.
// The main goroutine calls New to set the limit of concurrent goroutines
// and Add to set the number of goroutines to wait for.
// Then each of the goroutines runs and calls Done when finished.
// At the same time, Wait can be used to block until all goroutines have finished.
type GoLimit struct {
	c chan struct{}
}

// New create a new Golimit with max concurrency number.
func New(max int) *GoLimit {
	return &GoLimit{
		c: make(chan struct{}, max),
	}
}

// Add adds delta, which may be negative, to the GoLimit counter.
//
// Note that calls with a positive delta that occur when the counter is zero
// must happen before a Wait. Calls with a negative delta, or calls with a
// positive delta that start when the counter is greater than zero, may happen
// at any time.
// Typically this means the calls to Add should execute before the statement
// creating the goroutine or other event to be waited for.
// If a GoLimit is reused to wait for several independent sets of events,
// new Add calls must happen after all previous Wait calls have returned.
// See the GoLimit example.
func (l GoLimit) Add(delta int) {
	if delta > 0 {
		for i := 0; i < delta; i++ {
			l.c <- struct{}{}
		}
	}
	if delta < 0 {
		for i := 0; i < (delta * -1); i++ {
			<-l.c
		}
	}
}

// Done decrements the GoLimit counter by one.
func (l GoLimit) Done() {
	l.Add(-1)
}

// Wait blocks until the GoLimit counter is zero.
func (l GoLimit) Wait() {
	// wait until all go are finished
	l.Add(cap(l.c))
	// allow new go by empty channel
	l.Add(len(l.c) * -1)
}

// Running return the number of running goroutine
func (l GoLimit) Running() int {
	return len(l.c)
}

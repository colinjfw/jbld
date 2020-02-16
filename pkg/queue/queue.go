package queue

import (
	"errors"
	"sync"
)

// ProcessFunc is the main unit of processing.
type ProcessFunc func(f string) ([]string, error)

// Run will kick off a queue.
func Run(workers int, initial []string, fn ProcessFunc) (int, error) {
	return New(workers, fn).Run(initial)
}

const (
	okState = iota
	errorState
	doneState
)

type state struct {
	state uint
	next  []string
	err   error
}

// New initializes a new queue.
func New(workers int, process ProcessFunc) *Queue {
	return &Queue{
		process: process,
		w:       workers,
	}
}

// Queue represents a generic work queue.
type Queue struct {
	process ProcessFunc
	w       int
	out     chan string
	state   chan state
}

// Dealer reads files from the in channel and closes in once all work is
// complete. We know that work is complete if len(in) == 0 and len(out) == 0 and
// there is no in flight work.
func (q *Queue) dealer(initial []string) (int, error) {
	var (
		count   int
		closing bool
		err     error
		cache   = map[string]bool{}
		buf     = initial
	)

	for {
		select {
		case state := <-q.state:
			if state.state == doneState {
				close(q.state)
				return count, err
			}
			if closing {
				continue
			}
			if state.state == errorState {
				close(q.out)
				closing = true
				err = state.err
				continue
			}

			if len(state.next) > 0 {
				buf = append(buf, state.next...)
			}

			count++
		default:
			if closing {
				continue
			}
			if len(buf) == 0 {
				if count == len(cache) {
					close(q.out)
					closing = true
				}
				continue
			}
			var f string
			f, buf = buf[0], buf[1:]
			if cache[f] {
				continue
			}
			cache[f] = true
			q.out <- f
		}
	}
}

func (q *Queue) worker(w int, wg *sync.WaitGroup) error {
	defer wg.Done()

	for f := range q.out {
		out, err := q.process(f)
		if err != nil {
			q.state <- state{state: errorState, err: err}
			continue
		}
		q.state <- state{state: okState, next: out}
	}
	return nil
}

// Run will traverse the graph.
func (q *Queue) Run(initial []string) (int, error) {
	if q.w == 0 {
		return 0, errors.New("queue: must have non-zero workers")
	}
	q.out = make(chan string, q.w)
	q.state = make(chan state, q.w)
	wg := &sync.WaitGroup{}
	wg.Add(q.w)
	for i := 0; i < q.w; i++ {
		go q.worker(i, wg)
	}

	go func() {
		wg.Wait()
		q.state <- state{state: doneState}
	}()
	return q.dealer(initial)
}

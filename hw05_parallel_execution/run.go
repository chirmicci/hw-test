package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	mErr := int64(m)
	cErr := int64(0)

	queue := make(chan Task)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(queue chan Task, wg *sync.WaitGroup, cErr *int64) {
			defer wg.Done()

			for task := range queue {
				if task() != nil {
					atomic.AddInt64(cErr, 1)
				}
			}
		}(queue, &wg, &cErr)
	}

	for _, task := range tasks {
		if atomic.LoadInt64(&cErr) >= mErr && mErr > 0 {
			break
		}
		queue <- task
	}

	close(queue)
	wg.Wait()

	if cErr >= mErr && mErr > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

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
	if n <= 0 {
		return nil
	}

	ignoreErrors := false

	if m <= 0 {
		ignoreErrors = true
	}

	var wg sync.WaitGroup

	var chTasks = make(chan Task)

	var numErrors int32

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range chTasks {
				err := task()
				if !ignoreErrors && err != nil {
					atomic.AddInt32(&numErrors, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if !ignoreErrors && atomic.LoadInt32(&numErrors) >= int32(m) {
			break
		}
		chTasks <- task
	}

	close(chTasks)
	wg.Wait()

	if !ignoreErrors && numErrors >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

package util

import (
	"context"
	"errors"
	"sync"
)

var ErrAllTasksFailed = errors.New("all tasks failed")

// FirstSuccess runs multiple tasks concurrently and returns the result
// from the first task that succeeds. If all tasks fail, it returns
// ErrAllTasksFailed along with the last error encountered.
//
// Tasks are functions that return a result and an error. The function
// will return as soon as one task succeeds, and will wait for all
// tasks to complete if all fail.
//
// Example:
//
//	result, err := FirstSuccess(
//		func() (string, error) { return fetchFromSource1() },
//		func() (string, error) { return fetchFromSource2() },
//		func() (string, error) { return fetchFromSource3() },
//	)
func FirstSuccess[T any](tasks ...func() (T, error)) (T, error) {
	return FirstSuccessWithContext(context.Background(), tasks...)
}

// FirstSuccessWithContext is like FirstSuccess but accepts a context
// for cancellation. If the context is cancelled, the function will
// return the context error.
func FirstSuccessWithContext[T any](ctx context.Context, tasks ...func() (T, error)) (T, error) {
	var zero T
	if len(tasks) == 0 {
		return zero, errors.New("no tasks provided")
	}

	type result struct {
		value T
		err   error
	}

	resultChan := make(chan result, len(tasks))
	var wg sync.WaitGroup

	// Launch all tasks concurrently
	for _, task := range tasks {
		wg.Add(1)
		go func(t func() (T, error)) {
			defer wg.Done()
			val, err := t()
			select {
			case resultChan <- result{value: val, err: err}:
			case <-ctx.Done():
			}
		}(task)
	}

	// Wait for the first success or all failures
	failureCount := 0

	// Create a goroutine to close the channel when all tasks complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for {
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case res, ok := <-resultChan:
			if !ok {
				// Channel closed, all tasks completed without success
				return zero, ErrAllTasksFailed
			}

			if res.err == nil {
				// First success found, return immediately
				return res.value, nil
			}

			// Track failures
			failureCount++
			if failureCount == len(tasks) {
				// All tasks failed
				return zero, ErrAllTasksFailed
			}
		}
	}
}

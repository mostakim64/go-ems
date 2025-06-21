package worker

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Executor interface {
	// Execute runs the task with the given name and payload.
	Execute() error
	// OnError handles errors that occur during task execution.
	OnError(error)
	// MaxRetries returns the maximum number of retries for the task.
	MaxRetries() int
}

type Pool struct {
	numWorkers int
	tasks      chan Executor
	start      sync.Once
	stop       sync.Once
	quit       chan struct{}
}

func NewPool(numWorkers int, taskChannelSize int) *Pool {
	if numWorkers <= 0 {
		panic("num workers must be greater than zero")
	}
	if taskChannelSize < 0 {
		panic("channel size cannot be negative")
	}

	return &Pool{
		numWorkers: numWorkers,
		tasks:      make(chan Executor, taskChannelSize),
		start:      sync.Once{},
		stop:       sync.Once{},
		quit:       make(chan struct{}),
	}
}

func (p *Pool) Start() {
	p.start.Do(func() {
		backgroundCtx := context.Background()
		p.startWorkers(backgroundCtx)
	})
}

func (p *Pool) Stop() {
	p.stop.Do(func() {
		close(p.quit)  // First signal to stop accepting new tasks
		close(p.tasks) // Then close the task channel to stop running tasks and return from infinite loop
	})
}

func (p *Pool) StopWithContext(ctx context.Context) {
	p.stop.Do(func() {
		close(p.quit)  // First signal to stop accepting new tasks
		close(p.tasks) // Then close the task channel to stop running tasks and return from infinite loop

		// Wait for context timeout
		<-ctx.Done()
		fmt.Println("Worker pool stopped with context timeout")
	})
}

func (p *Pool) AddTask(t Executor) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
}

func (p *Pool) startWorkers(ctx context.Context) {
	for i := 0; i < p.numWorkers; i++ {
		go func(workerNum int) {
			fmt.Printf("worker number %d started\n", workerNum)
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("worker number %d stopping due to context cancellation\n", workerNum)
					return
				case <-p.quit:
					fmt.Printf("worker number %d stopping due to quit signal\n", workerNum)
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}

					// Retry logic
					var err error
					for retry := 0; retry <= task.MaxRetries(); retry++ {
						if err = task.Execute(); err == nil {
							break
						}
						if retry < task.MaxRetries() {
							time.Sleep(time.Duration(retry+1) * time.Second)
							fmt.Printf("worker %d retrying task (attempt %d/%d)\n", workerNum, retry+1, task.MaxRetries())
						}
					}

					if err != nil {
						task.OnError(err)
					}

					// if err := task.Execute(); err != nil {
					// 	task.OnError(err)
					// }

					fmt.Printf("worker number %d finished task\n", workerNum)
				}
			}
		}(i)
	}
}

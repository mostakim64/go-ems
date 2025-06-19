package worker

import "time"

// Scheduler manages a single recurring task.
type Scheduler struct {
	ticker *time.Ticker
	done   chan bool
}

// NewScheduler creates a new scheduler for a given interval.
func NewScheduler(interval time.Duration) *Scheduler {
	return &Scheduler{
		ticker: time.NewTicker(interval),
		done:   make(chan bool),
	}
}

// Start begins the execution of the scheduled task.
func (s *Scheduler) Start(task func()) {
	go func() {
		for {
			select {
			case <-s.done:
				return
			case <-s.ticker.C:
				task()
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	s.ticker.Stop()
	s.done <- true
}

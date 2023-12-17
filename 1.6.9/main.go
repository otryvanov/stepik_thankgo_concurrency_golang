package main

import (
	"fmt"
	"time"
)

// начало решения

type cancellableTask struct {
	cancel   chan struct{}
	interval time.Duration
	fn       func()
}

func schedule(dur time.Duration, fn func()) func() {
	task := &cancellableTask{
		cancel:   make(chan struct{}),
		interval: dur,
		fn:       fn,
	}

	go func() {
		ticker := time.NewTicker(task.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				task.fn()
			case <-task.cancel:
				return
			}
		}
	}()

	return func() {
		select {
		case <-task.cancel:
		//уже отменено
		default:
			close(task.cancel)
		}
	}
}

// конец решения

func main() {
	work := func() {
		at := time.Now()
		fmt.Printf("%s: work done\n", at.Format("15:04:05.000"))
	}

	cancel := schedule(50*time.Millisecond, work)
	defer cancel()

	// хватит на 5 тиков
	time.Sleep(260 * time.Millisecond)
}

package concurrency

import (
	"sync"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(handler func()) {
	w.Add(1)
	go func() {
		handler()
		w.Done()
	}()
}

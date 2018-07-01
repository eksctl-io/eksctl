package utils

import "time"

type InstantTicker struct {
	innerTicker *time.Ticker
	C           <-chan time.Time
	stopChan    chan<- bool
}

func (it *InstantTicker) Stop() {
	it.innerTicker.Stop()
	it.stopChan <- true
	close(it.stopChan)
}

func replicateEvents(out chan<- time.Time, in <-chan time.Time, stopChan <-chan bool) {
	out <- time.Now()

	for {
		select {
		case t := <-in:
			out <- t
		case <-stopChan:
			return
		}
	}
}

func NewInstantTicker(d time.Duration) *InstantTicker {
	t := time.NewTicker(d)
	instantChan := make(chan time.Time)
	stopChan := make(chan bool)
	go replicateEvents(instantChan, t.C, stopChan)

	return &InstantTicker{
		innerTicker: t,
		C:           instantChan,
		stopChan:    stopChan,
	}
}

package metrics

import "time"

// Small helper so your handlers stay clean.
type Timer struct {
	start time.Time
}

func StartTimer() Timer {
	return Timer{start: time.Now()}
}

func (t Timer) Seconds() float64 {
	return time.Since(t.start).Seconds()
}

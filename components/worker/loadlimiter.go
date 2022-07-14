package worker

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/load"
)

// LoadLimiter stalls the system to limit load.
type LoadLimiter interface {
	LoadLimit() error
}

// NewLoadLimiter initializes and returns a load limiter based on 1-min system load per CPU, doubling sleep time up until maximum until load is below maximum.
func NewLoadLimiter(name string, maxRatio float64, throttleMin time.Duration, throttleMax time.Duration) LoadLimiter {
	return &loadLimiter{
		name:            name,
		maxLoad:         maxRatio * float64(runtime.NumCPU()),
		throttleCurrent: throttleMin,
		throttleMin:     throttleMin,
		throttleMax:     throttleMax,
		throttleJitter:  throttleMin,
		throttle:        throttle,
		getLoad:         getLoad1,
	}
}

type loadLimiter struct {
	name            string
	maxLoad         float64
	throttleMin     time.Duration
	throttleCurrent time.Duration
	throttleMax     time.Duration
	throttleJitter  time.Duration
	throttle        func(time.Duration)
	getLoad         func() (float64, error)
}

func throttle(t time.Duration) {
	time.Sleep(t)
}

func getLoad1() (float64, error) {
	load, err := load.Avg()
	if err != nil {
		return 0.0, err
	}

	return load.Load1, nil
}

func (l *loadLimiter) getJitter() time.Duration {
	nanos := l.throttleJitter.Nanoseconds()

	var jitter int64
	if nanos > 0 {
		jitter = rand.Int63n(nanos)
	} else {
		jitter = 0
	}

	return time.Duration(jitter)
}

func (l *loadLimiter) throttleDouble() {
	// Double next throttle time
	if next := l.throttleCurrent * 2; next > l.throttleMax {
		l.throttleCurrent = l.throttleMax
	} else {
		l.throttleCurrent = next
	}
}

func (l *loadLimiter) throttleReset() {
	l.throttleCurrent = l.throttleMin
}

func (l *loadLimiter) LoadLimit() error {
	load, err := l.getLoad()
	if err != nil {
		return err
	}
	if load > l.maxLoad {
		throttleTime := l.throttleCurrent + l.getJitter()
		log.Printf("(%s) Load %.2f above threshold %.2f, throttling for %.2fs.", l.name, load, l.maxLoad, throttleTime.Seconds())
		l.throttle(throttleTime)
		l.throttleDouble()

		return nil
	}

	l.throttleReset()
	return nil
}

// Compile-time assurance that implementation satisfies interface.
var _ LoadLimiter = &loadLimiter{}

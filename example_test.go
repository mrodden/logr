package logr_test

import (
	"sync"

	log "github.com/mrodden/logr"
)

func Example() {

	log.Info("Hello Logging")

	log.Debug("Debug message output")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Warn("Warning from inside a goroutine")
		wg.Done()
	}()

	log.Infof("Info format: %#v", wg)

	wg.Wait()
}

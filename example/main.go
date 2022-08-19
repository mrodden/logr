package main

import (
	"sync"

	log "github.com/mrodden/logr"
)

func main() {

	l := log.Default()

	l.Info("Hello Info")
	l.Debug("Debug message")
	l.Infof("Hello Infof: %v", l)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		l.Warn("Hello Warn from inside a goroutine")
		wg.Done()
	}()

	wg.Wait()
}

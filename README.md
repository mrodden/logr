# logr

Reasonable logging facilities for Go

Simple interface with Time-Thread-Condition-Component based output, with automatic colored output if supported.

## Features

 - Time: ISO 8901 format timestamps in UTC
 - Thread: Log source Goroutine identifiers
 - Condition: Log level based logging [DEBUG, INFO, WARN, ERROR, CRITICAL]
 - Component: Filename and line number
 - Colors: Automatically detects support for ANSI terminal colorization


## Example

```go
import (
	"sync"

	log "github.com/mrodden/logr"
)

func main() {

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
```

![Colored output](logr.png)

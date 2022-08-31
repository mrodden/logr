package fmt

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/mrodden/logr/logger"
)

type jsonOut struct {
	Level       string                 `json:"level,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Target      string                 `json:"target,omitempty"`
	Filename    string                 `json:"filename,omitempty"`
	LineNumber  int                    `json:"line_number,omitempty"`
	GoroutineId string                 `json:"goroutine_id,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

var JSON = FormatEvent(jsonEvent)

func jsonEvent(f Format, w *SyncWriter, event logger.Record) error {
	// level
	// timestamp
	// target/module
	// filename
	// line_number
	// fields
	//   message
	//   ...

	out := jsonOut{
		Level:      event.Metadata().Level().String(),
		Timestamp:  time.Now().UTC().Format(ISO),
		Target:     event.Metadata().Target(),
		Filename:   event.Filename(),
		LineNumber: int(event.LineNumber()),
		Fields: map[string]interface{}{
			"message": fmt.Sprint(event.Args()...),
		},
	}

	if f.DisplayGoroutineId {
		// goroutine lookup adds 17-20 microseconds from testing
		out.GoroutineId = GoRoutineID()
	}

	b, err := json.Marshal(out)
	if err != nil {
		fmt.Printf("things broke: %v\n", err)
		return err
	}

	b = append(b, "\n"...)

	_, err = w.Write(b)
	return err
}

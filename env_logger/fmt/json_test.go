package fmt

import (
	"bytes"
	"testing"

	"github.com/mrodden/logr/logger"
)

var result error

func BenchmarkJson(b *testing.B) {
	buf := new(bytes.Buffer)
	w := NewSyncWriter(buf)

	f := Format{Format: JSON, DisplayGoroutineId: false}

	rec := logger.NewRecordBuilder().
		Level(logger.INFO).
		Target("testing").
		Args("test message").
		PackagePath("testing.package.path").
		File("json_test.go").
		Line(1).
		Build()

	// reset after setup
	b.ResetTimer()

	var r error
	for n := 0; n < b.N; n++ {
		r = jsonEvent(f, w, rec)
	}

	result = r
}

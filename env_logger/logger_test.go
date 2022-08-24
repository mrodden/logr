package env_logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrodden/logr/logger"
)

func TestParse(t *testing.T) {
	ds := "info,logger=debug,logr,amodule=info"

	exp := []Directive{
		Directive{"logr", logger.TRACE},
		Directive{"logger", logger.DEBUG},
		Directive{"amodule", logger.INFO},
		Directive{"", logger.INFO},
	}

	assert.Equal(t, exp, parse(ds))
}

func TestParseNoEnv(t *testing.T) {
	ds := ""

	var exp []Directive

	assert.Equal(t, exp, parse(ds))
}

func TestNoColor(t *testing.T) {
	buf := new(bytes.Buffer)
	log := envLogger{
		directives: []Directive{
			Directive{"", logger.TRACE},
		},
		out: buf,
	}

	r := logger.NewRecordBuilder().
		Level(logger.INFO).
		Args("Stuff").
		Build()

	log.Log(r)

	assert.Contains(t, buf.String(), "Stuff")
	assert.NotContains(t, buf.String(), "\x1B[")
}

func TestColorPaint(t *testing.T) {
	assert.Equal(t, Green.Paint("blah"), "\x1B[32mblah\x1B[0m")
	assert.Equal(t, Green.Dimmed().Paint("blah"), "\x1B[2;32mblah\x1B[0m")
}

package logr

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoColor(t *testing.T) {
	buf := new(bytes.Buffer)
	log := New(buf)
	log.Info("Stuff")

	assert.NotContains(t, buf.String(), "\x1B[")
}

func TestCallerPackageName(t *testing.T) {
	assert.Equal(t, CallerPackageName(1), "logr")
	assert.Equal(t, CallerPackageName(2), "testing")
}

func TestColorPaint(t *testing.T) {
	assert.Equal(t, Green.Paint("blah"), "\x1B[32mblah\x1B[0m")
	assert.Equal(t, Green.Dimmed().Paint("blah"), "\x1B[2;32mblah\x1B[0m")
}

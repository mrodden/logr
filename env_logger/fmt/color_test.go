package fmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorPaint(t *testing.T) {
	assert.Equal(t, Green.Paint("blah"), "\x1B[32mblah\x1B[0m")
	assert.Equal(t, Green.Dimmed().Paint("blah"), "\x1B[2;32mblah\x1B[0m")
}

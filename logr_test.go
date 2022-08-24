package logr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallerPackageName(t *testing.T) {
	assert.Equal(t, "github.com/mrodden/logr", CallerPackageName(1))
	assert.Equal(t, "testing", CallerPackageName(2))
}

func TestCleanupPackageName(t *testing.T) {
	assert.Equal(t, "github.com/mrodden/logr/logger", cleanupPackageName("github.com/mrodden/logr/logger.(*Logger).Log"))
	assert.Equal(t, "github.com/stretchr/testify/assert", cleanupPackageName("github.com/stretchr/testify/assert.Equal"))
	assert.Equal(t, "main", cleanupPackageName("main.main"))
}

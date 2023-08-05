package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	t.Run("assert new app returns no error", func(t *testing.T) {
		_, err := NewApp("../flyspray.db")
		assert.NoError(t, err)
	})
}

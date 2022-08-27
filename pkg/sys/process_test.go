package sys

import (
	"testing"
)

func TestNewProcess(t *testing.T) {
	t.Run("echo", func(t *testing.T) {
		err := NewProcess("echo", []string{"hello", "world"})
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("unknownCommand", func(t *testing.T) {
		err := NewProcess("unknownCommand", []string{"hello", "world"})
		if err == nil {
			t.Error(err)
		}
	})
}

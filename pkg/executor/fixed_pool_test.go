package executor

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedPool(t *testing.T) {
	pool := NewFixedPool(100)
	assert.Equal(t, uint32(100), pool.threadNum)
	assert.GreaterOrEqual(t, runtime.NumGoroutine(), 100)

	for i := 0; i < 101; i++ {
		pool.SubmitFunc(func() {
			time.Sleep(time.Second)
		})
	}

}

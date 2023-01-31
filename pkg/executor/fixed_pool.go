package executor

import (
	"context"
	"sync/atomic"
	"time"
)

type fixedThreadPool struct {
	threadNum uint32
	taskChan  chan task
	running   int32
}

var _ Executor = (*fixedThreadPool)(nil)

// NewFixedPool 创建一个固定大小的"thread pool"
func NewFixedPool(n uint32) *fixedThreadPool {
	ftp := &fixedThreadPool{
		threadNum: n,
		taskChan:  make(chan task),
		running:   0,
	}
	for i := n; i > 0; i-- {
		ftp.startThread()
	}

	return ftp
}

// Submit 提交一个task接口类型的任务
func (ftp *fixedThreadPool) Submit(t task) {
	if t == nil {
		return
	}
	ftp.taskChan <- t
}

// SubmitFunc 提交一个func类型的任务
func (ftp *fixedThreadPool) SubmitFunc(f func()) {
	if f == nil {
		return
	}
	ftp.taskChan <- &funcWrapper{f: f}
}

// Wait 等待之前submit的任务都结束
func (ftp *fixedThreadPool) Wait(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			if atomic.LoadInt32(&ftp.running) == 0 {
				return nil
			}
			return context.DeadlineExceeded
		case <-ticker.C:
			if atomic.LoadInt32(&ftp.running) == 0 {
				return nil
			}
		}
	}
}

func (ftp *fixedThreadPool) startThread() {
	go func() {
		for t := range ftp.taskChan {
			if t == nil {
				// 收到结束信号
				return
			}
			atomic.AddInt32(&ftp.running, 1)
			t.Run()
			atomic.AddInt32(&ftp.running, -1)
		}
	}()
}

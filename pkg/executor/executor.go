package executor

import "context"

type Executor interface {
	Submit(t task)
	SubmitFunc(f func())
	Wait(context.Context) error
}

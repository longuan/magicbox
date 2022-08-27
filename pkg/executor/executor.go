package executor

import "context"

type executor interface {
	Submit()
	SubmitFunc(f func())
	Wait(context.Context) error
}

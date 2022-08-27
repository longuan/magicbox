package executor

type task interface {
	Run()
}

type funcWrapper struct {
	f func()
}

func (fw *funcWrapper) Run() {
	fw.f()
}

package worker

type Task struct {
	execute      func() error
	errorHandler func(error)
	maxRetries   int
}

func NewTask(execute func() error, errorHandler func(error), maxRetries int) *Task {
	return &Task{
		execute:      execute,
		errorHandler: errorHandler,
		maxRetries:   maxRetries,
	}
}

func (t *Task) Execute() error {
	return t.execute()
}

func (t *Task) OnError(err error) {
	t.errorHandler(err)
}

func (t *Task) MaxRetries() int {
	return t.maxRetries
}

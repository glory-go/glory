package cronjob

import (
	"fmt"
	"sync"

	"github.com/alexflint/go-arg"
)

type Job struct {
	tasks sync.Map
}

type Runner func() error

type jobArg struct {
	Name string `arg:"required"`
}

func NewJobContainer() *Job {
	return &Job{}
}

func (j *Job) Register(name string, runner Runner) {
	j.tasks.Store(name, runner)
}

func (j *Job) Run(name string) error {
	v, ok := j.tasks.Load(name)
	if !ok {
		return fmt.Errorf("job with name %v not found", name)
	}
	task, ok := v.(Runner)
	if !ok {
		return fmt.Errorf("job with name %v, type is %T, is not valid", name, v)
	}
	return task()
}

func (j *Job) RunWithArgs() error {
	args := &jobArg{}
	if err := arg.Parse(&args); err != nil {
		return err
	}

	return j.Run(args.Name)
}

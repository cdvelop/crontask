//go:build wasm

package crontask

import (
	"errors"
	"syscall/js"
)

type crontabAdapter struct{}

func AddNewTasks(tasksPath ...string) (*cronTask, error) {
	return newCronTask(&crontabAdapter{}, tasksPath...)
}

func (a *crontabAdapter) AddJob(schedule string, fn any, args ...any) error {
	jsFn := js.ValueOf(fn)
	jsArgs := make([]any, len(args))
	for i, arg := range args {
		jsArgs[i] = arg
	}

	js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) any {
		jsFn.Invoke(jsArgs...)
		return nil
	}), 0)

	return nil
}

func (a crontabAdapter) GetTasksFromPath(tasksPath ...string) ([]Tasks, error) {
	return nil, errors.New("frontend crontabAdapter not implemented GetTasksFromPath")
}

func (a *crontabAdapter) ExecuteCmd(cmd Task) error {
	js.Global().Call(cmd.Command, cmd.Args)
	return nil
}

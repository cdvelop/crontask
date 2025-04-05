//go:build wasm

package crontask

import (
	"errors"
	"strconv"
	"strings"
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
	filePath := "crontasks.yml" // Default path
	if len(tasksPath) > 0 && tasksPath[0] != "" {
		filePath = tasksPath[0]
	}

	// If path doesn't start with http, assume it's relative to current path
	if !strings.HasPrefix(filePath, "http://") && !strings.HasPrefix(filePath, "https://") {
		filePath = "https://localhost/" + filePath
	}

	// Use XMLHttpRequest for synchronous requests (since we need to return the result)
	xhr := js.Global().Get("XMLHttpRequest").New()
	xhr.Call("open", "GET", filePath, false)
	xhr.Call("send")

	status := xhr.Get("status").Int()
	if status != 200 {
		return nil, errors.New("failed to fetch YAML config: HTTP " + strconv.Itoa(status))
	}

	yamlText := xhr.Get("responseText").String()

	// Convert string to []byte for parsing
	parser := YAMLParser{}
	tasks, err := parser.ParseYAML([]byte(yamlText))
	if err != nil {
		return nil, err
	}

	// Wrap in Tasks slice to maintain compatibility
	return []Tasks{tasks}, nil
}

func (a *crontabAdapter) ExecuteCmd(cmd Task) error {
	js.Global().Call(cmd.Command, cmd.Args)
	return nil
}

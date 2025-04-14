//go:build wasm

package crontask

import (
	"errors"
	"strconv"
	"strings"
	"syscall/js"
)

// Inicializador específico para WASM
func newCronAdapter() cronAdapter {
	return &wasmAdapter{}
}

// Adaptador para entorno WASM
type wasmAdapter struct{}

func (a *wasmAdapter) AddJob(schedule string, fn any, args ...any) error {
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

func (a *wasmAdapter) GetTasksFromPath(tasksPath string) ([]Tasks, error) {
	filePath := getDefaultFilePathTasks(tasksPath)

	// If path doesn't start with http or https, assume it's relative to current path
	// or if it begins with "/" assume it's relative to domain root
	if !strings.HasPrefix(filePath, "http://") && !strings.HasPrefix(filePath, "https://") {
		// Get current location from window.location
		location := js.Global().Get("window").Get("location")
		origin := location.Get("origin").String()

		if strings.HasPrefix(filePath, "/") {
			// Absolute path from domain root
			filePath = origin + filePath
		} else {
			// Relative path from current directory
			pathname := location.Get("pathname").String()
			// Get the directory part of the current path
			lastSlash := strings.LastIndex(pathname, "/")
			if lastSlash > 0 {
				pathname = pathname[:lastSlash+1]
			}
			filePath = origin + pathname + filePath
		}
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
	parser := ymlParser{}
	tasks, err := parser.ParseYAML([]byte(yamlText))
	if err != nil {
		return nil, err
	}

	// Wrap in Tasks slice to maintain compatibility
	return []Tasks{tasks}, nil
}

func (a *wasmAdapter) ExecuteCmd(cmd Task) error {
	js.Global().Call(cmd.Command, cmd.Args)
	return nil
}

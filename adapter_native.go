//go:build !wasm

package crontask

import (
	"os"
	"os/exec"
)

// Inicializador específico para entornos no-WASM
func newCronAdapter() cronAdapter {
	return &nativeAdapter{ctab: newCrontab()}
}

// Adaptador para entornos nativos (no-WASM)
type nativeAdapter struct {
	ctab *crontab
}

func (a *nativeAdapter) AddJob(schedule string, fn any, args ...any) error {
	jobFunc, ok := fn.(func())
	if !ok {
		return newErr("invalid function type")
	}
	return a.ctab.AddJob(schedule, jobFunc, args...)
}

func (a *nativeAdapter) RunAll() {
	a.ctab.RunAll()
}

func (a *nativeAdapter) GetBasePath() string {
	// Get the current working directory as the base path
	dir, err := os.Getwd()
	if err != nil {
		// If there's an error, return empty string (current directory)
		return ""
	}
	return dir
}

func (a *nativeAdapter) GetTasksFromPath(tasksPath string) ([]Tasks, error) {

	// Read file contents
	data, err := os.ReadFile(tasksPath)
	if err != nil {
		return nil, err
	}

	// Parse YAML data
	parser := ymlParser{}
	tasks, err := parser.ParseYAML(data)
	if err != nil {
		return nil, err
	}

	// Wrap in Tasks slice to maintain compatibility
	return []Tasks{tasks}, nil
}

func (a *nativeAdapter) ExecuteCmd(cmd Task) error {
	command := exec.Command(cmd.Command, cmd.Args)
	return command.Run()
}

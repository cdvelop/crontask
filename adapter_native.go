//go:build !wasm

package crontask

import (
	"os"
	"os/exec"

	"github.com/goccy/go-yaml"
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

	// Parse YAML data using go-yaml
	var tasks []Task
	err = yaml.Unmarshal(data, &tasks)
	if err != nil {
		// Try to parse with wrapper (tasks: [...])
		var wrapper struct {
			Tasks []Task `yaml:"tasks"`
		}
		err2 := yaml.Unmarshal(data, &wrapper)
		if err2 != nil || len(wrapper.Tasks) == 0 {
			return nil, err // return original error
		}
		tasks = wrapper.Tasks
	}

	return []Tasks{tasks}, nil
}

func (a *nativeAdapter) ExecuteCmd(cmd Task) error {
	// Split args string to proper arguments array
	args := []string{}
	if cmd.Args != "" {
		// Simple parsing - this could be improved with proper argument parsing
		inQuote := false
		current := ""
		for _, c := range cmd.Args {
			if c == '"' || c == '\'' {
				inQuote = !inQuote
				continue
			}
			if c == ' ' && !inQuote {
				if current != "" {
					args = append(args, current)
					current = ""
				}
				continue
			}
			current += string(c)
		}
		if current != "" {
			args = append(args, current)
		}
	}

	command := exec.Command(cmd.Command, args...)
	return command.Run()
}

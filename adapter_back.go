//go:build !wasm

package crontask

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mileusna/crontab"
)

type crontabAdapter struct {
	ctab *crontab.Crontab
}

func AddNewTasks(tasksPath ...string) (*cronTask, error) {

	newAdapter := &crontabAdapter{ctab: crontab.New()}

	return newCronTask(newAdapter, tasksPath...)
}

func (a *crontabAdapter) AddJob(schedule string, fn any, args ...any) error {
	jobFunc, ok := fn.(func())
	if !ok {
		return fmt.Errorf("invalid function type")
	}
	return a.ctab.AddJob(schedule, jobFunc, args...)
}

func (a crontabAdapter) GetTasksFromPath(tasksPath ...string) ([]Tasks, error) {
	filePath := ""
	if len(tasksPath) > 0 {
		filePath = tasksPath[0]
	}

	if filePath == "" {
		return nil, nil
	}

	// Read file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse YAML data
	parser := YAMLParser{}
	tasks, err := parser.ParseYAML(data)
	if err != nil {
		return nil, err
	}

	// Wrap in Tasks slice to maintain compatibility
	return []Tasks{tasks}, nil
}

func (a *crontabAdapter) ExecuteCmd(cmd Task) error {
	command := exec.Command(cmd.Command, cmd.Args)
	return command.Run()
}

// Check for default configuration file in current directory
func GetDefaultConfigPath() string {
	// Try current directory first
	if _, err := os.Stat("crontasks.yml"); err == nil {
		return "crontasks.yml"
	}

	// Try executable directory next
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		defaultPath := filepath.Join(exeDir, "crontasks.yml")
		if _, err := os.Stat(defaultPath); err == nil {
			return defaultPath
		}
	}

	return ""
}

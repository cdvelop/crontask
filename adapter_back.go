//go:build !wasm

package crontask

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mileusna/crontab"
	"gopkg.in/yaml.v3"
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

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tasksCmd []Tasks
	if err := yaml.Unmarshal(data, &tasksCmd); err != nil {
		return nil, err
	}

	return tasksCmd, nil
}

func (a *crontabAdapter) ExecuteCmd(cmd Task) error {
	command := exec.Command(cmd.Command, cmd.Args)
	return command.Run()
}

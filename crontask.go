package crontask

import "errors"

type cronAdapter interface {
	AddJob(schedule string, fn any, args ...any) error
	GetTasksFromPath(tasksPath ...string) ([]Tasks, error)
}

type cronTask struct {
	adapter cronAdapter
	tasks   []Task
}

// Example of tasksPath: "C:/tasks/tasks.yaml" or "/etc/tasks/tasks.yaml"
// The tasksPath parameter is used to load tasks from a YAML file
func newCronTask(a cronAdapter, tasksPath ...string) (*cronTask, error) {
	var ts []Tasks
	var err error

	if len(tasksPath) > 0 {
		ts, err = a.GetTasksFromPath(tasksPath[0])
		if err != nil {
			return nil, errors.New("newCronTask " + err.Error())
		}
	}

	c := &cronTask{
		adapter: a,
		tasks:   make([]Task, 0),
	}

	for _, t := range ts {
		c.tasks = append(c.tasks, t...)
	}

	return c, nil
}

func (c *cronTask) AddJob(schedule string, fn any, args ...any) error {
	return c.adapter.AddJob(schedule, fn, args...)
}

type Tasks []Task

type Task struct {
	Name     string `yaml:"name"`     // eg: "Backup system"
	Schedule string `yaml:"schedule"` // eg: "0 7 * * 1,4" (2 times a week, monday and thursday)
	Command  string `yaml:"command"`  // eg: "C:\Program Files\FreeFileSync\FreeFileSync.exe"
	Args     string `yaml:"args"`     // eg: "D:\Backup\SystemBackup.ffs_batch"
}

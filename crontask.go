package crontask

type cronAdapter interface {
	AddJob(schedule string, fn any, args ...any) error
	GetTasksFromPath(tasksPath string) ([]Tasks, error)
	ExecuteCmd(cmd Task) error
}

const filePathDefault = "crontasks.yml"

func getDefaultFilePathTasks(tasksPath string) string {
	// 1.yml = min 5 characters
	if len(tasksPath) >= 5 {
		return tasksPath
	}

	return filePathDefault
}

type Tasks []Task

type Task struct {
	Name     string `yaml:"name"`     // eg: "Backup system"
	Schedule string `yaml:"schedule"` // eg: "0 7 * * 1,4" (2 times a week, monday and thursday)
	Command  string `yaml:"command"`  // eg: "C:\Program Files\FreeFileSync\FreeFileSync.exe"
	Args     string `yaml:"args"`     // eg: "D:\Backup\SystemBackup.ffs_batch"
}

type CronTaskEngine struct {
	adapter cronAdapter
	tasks   []Task
}

// NewCronTaskEngine creates a new CronTaskEngine instance.
// It automatically selects the appropriate adapter based on the build environment.
// If tasksPath is provided, it will try to load tasks from that path.
// Example of tasksPath: "C:/tasks/tasks.yaml" or "/etc/tasks/tasks.yaml"
// default: "crontasks.yml"
func NewCronTaskEngine(tasksPath ...string) (*CronTaskEngine, error) {
	// The adapter initialization is handled by build-specific files
	a := newCronAdapter()

	var ts []Tasks
	var err error

	pathTasks := "crontasks.yml"

	if len(tasksPath) > 0 && tasksPath[0] != "" {
		pathTasks = tasksPath[0]
	}

	ts, err = a.GetTasksFromPath(pathTasks)
	if err != nil {
		return nil, newErr("NewCronTaskEngine:", err)
	}

	c := &CronTaskEngine{
		adapter: a,
		tasks:   make([]Task, 0),
	}

	for _, t := range ts {
		c.tasks = append(c.tasks, t...)
	}

	return c, nil
}

// AddJob adds a new scheduled job to the cron task
func (c *CronTaskEngine) AddJob(schedule string, fn any, args ...any) error {
	return c.adapter.AddJob(schedule, fn, args...)
}

// ScheduleAllTasks schedules all loaded tasks to be executed according to their schedule
func (c *CronTaskEngine) ScheduleAllTasks() error {
	for _, task := range c.tasks {
		taskCopy := task // Create a copy to avoid closure issues
		err := c.AddJob(task.Schedule, func() {
			c.adapter.ExecuteCmd(taskCopy)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

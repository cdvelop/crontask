package crontask

type cronAdapter interface {
	AddJob(schedule string, fn any, args ...any) error
	GetTasksFromPath(tasksPath string) ([]Tasks, error)
	ExecuteCmd(cmd Task) error
	GetBasePath() string
	RunAll()
}

const filePathDefault = "crontasks.yml"

type Tasks []Task

type Task struct {
	Name     string `yaml:"name"`     // eg: "Backup system"
	Schedule string `yaml:"schedule"` // eg: "0 7 * * 1,4" (2 times a week, monday and thursday)
	Command  string `yaml:"command"`  // eg: "C:\Program Files\FreeFileSync\FreeFileSync.exe"
	Args     string `yaml:"args"`     // eg: "D:\Backup\SystemBackup.ffs_batch"
}

// Config contains all configuration options for the CronTaskEngine
type Config struct {
	Logger    func(...any) // Logger function
	TasksPath string       // Path to tasks file, default: "crontasks.yml"
}

type CronTaskEngine struct {
	adapter  cronAdapter
	tasks    []Task
	logger   func(...any) // Logger function
	basePath string       // Base path for execution and file lookup
}

// NewCronTaskEngine creates a new CronTaskEngine instance.
// It automatically selects the appropriate adapter based on the build environment.
// Example: NewCronTaskEngine(Config{Logger: log.Printf, TasksPath: "tasks/tasks.yaml"})
func NewCronTaskEngine(config Config) *CronTaskEngine {
	// The adapter initialization is handled by build-specific files
	a := newCronAdapter()

	c := &CronTaskEngine{
		adapter:  a,
		tasks:    make([]Task, 0),
		logger:   config.Logger,
		basePath: a.GetBasePath(), // Get base path from adapter
	}

	// Set default tasks path if not provided
	pathTasks := filePathDefault
	if config.TasksPath != "" {
		pathTasks = config.TasksPath
	}

	ts, err := a.GetTasksFromPath(pathTasks)
	if err != nil {
		c.logger("no task from path:", pathTasks)
	} else {
		for _, t := range ts {
			c.tasks = append(c.tasks, t...)
		}
	}

	return c
}

// AddJob adds a new scheduled job to the cron task
func (c *CronTaskEngine) AddJob(schedule string, fn any, args ...any) error {
	return c.adapter.AddJob(schedule, fn, args...)
}

// ScheduleAllTasks schedules all loaded tasks to be executed according to their schedule
func (c *CronTaskEngine) ScheduleAllTasks() error {

	if len(c.tasks) == 0 {
		return newErr("no tasks to schedule")
	}

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

// RunAll executes all scheduled tasks immediately
func (c *CronTaskEngine) RunAll() {
	c.adapter.RunAll()
}

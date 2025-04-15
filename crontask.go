package crontask

import (
	"fmt"
	"path/filepath"
)

type cronAdapter interface {
	AddProgramTask(schedule string, fn any, args ...any) error
	GetTasksFromPath(tasksPath string) ([]Tasks, error)
	ExecuteCmd(cmd Task) error
	GetBasePath() string // without / eg: "path/to/base"
	RunAllAdapterTasks()
	Log(...any) // Logger function
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
	TasksPath      string // Path to tasks file, default: "crontasks.yml"
	NoAutoSchedule bool   // Set to true to disable automatic task scheduling
	testFolderPath string // Base path for execution and file lookup eg: "test/uc01_test", default: ""
}

type CronTaskEngine struct {
	adapter cronAdapter
	tasks   []Task
	Log     func(...any) // Logger function
}

// NewCronTaskEngine creates a new CronTaskEngine instance.
// It automatically selects the appropriate adapter based on the build environment
// and schedules all tasks by default.
// Examples:
//
//	engine := NewCronTaskEngine()            // Uses all defaults
//	engine := NewCronTaskEngine(Config{})    // Uses all defaults (explicit)
//	engine := NewCronTaskEngine(Config{TasksPath: "custom.yml"}) // Custom config
func NewCronTaskEngine(configs ...Config) *CronTaskEngine {
	// Default config
	config := Config{}

	// Use first config if provided
	if len(configs) > 0 {
		config = configs[0]
	}

	// The adapter initialization is handled by build-specific files
	a := newCronAdapter()

	var testFolderPath string
	if config.testFolderPath != "" {
		testFolderPath = config.testFolderPath
	}

	c := &CronTaskEngine{
		adapter: a,
		tasks:   make([]Task, 0),
		Log:     a.Log,
	}

	// Set default tasks path if not provided
	pathTasks := filePathDefault
	if config.TasksPath != "" {
		pathTasks = config.TasksPath
	}

	fullPath := filepath.Join(a.GetBasePath(), testFolderPath, pathTasks)
	c.Log("Loading tasks from", fullPath)

	ts, err := a.GetTasksFromPath(fullPath)
	if err != nil {
		c.Log("No tasks loaded from path:", fullPath, "Error:", err)
	} else {
		for _, t := range ts {
			c.tasks = append(c.tasks, t...)
		}

		// Display loaded tasks
		for i, task := range c.tasks {
			c.Log(fmt.Sprintf("Task %d: %s (Schedule: %s)", i+1, task.Name, task.Schedule))
		}
	}

	// Auto-schedule tasks unless explicitly disabled
	if !config.NoAutoSchedule {
		if err := c.ScheduleAllTasks(); err != nil {
			c.Log("Error scheduling tasks:", err)
		} else {
			c.Log("All tasks scheduled successfully")
		}
	}

	return c
}

// AddJob adds a new scheduled job to the cron task
func (c *CronTaskEngine) AddTaskSchedule(schedule string, fn any, args ...any) error {
	c.Log("Adding job with schedule:", schedule)
	return c.adapter.AddProgramTask(schedule, fn, args...)
}

// ScheduleAllTasks schedules all loaded tasks to be executed according to their schedule
func (c *CronTaskEngine) ScheduleAllTasks() error {
	if len(c.tasks) == 0 {
		return newErr("no tasks to schedule")
	}

	c.Log("Scheduling", len(c.tasks), "tasks")
	for _, task := range c.tasks {
		taskCopy := task // Create a copy to avoid closure issues
		c.Log("Scheduling task:", task.Name, "with schedule:", task.Schedule)
		err := c.adapter.AddProgramTask(task.Schedule, func() {
			c.Log("Executing scheduled task:", taskCopy.Name)
			c.adapter.ExecuteCmd(taskCopy)
		})
		if err != nil {
			c.Log("Error scheduling task:", task.Name, "Error:", err)
			return err
		}
	}
	return nil
}

// RunAll executes all scheduled tasks immediately
func (c *CronTaskEngine) RunAllTasks() {
	c.Log("Running all scheduled tasks")
	c.adapter.RunAllAdapterTasks()
}

// ExecuteTask executes a specific task by its name
func (c *CronTaskEngine) ExecuteTask(taskName string) error {
	c.Log("Executing task:", taskName)
	for _, task := range c.tasks {
		if task.Name == taskName {
			return c.adapter.ExecuteCmd(task)
		}
	}
	return newErr("task not found: " + taskName)
}

// GetTasks returns a copy of all loaded tasks
func (c *CronTaskEngine) GetTasks() []Task {
	tasksCopy := make([]Task, len(c.tasks))
	copy(tasksCopy, c.tasks)
	return tasksCopy
}

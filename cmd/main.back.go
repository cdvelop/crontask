//go:build !wasm

package main

import (
	"fmt"
	"time"

	"github.com/cdvelop/crontask"
)

func main() {
	// Create a crontask engine with minimal configuration
	// It will automatically:
	// - Load tasks from "crontasks.yml"
	// - Schedule all tasks
	// - Log operations
	cron := crontask.NewCronTaskEngine()

	// Add a programmatic task if needed
	err := cron.AddTaskSchedule("* * * * *", func() {
		fmt.Println("Executing programmatic task every minute:", time.Now())
	})

	if err != nil {
		fmt.Println("Error adding job:", err)
		return
	}

	fmt.Println("Cron server started at", time.Now())
	fmt.Println("Press Ctrl+C to stop")

	// Keep the program running
	select {}
}

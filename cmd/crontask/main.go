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
	crontask.NewCronTaskEngine()

	fmt.Printf("Cron server started %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("Press Ctrl+C to stop")

	// Keep the program running
	select {}
}

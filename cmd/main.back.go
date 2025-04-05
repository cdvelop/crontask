//go:build !wasm

package main

import (
	"fmt"
	"time"

	"github.com/cdvelop/crontask"
)

func main() {
	// Check for default config file
	configPath := crontask.GetDefaultConfigPath()

	cron, err := crontask.AddNewTasks(configPath)
	if err != nil {
		fmt.Println("Error initializing cron:", err)
		return
	}

	// Schedule tasks from config file if any
	if configPath != "" {
		if err := cron.ScheduleAllTasks(); err != nil {
			fmt.Println("Error scheduling tasks:", err)
		} else {
			fmt.Println("Tasks scheduled from", configPath)
		}
	}

	// Add programmatic tasks
	err = cron.AddJob("* * * * *", func() {
		fmt.Println("Ejecutando tarea cada minuto:", time.Now())
	})

	if err != nil {
		fmt.Println("Error al agregar job:", err)
		return
	}

	fmt.Println("iniciado")

	// Mantener el programa en ejecución
	select {}
}

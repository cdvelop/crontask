//go:build !wasm

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cdvelop/crontask"
)

func main() {

	cron := crontask.NewCronTaskEngine(crontask.Config{
		Logger: log.Println,
	})

	// Schedule tasks from config file if any
	if err := cron.ScheduleAllTasks(); err != nil {
		fmt.Println("Error scheduling tasks:", err)
	}

	// Add programmatic tasks
	err := cron.AddJob("* * * * *", func() {
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

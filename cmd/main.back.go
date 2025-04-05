//go:build !wasm

package main

import (
	"fmt"
	"time"

	"github.com/cdvelop/crontask"
)

func main() {
	cron, err := crontask.AddNewTasks()
	if err == nil {

		// Ejemplo de uso con el struct Crontask
		err = cron.AddJob("* * * * *", func() {
			fmt.Println("Ejecutando tarea cada minuto:", time.Now())
		})

		if err != nil {
			fmt.Println("Error al agregar job:", err)
			return
		}

		fmt.Println("iniciado")

	} else {
		fmt.Println(err)
	}
	// Mantener el programa en ejecución
	select {}
}

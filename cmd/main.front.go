//go:build wasm

package main

import (
	"syscall/js"

	"github.com/cdvelop/crontask"
)

func main() {
	// Registrar funciones en JS
	js.Global().Set("addCronJob", js.FuncOf(addCronJob))

	// Mantener el programa en ejecución
	select {}
}

func addCronJob(this js.Value, args []js.Value) any {
	adapter, err := crontask.AddNewTasks()
	if err != nil {
		return err.Error()
	}

	if len(args) < 2 {
		return "Se requieren 2 argumentos: schedule y callback"
	}

	schedule := args[0].String()
	callback := args[1]

	err = adapter.AddJob(schedule, func() {
		callback.Invoke()
	})

	if err != nil {
		return err.Error()
	}

	return nil
}

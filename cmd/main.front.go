//go:build wasm

package main

import (
	"syscall/js"

	"github.com/cdvelop/crontask"
)

func main() {
	// register the functions to be called from JavaScript
	js.Global().Set("addCronJob", js.FuncOf(addCronJob))
	js.Global().Set("loadCronTasks", js.FuncOf(loadCronTasks))

	// Try to load default config
	loadDefaultConfig()

	// maintain the program running
	select {}
}

func loadDefaultConfig() {
	// Try to load from default locations
	adapter, err := crontask.AddNewTasks("crontasks.yml")
	if err == nil {
		adapter.ScheduleAllTasks()
		js.Global().Call("console.log", "Cron tasks loaded from default config")
	}
}

func addCronJob(this js.Value, args []js.Value) any {
	adapter, err := crontask.AddNewTasks()
	if err != nil {
		return err.Error()
	}

	if len(args) < 2 {
		return "required minimum 2 arguments: schedule and callback function eg: addCronJob('* * * * *', function() { console.log('Hello World') })"
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

func loadCronTasks(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return "Se requiere la ruta del archivo de tareas"
	}

	configPath := args[0].String()
	adapter, err := crontask.AddNewTasks(configPath)
	if err != nil {
		return err.Error()
	}

	if err := adapter.ScheduleAllTasks(); err != nil {
		return err.Error()
	}

	return "Tareas programadas correctamente"
}

// filepath: c:\Users\Cesar\Packages\Internal\crontask\taskyml_test.go
package crontask

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTasksFromYaml(t *testing.T) {
	// esta preuna se trata de cargar un archivo yaml y verificar que se carguen las tareas correctamente
	// y que los valores sean correctos.
	// Definir la ruta del directorio de prueba
	testDirPath := filepath.Join("test", "uc02_load_yml_task")

	// Verificar que el archivo existe antes de continuar
	yamlPath := filepath.Join(testDirPath, filePathDefault)
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		t.Fatalf("El archivo de tareas YAML no existe en la ruta: %s", yamlPath)
	}
	// Crear una nueva instancia de CronTaskEngine con la configuración adecuada
	// La ruta se configura de manera que la concatenación simple de cadenas
	// que hace el motor de tareas funcione correctamente.
	// No incluimos separador de directorio al final ya que el motor concatena directamente
	cron := NewCronTaskEngine(Config{
		testFolderPath: testDirPath, // Sin separador al final
	})

	// Verificar que se hayan cargado exactamente 2 tareas
	if len(cron.tasks) != 2 {
		t.Errorf("Se esperaban 2 tareas, pero se encontraron %d", len(cron.tasks))
		return
	}

	// Verificar la primera tarea
	if cron.tasks[0].Name != "task1" {
		t.Errorf("Nombre incorrecto para la tarea 1. Se esperaba 'task1', se obtuvo '%s'", cron.tasks[0].Name)
	}
	if cron.tasks[0].Schedule != "*/5 * * * *" {
		t.Errorf("Programación incorrecta para la tarea 1. Se esperaba '*/5 * * * *', se obtuvo '%s'", cron.tasks[0].Schedule)
	}
	if cron.tasks[0].Command != "echo" {
		t.Errorf("Comando incorrecto para la tarea 1. Se esperaba 'echo', se obtuvo '%s'", cron.tasks[0].Command)
	}
	if cron.tasks[0].Args != "hello world" {
		t.Errorf("Argumentos incorrectos para la tarea 1. Se esperaba 'hello world', se obtuvo '%s'", cron.tasks[0].Args)
	}

	// Verificar la segunda tarea
	if cron.tasks[1].Name != "task2" {
		t.Errorf("Nombre incorrecto para la tarea 2. Se esperaba 'task2', se obtuvo '%s'", cron.tasks[1].Name)
	}
	if cron.tasks[1].Schedule != "0 12 * * *" {
		t.Errorf("Programación incorrecta para la tarea 2. Se esperaba '0 12 * * *', se obtuvo '%s'", cron.tasks[1].Schedule)
	}
	if cron.tasks[1].Command != "ls" {
		t.Errorf("Comando incorrecto para la tarea 2. Se esperaba 'ls', se obtuvo '%s'", cron.tasks[1].Command)
	}
	if cron.tasks[1].Args != "-la" {
		t.Errorf("Argumentos incorrectos para la tarea 2. Se esperaba '-la', se obtuvo '%s'", cron.tasks[1].Args)
	}
}

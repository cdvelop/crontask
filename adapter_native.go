//go:build !wasm

package crontask

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/goccy/go-yaml"
)

// Adaptador para entornos nativos (no-WASM)
type nativeAdapter struct {
	ctab   *crontab
	logger *log.Logger
}

// Inicializador específico para entornos no-WASM
func newCronAdapter() cronAdapter {
	// Abrir o crear archivo log.txt para escritura y append
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Si hay error al abrir el archivo, usar solo stdout
		return &nativeAdapter{
			ctab:   newCrontab(),
			logger: log.New(os.Stdout, "CRONTASK: ", log.Ldate|log.Ltime),
		}
	}

	// Crear un MultiWriter para escribir en stdout y en el archivo
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Crear logger con formato y destino configurados
	logger := log.New(multiWriter, "CRONTASK: ", log.Ldate|log.Ltime|log.Lshortfile)

	return &nativeAdapter{
		ctab:   newCrontab(),
		logger: logger,
	}
}

func (a *nativeAdapter) Log(args ...any) {
	// Usar el logger configurado en lugar de log.Println directamente
	a.logger.Println(args...)
}

func (a *nativeAdapter) AddProgramTask(schedule string, fn any, args ...any) error {
	jobFunc, ok := fn.(func())
	if !ok {
		return newErr("invalid function type")
	}
	return a.ctab.AddJob(schedule, jobFunc, args...)
}

func (a *nativeAdapter) RunAllAdapterTasks() {
	a.ctab.RunAll()
}

func (a *nativeAdapter) GetBasePath() string {
	// Get the current working directory as the base path
	dir, err := os.Getwd()
	if err != nil {
		// If there's an error, return empty string (current directory)
		return ""
	}
	return dir
}

func (a *nativeAdapter) GetTasksFromPath(tasksPath string) ([]Tasks, error) {
	// Read file contents
	data, err := os.ReadFile(tasksPath)
	if err != nil {
		return nil, err
	}

	// Parse YAML data using go-yaml
	var tasks []Task
	err = yaml.Unmarshal(data, &tasks)
	if err != nil {
		// Try to parse with wrapper (tasks: [...])
		var wrapper struct {
			Tasks []Task `yaml:"tasks"`
		}
		err2 := yaml.Unmarshal(data, &wrapper)
		if err2 != nil || len(wrapper.Tasks) == 0 {
			return nil, err // return original error
		}
		tasks = wrapper.Tasks
	}

	return []Tasks{tasks}, nil
}

func (a *nativeAdapter) ExecuteCmd(cmd Task) error {
	// Split args string to proper arguments array
	args := []string{}
	if cmd.Args != "" {
		// Simple parsing - this could be improved with proper argument parsing
		inQuote := false
		current := ""
		for _, c := range cmd.Args {
			if c == '"' || c == '\'' {
				inQuote = !inQuote
				continue
			}
			if c == ' ' && !inQuote {
				if current != "" {
					args = append(args, current)
					current = ""
				}
				continue
			}
			current += string(c)
		}
		if current != "" {
			args = append(args, current)
		}
	}

	// Expand environment variables in command and args
	command := os.ExpandEnv(cmd.Command)
	expandedArgs := make([]string, len(args))
	for i, arg := range args {
		expandedArgs[i] = os.ExpandEnv(arg)
	}

	// Use exec.Command directly
	execCmd := exec.Command(command, expandedArgs...)

	// Capture command output for better debugging
	output, err := execCmd.CombinedOutput()
	if err != nil {
		a.Log("Command execution failed:", err, "Output:", string(output))
		return err
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("*- [%s] %v completed.\n%s\n", currentTime, cmd.Name, string(output))
	return nil
}

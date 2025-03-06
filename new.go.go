package crontask

import (
	"os"

	"github.com/mileusna/crontab"
	"gopkg.in/yaml.v3"
)

// Crontask maneja la programación de tareas
type Crontask struct {
	ctab *crontab.Crontab
}

// New crea una nueva instancia de Crontask
func New() *Crontask {
	return &Crontask{
		ctab: crontab.New(),
	}
}

// Add agrega una tarea programada
func (c *Crontask) Add(schedule string, job func()) error {
	return c.ctab.AddJob(schedule, job)
}

type TaskConfig struct {
	Tasks []struct {
		Schedule string `yaml:"schedule"`
		Command  string `yaml:"command"`
	} `yaml:"tasks"`
}

func LoadConfig(filename string) (*TaskConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config TaskConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

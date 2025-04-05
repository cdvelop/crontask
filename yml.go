package crontask

import (
	"bytes"
	"errors"
	"regexp"
)

// YAMLParser handles parsing of YAML content for both frontend and backend
type YAMLParser struct{}

// ParseYAML parses YAML content bytes into Tasks
func (p YAMLParser) ParseYAML(data []byte) (Tasks, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, errors.New("empty YAML content")
	}

	// Parse YAML content with regex for best performance
	return p.parseWithRegex(data)
}

// parseWithRegex uses regular expressions to extract tasks from YAML formats
func (p YAMLParser) parseWithRegex(data []byte) (Tasks, error) {
	var tasks []Task

	// Pattern for a task with quoted or unquoted values
	pattern := regexp.MustCompile(`(?m)- *name: *["']?([^"'\n]+)["']?\n *schedule: *["']?([^"'\n]+)["']?\n *command: *["']?([^"'\n]+)["']?\n *(?:args: *["']?([^"'\n]*)["']?)?`)

	// Direct list pattern
	matches := pattern.FindAllSubmatch(data, -1)

	// If no direct matches, try looking inside a tasks: block
	if len(matches) == 0 {
		// Check if there's a tasks: section and try inside it
		tasksSection := regexp.MustCompile(`(?s)tasks:(.*$)`).FindSubmatch(data)
		if len(tasksSection) > 1 {
			matches = pattern.FindAllSubmatch(tasksSection[1], -1)
		}
	}

	// Process all matches
	for _, match := range matches {
		if len(match) >= 4 {
			task := Task{
				Name:     string(match[1]),
				Schedule: string(match[2]),
				Command:  string(match[3]),
			}

			// Args are optional
			if len(match) >= 5 && len(match[4]) > 0 {
				task.Args = string(match[4])
			}

			tasks = append(tasks, task)
		}
	}

	// Check if we found any tasks
	if len(tasks) == 0 {
		return nil, errors.New("no valid tasks found in YAML")
	}

	return tasks, nil
}

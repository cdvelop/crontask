package crontask

import (
	"testing"
)

func TestYmlParser(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    int // Expected number of tasks
		wantErr bool
	}{
		{
			name: "Basic task list",
			yaml: `
- name: "task1"
  schedule: "* * * * *"
  command: "echo"
  args: "hello world"
- name: "task2"
  schedule: "*/5 * * * *"
  command: "ls"
  args: "-la"
`,
			want:    2,
			wantErr: false,
		},
		{
			name: "Tasks with wrapper",
			yaml: `
tasks:
  - name: "task1"
    schedule: "* * * * *"
    command: "echo"
    args: "hello world"
  - name: "task2"
    schedule: "*/5 * * * *"
    command: "ls"
    args: "-la"
`,
			want:    2,
			wantErr: false,
		},
		{
			name: "Single task",
			yaml: `
- name: "task1"
  schedule: "* * * * *"
  command: "echo"
  args: "hello world"
`,
			want:    1,
			wantErr: false,
		},
		{
			name:    "Empty content",
			yaml:    "",
			want:    0,
			wantErr: true,
		},
		{
			name: "Invalid YAML (missing required fields)",
			yaml: `
- name: "task1"
  schedule: "* * * * *"
  # missing command
`,
			want:    0,
			wantErr: true,
		},
		{
			name: "With quotes",
			yaml: `
- name: "my task"
  schedule: '*/10 * * * *'
  command: "/bin/bash"
  args: "-c 'echo hello'"
`,
			want:    1,
			wantErr: false,
		},
	}

	parser := ymlParser{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseYAML([]byte(tt.yaml))

			// Check error expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If expecting error, no need to check content
			if tt.wantErr {
				return
			}

			// Check number of tasks
			if len(got) != tt.want {
				t.Errorf("ParseYAML() got %d tasks, want %d", len(got), tt.want)
			}

			// Validate first task if any expected
			if tt.want > 0 {
				task := got[0]
				if task.Name == "" || task.Schedule == "" || task.Command == "" {
					t.Errorf("ParseYAML() got incomplete task: %+v", task)
				}
			}
		})
	}
}

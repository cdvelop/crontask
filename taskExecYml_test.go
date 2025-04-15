package crontask

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

const testAppName = "test_app"
const testDirPath = "test/uc03_exec_yml_task"
const testFileContent = "This is a test file content"

// createTestApp creates a Go application that can create or delete files
func createTestApp(t *testing.T) string {
	// Determine the executable extension based on the OS
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	appPath := filepath.Join(testDirPath, testAppName+ext)

	// Check if the test app already exists
	if _, err := os.Stat(appPath); err == nil {
		return appPath // App already exists, no need to rebuild
	}

	// Create test directory if it doesn't exist
	if err := os.MkdirAll(testDirPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a temporary Go file for the test application
	tempAppFile := filepath.Join(testDirPath, "main.go")
	appCode := `package main

import (
	"fmt"
	"os"

)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: test_app [+/-] [filename] [content (if +)]")
		os.Exit(1)
	}

	operation := os.Args[1]
	filename := os.Args[2]

	if operation == "+" {
		if len(os.Args) < 4 {
			fmt.Println("Error: Content required for file creation")
			os.Exit(1)
		}
		content := os.Args[3]
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File created: %s\n", filename)
	} else if operation == "-" {
		err := os.Remove(filename)
		if err != nil {
			fmt.Printf("Error removing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File removed: %s\n", filename)
	} else {
		fmt.Println("Invalid operation. Use + to create or - to delete")
		os.Exit(1)
	}
}
`
	if err := os.WriteFile(tempAppFile, []byte(appCode), 0644); err != nil {
		t.Fatalf("Failed to create test app source: %v", err)
	}

	// Compile the Go application
	buildCmd := exec.Command("go", "build", "-o", appPath, tempAppFile)
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build test app: %v", err)
	}

	return appPath
}

// createYamlFile creates a YAML file with tasks to create and delete a test file
func createYamlFile(t *testing.T, appPath string) string {
	yamlPath := filepath.Join(testDirPath, filePathDefault)
	testFilePath := filepath.Join(testDirPath, testFileName)

	// Create the YAML content
	yamlContent := `- name: "create_file"
  schedule: "*/1 * * * *"
  command: "` + appPath + `"
  args: "+ ` + testFilePath + ` \"` + testFileContent + `\""
- name: "delete_file"
  schedule: "*/1 * * * *"
  command: "` + appPath + `"
  args: "- ` + testFilePath + `"
`

	// Create the YAML file
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create YAML file: %v", err)
	}

	return yamlPath
}

// TestTaskExecution implements UC03 (Use Case 03): Testing YAML task execution functionality
// This test verifies the CronTaskEngine's ability to:
// 1. Load tasks from a YAML configuration file
// 2. Execute a specific task by name ("create_file") which creates a test file with content
// 3. Verify the file was created with the correct content
// 4. Execute another task ("delete_file") which removes the test file
// 5. Verify the file was properly deleted
// This demonstrates the end-to-end capability of the task execution system
// with file operations as the test actions.
func TestTaskExecution(t *testing.T) {
	// Create the test app
	appPath := createTestApp(t)

	// Create the YAML file with tasks
	createYamlFile(t, appPath)

	// Initialize the CronTaskEngine
	cron := NewCronTaskEngine(Config{
		Logger:         t.Log,
		TasksPath:      filePathDefault,
		testFolderPath: testDirPath,
	})

	// Verify that tasks were loaded correctly
	if len(cron.tasks) != 2 {
		t.Fatalf("Expected 2 tasks, but got %d", len(cron.tasks))
	}

	// Verify the test file doesn't exist initially
	testFilePath := filepath.Join(testDirPath, testFileName)
	if _, err := os.Stat(testFilePath); err == nil {
		// Remove it if it exists from a previous test run
		os.Remove(testFilePath)
	}

	// Execute the create_file task
	t.Log("Executing create_file task")
	err := cron.ExecuteTask("create_file")
	if err != nil {
		t.Fatalf("Failed to execute create_file task: %v", err)
	}

	// Wait briefly to ensure file creation
	time.Sleep(100 * time.Millisecond)

	// Verify the file was created
	fileContents, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read test file after creation: %v", err)
	}

	if string(fileContents) != testFileContent {
		t.Errorf("File content mismatch. Expected: %s, Got: %s", testFileContent, string(fileContents))
	}

	// Execute the delete_file task
	t.Log("Executing delete_file task")
	err = cron.ExecuteTask("delete_file")
	if err != nil {
		t.Fatalf("Failed to execute delete_file task: %v", err)
	}

	// Wait briefly to ensure file deletion
	time.Sleep(100 * time.Millisecond)

	// Verify the file was deleted
	if _, err := os.Stat(testFilePath); err == nil {
		t.Errorf("File still exists after delete task execution")
	}
}

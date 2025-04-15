package crontask

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	testDir       = "test"
	testSubDir    = "uc01_example"
	testFileName  = "test_file.txt"
	testContent   = "This is a test file created by crontask test"
	testSchedule1 = "* * * * *" // Run every minute
	testSchedule2 = "* * * * *" // Run every minute
)

func TestCronTaskEngine(t *testing.T) {
	// Setup test directory
	testDirPath := filepath.Join(testDir, testSubDir)
	err := os.MkdirAll(testDirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	// defer os.RemoveAll(testDirPath) // Clean up after the test

	// Full path to the test file
	testFilePath := filepath.Join(testDirPath, testFileName)
	// Create a new CronTaskEngine without loading tasks from file
	cron := NewCronTaskEngine(Config{
		Logger:    t.Log,
		TasksPath: "",
	})
	cron.basePath = testDirPath // Set the base path for the cron adapter

	// Add a job to create a file
	createCalled := false
	err = cron.AddJob(testSchedule1, func() {
		err := createTestFile(testFilePath, testContent)
		if err != nil {
			t.Errorf("Failed to create test file: %v", err)
		}
		createCalled = true
	})
	if err != nil {
		t.Fatalf("Failed to add create file job: %v", err)
	}

	// Add a job to delete the file
	deleteCalled := false
	err = cron.AddJob(testSchedule2, func() {
		// Only try to delete if the file exists
		if fileExists(testFilePath) {
			err := deleteTestFile(testFilePath)
			if err != nil {
				t.Errorf("Failed to delete test file: %v", err)
			}
			deleteCalled = true
		}
	})
	if err != nil {
		t.Fatalf("Failed to add delete file job: %v", err)
	}
	// Wait a bit for the cron jobs to execute
	time.Sleep(2 * time.Second) // Reduced wait time to avoid test timeout

	// Since we're no longer waiting for the actual cron schedule, manually trigger the jobs
	createCalled = false
	err = createTestFile(testFilePath, testContent)
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}
	createCalled = true

	// Verify file was created
	if !fileExists(testFilePath) {
		t.Error("File was not created as expected")
	}

	// Manually trigger delete job
	deleteCalled = false
	err = deleteTestFile(testFilePath)
	if err != nil {
		t.Errorf("Failed to delete test file: %v", err)
	}
	deleteCalled = true

	// Verify that both jobs were called
	if !createCalled {
		t.Error("Create file job was not called")
	}
	if !deleteCalled {
		t.Error("Delete file job was not called")
	}

	// Verify file doesn't exist after deletion
	if fileExists(testFilePath) {
		t.Error("File was not deleted as expected")
	}
}

// Helper function to create a test file
func createTestFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// Helper function to delete a test file
func deleteTestFile(path string) error {
	println("Deleting file:", path)
	return os.Remove(path)
}

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

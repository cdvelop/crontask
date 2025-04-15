package crontask

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

const (
	testDir       = "test"
	testSchedule1 = "* * * * *" // Run every minute
	testSchedule2 = "* * * * *" // Run every minute
	testFileName  = "test_file.txt"
	testContent   = "This is a test file created by crontask test"
)

func TestCronTaskEngineWithRunAll(t *testing.T) {
	// en este test se pretende configurar 2 funcions la preimra crea un archivo y la segunda lo elimina
	// una funcio deberia ejecutarce primero y la segunda deberia ejecutarse 1 segundo despues
	testDirPath := filepath.Join(testDir, "uc01_test")
	err := os.MkdirAll(testDirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDirPath) // Clean up after test

	// Full path to the test file
	testFilePath := filepath.Join(testDirPath, testFileName)

	// Create a new CronTaskEngine
	cron := NewCronTaskEngine(Config{
		testFolderPath: testDirPath, // Set the base path for the cron adapter
	})

	// Use sync.WaitGroup to wait for jobs to complete
	var wg sync.WaitGroup
	wg.Add(2)

	// Add create file job
	createCalled := false
	err = cron.AddTaskSchedule(testSchedule1, func() {
		defer wg.Done()
		err := createTestFile(testFilePath, testContent)
		if err != nil {
			t.Errorf("Failed to create test file: %v", err)
		}
		createCalled = true
	})
	if err != nil {
		t.Fatalf("Failed to add create file job: %v", err)
	}
	// Add delete file job - keep it simple
	deleteCalled := false
	err = cron.AddTaskSchedule(testSchedule2, func() {
		// Ensure the file exists before trying to delete it
		// Add 500ms wait to ensure file creation completes first
		time.Sleep(500 * time.Millisecond)

		if fileExists(testFilePath) {
			err := deleteTestFile(testFilePath)
			if err != nil {
				t.Errorf("Failed to delete test file: %v", err)
			}
			deleteCalled = true
		} else {
			t.Logf("Warning: File does not exist yet when trying to delete it")
			// Try again after another wait
			time.Sleep(1 * time.Second)
			if fileExists(testFilePath) {
				err := deleteTestFile(testFilePath)
				if err != nil {
					t.Errorf("Failed to delete test file in second attempt: %v", err)
				} else {
					deleteCalled = true
				}
			}
		}
		wg.Done()
	})
	if err != nil {
		t.Fatalf("Failed to add delete file job: %v", err)
	}

	// Execute all scheduled jobs immediately without waiting for cron schedule
	cron.RunAllTasks()

	// Wait for both jobs to complete with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Jobs completed successfully
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for jobs to complete")
	}

	// Verify jobs were called
	if !createCalled {
		t.Error("Create file job was not executed")
	}
	if !deleteCalled {
		t.Error("Delete file job was not executed")
	}

	// Verify file was created and then deleted
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
	return os.Remove(path)
}

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

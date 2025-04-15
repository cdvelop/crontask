package crontask_test // Use a test package name (e.g., main_test or your_package_test)

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func TestInstallScript(t *testing.T) {

	const installScriptPath = "./build.sh"

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// En Windows, intentamos primero con Git Bash
		gitBashPath := "C:\\Program Files\\Git\\bin\\bash.exe"
		if _, err := os.Stat(gitBashPath); err == nil {
			cmd = exec.Command(gitBashPath, "-c", installScriptPath)
		} else {
			// Si no encuentra Git Bash, intenta con WSL
			cmd = exec.Command("wsl", installScriptPath)
		}
	default:
		// Para Linux y macOS
		cmd = exec.Command("bash", installScriptPath)
	}

	// Configurar el directorio de trabajo
	cmd.Dir = "." // Asegúrate de que esto apunta al directorio correcto

	// Capturar tanto stdout como stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Failed to execute install.sh: %v\nOutput: %s", err, string(output))
	}
}

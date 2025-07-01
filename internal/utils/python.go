package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lmagdanello/bluebanquise-installer/internal/system"
)

// DownloadRequirements downloads Python packages without installing them.
func DownloadRequirements(requirements []string, downloadPath string) error {
	LogInfo("Downloading Python requirements", "requirements", requirements, "path", downloadPath)

	if len(requirements) == 0 {
		LogError("No requirements provided", nil)
		return fmt.Errorf("no requirements provided")
	}

	// Create download directory
	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		LogError("Failed to create download directory", err, "path", downloadPath)
		return fmt.Errorf("failed to create download directory: %v", err)
	}

	// Create temporary requirements.txt file
	requirementsFile := filepath.Join(downloadPath, "requirements.txt")
	requirementsContent := strings.Join(requirements, "\n")
	if err := os.WriteFile(requirementsFile, []byte(requirementsContent), 0644); err != nil {
		LogError("Failed to create requirements.txt", err, "file", requirementsFile)
		return fmt.Errorf("failed to create requirements.txt: %v", err)
	}

	LogInfo("Created requirements.txt", "file", requirementsFile, "content", requirementsContent)

	// Get the correct Python command for this OS
	pythonCmd, err := system.GetPythonCommand()
	if err != nil {
		LogError("Failed to get Python command", err)
		return fmt.Errorf("failed to get Python command: %v", err)
	}

	// Download packages using the OS-specific Python
	LogCommand(pythonCmd, "-m", "pip", "download", "-r", requirementsFile, "-d", downloadPath)
	cmd := exec.Command(pythonCmd, "-m", "pip", "download", "-r", requirementsFile, "-d", downloadPath)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		LogError("Failed to download requirements", err, "requirements", requirements, "path", downloadPath, "output", string(output))
		return fmt.Errorf("failed to download requirements: %v, output: %s", err, string(output))
	}

	LogInfo("pip download completed", "output", string(output))

	// Verify that packages were downloaded
	entries, err := os.ReadDir(downloadPath)
	if err != nil {
		LogError("Failed to read download directory", err, "path", downloadPath)
		return fmt.Errorf("failed to read download directory: %v", err)
	}

	packageCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(name, ".whl") || strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".tgz") {
				packageCount++
				LogInfo("Downloaded package", "name", name)
			}
		}
	}

	if packageCount == 0 {
		LogError("No packages were downloaded", nil, "path", downloadPath, "entries", len(entries))
		return fmt.Errorf("no packages were downloaded to %s", downloadPath)
	}

	LogInfo("Requirements downloaded successfully", "path", downloadPath, "requirements", requirements, "packages", packageCount)
	return nil
}

// InstallRequirementsOffline installs Python packages from local directory.
func InstallRequirementsOffline(venvPath, requirementsPath string) error {
	LogInfo("Installing Python requirements offline", "venv", venvPath, "requirements_path", requirementsPath)

	if _, err := os.Stat(requirementsPath); os.IsNotExist(err) {
		LogError("Requirements path does not exist", err, "path", requirementsPath)
		return fmt.Errorf("requirements path does not exist: %s", requirementsPath)
	}

	requirementsFile := filepath.Join(requirementsPath, "requirements.txt")

	// Check if requirements.txt exists
	if _, err := os.Stat(requirementsFile); os.IsNotExist(err) {
		LogError("requirements.txt not found", err, "file", requirementsFile)
		return fmt.Errorf("requirements.txt not found: %s", requirementsFile)
	}

	// List contents of requirements directory for debug
	entries, err := os.ReadDir(requirementsPath)
	if err != nil {
		LogError("Cannot read requirements directory", err, "path", requirementsPath)
		return fmt.Errorf("cannot read requirements directory: %v", err)
	}

	LogInfo("Requirements directory contents", "path", requirementsPath, "entries", len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			LogInfo("Requirements file", "name", entry.Name())
		}
	}

	// Install packages from local directory using the OS-specific Python
	pythonCmd, err := system.GetPythonCommand()
	if err != nil {
		LogError("Failed to get Python command", err)
		return fmt.Errorf("failed to get Python command: %v", err)
	}

	args := []string{"-m", "pip", "install", "--no-index", "--find-links", requirementsPath, "-r", requirementsFile}

	fmt.Printf("Installing Python packages from local directory: %s\n", requirementsPath)
	LogCommand(pythonCmd, args...)
	cmd := exec.Command(pythonCmd, args...)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		LogError("Failed to install requirements offline", err, "venv", venvPath, "requirements_path", requirementsPath, "output", string(output))
		return fmt.Errorf("failed to install requirements offline: %v, output: %s", err, string(output))
	}

	LogInfo("pip install completed", "output", string(output))
	LogInfo("Requirements installed offline successfully", "venv", venvPath, "requirements_path", requirementsPath)
	return nil
}

// InstallRequirements installs Python packages in a virtual environment.
func InstallRequirements(venvPath string, requirements []string) error {
	LogInfo("Installing Python requirements", "venv", venvPath, "requirements", requirements)

	if len(requirements) == 0 {
		LogError("No requirements provided", nil)
		return fmt.Errorf("no requirements provided")
	}

	python3 := filepath.Join(venvPath, "bin", "python3")

	args := append([]string{"-m", "pip", "install", "--upgrade", "pip"}, requirements...)

	fmt.Printf("Installing Python packages: %s\n", strings.Join(requirements, " "))
	LogCommand(python3, args...)
	cmd := exec.Command(python3, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		LogError("Failed to install python packages", err, "venv", venvPath, "requirements", requirements)
		return fmt.Errorf("failed to install python packages: %v", err)
	}

	LogInfo("Python requirements installed successfully", "venv", venvPath, "requirements", requirements)
	return nil
}

// RHEL 7.
func ExportRHPython38(userHome string) error {
	LogInfo("Exporting RHEL7 Python 3.8 environment", "home", userHome)

	bashrc := filepath.Join(userHome, ".bashrc")
	lines := []string{
		"export LD_LIBRARY_PATH=/opt/rh/rh-python38/root/usr/lib64:$LD_LIBRARY_PATH",
		"export MANPATH=/opt/rh/rh-python38/root/usr/share/man:$MANPATH",
		"export PATH=/opt/rh/rh-python38/root/usr/local/bin:/opt/rh/rh-python38/root/usr/bin:$PATH",
		"export PKG_CONFIG_PATH=/opt/rh/rh-python38/root/usr/lib64/pkgconfig:$PKG_CONFIG_PATH",
		"export XDG_DATA_DIRS=/opt/rh/rh-python38/root/usr/share:$XDG_DATA_DIRS",
		"export X_SCLS=\"rh-python38 \"",
	}

	LogInfo("Opening .bashrc for RHEL7 Python configuration", "file", bashrc)
	f, err := os.OpenFile(bashrc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		LogError("Failed to open .bashrc for RHEL7 Python configuration", err, "file", bashrc)
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			LogWarning("Failed to close file", "error", closeErr)
		}
	}()

	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			LogError("Failed to write RHEL7 Python configuration line", err, "line", line)
			return err
		}
	}

	LogInfo("RHEL7 Python 3.8 environment exported successfully", "home", userHome)
	return nil
}

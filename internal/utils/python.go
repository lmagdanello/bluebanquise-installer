package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	// Download packages using pip
	LogCommand("pip", "download", "-r", requirementsFile, "-d", downloadPath)
	cmd := exec.Command("pip", "download", "-r", requirementsFile, "-d", downloadPath)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		LogError("Failed to download requirements", err, "requirements", requirements, "path", downloadPath)
		return fmt.Errorf("failed to download requirements: %v", err)
	}

	LogInfo("Requirements downloaded successfully", "path", downloadPath, "requirements", requirements)
	return nil
}

// InstallRequirementsOffline installs Python packages from local directory.
func InstallRequirementsOffline(venvPath, requirementsPath string) error {
	LogInfo("Installing Python requirements offline", "venv", venvPath, "requirements_path", requirementsPath)

	if _, err := os.Stat(requirementsPath); os.IsNotExist(err) {
		LogError("Requirements path does not exist", err, "path", requirementsPath)
		return fmt.Errorf("requirements path does not exist: %s", requirementsPath)
	}

	pip := filepath.Join(venvPath, "bin", "pip")
	requirementsFile := filepath.Join(requirementsPath, "requirements.txt")

	// Check if requirements.txt exists
	if _, err := os.Stat(requirementsFile); os.IsNotExist(err) {
		LogError("requirements.txt not found", err, "file", requirementsFile)
		return fmt.Errorf("requirements.txt not found: %s", requirementsFile)
	}

	// Install packages from local directory
	args := []string{"install", "--no-index", "--find-links", requirementsPath, "-r", requirementsFile}

	fmt.Printf("Installing Python packages from local directory: %s\n", requirementsPath)
	LogCommand(pip, args...)
	cmd := exec.Command(pip, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		LogError("Failed to install requirements offline", err, "venv", venvPath, "requirements_path", requirementsPath)
		return fmt.Errorf("failed to install requirements offline: %v", err)
	}

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

	pip := filepath.Join(venvPath, "bin", "pip")

	args := append([]string{"install", "--upgrade", "pip"}, requirements...)

	fmt.Printf("Installing Python packages: %s\n", strings.Join(requirements, " "))
	LogCommand(pip, args...)
	cmd := exec.Command(pip, args...)
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
	defer f.Close()

	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			LogError("Failed to write RHEL7 Python configuration line", err, "line", line)
			return err
		}
	}

	LogInfo("RHEL7 Python 3.8 environment exported successfully", "home", userHome)
	return nil
}

// Ubuntu 20.04.
func BuildPython311FromSource() error {
	LogInfo("Building Python 3.11 from source for Ubuntu 20.04")
	fmt.Println("Building Python 3.11 from source...")

	cmds := [][]string{
		{"wget", "https://www.python.org/ftp/python/3.11.4/Python-3.11.4.tgz"},
		{"tar", "-xf", "Python-3.11.4.tgz"},
		{"bash", "-c", "cd Python-3.11.4 && ./configure --enable-optimizations --with-ensurepip=install"},
		{"bash", "-c", "cd Python-3.11.4 && make -j"},
		{"bash", "-c", "cd Python-3.11.4 && make altinstall"},
		{"update-alternatives", "--install", "/usr/bin/python3", "python3", "/usr/local/bin/python3.11", "3"},
		{"update-alternatives", "--install", "/usr/bin/python", "python", "/usr/local/bin/python3.11", "3"},
		{"update-alternatives", "--install", "/usr/bin/pip3", "pip3", "/usr/local/bin/pip3.11", "3"},
		{"update-alternatives", "--install", "/usr/bin/pip", "pip", "/usr/local/bin/pip3.11", "3"},
	}

	for i, args := range cmds {
		LogCommand(args[0], args[1:]...)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Run(); err != nil {
			LogError("Failed to execute Python build command", err, "step", i+1, "command", args)
			return fmt.Errorf("failed to execute command: %v", args)
		}
		LogInfo("Python build step completed", "step", i+1, "command", args)
	}

	LogInfo("Python 3.11 built from source successfully")
	return nil
}

// OpenSUSE.
func LinkPython311AsDefault() error {
	LogInfo("Linking python3.11 as default in OpenSUSE")
	fmt.Println("Linking python3.11 as default in opensuse...")

	cmds := [][]string{
		{"update-alternatives", "--install", "/usr/bin/python3", "python3", "/usr/bin/python3.11", "3"},
		{"update-alternatives", "--install", "/usr/bin/python", "python", "/usr/bin/python3.11", "3"},
		{"update-alternatives", "--install", "/usr/bin/pip3", "pip3", "/usr/bin/pip3.11", "3"},
		{"update-alternatives", "--install", "/usr/bin/pip", "pip", "/usr/bin/pip3.11", "3"},
	}

	for i, args := range cmds {
		LogCommand(args[0], args[1:]...)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Run(); err != nil {
			LogError("Failed to link python3.11", err, "step", i+1, "command", args)
			return fmt.Errorf("failed to link python3.11: %v", err)
		}
		LogInfo("Python link step completed", "step", i+1, "command", args)
	}

	LogInfo("Python 3.11 linked as default successfully")
	return nil
}

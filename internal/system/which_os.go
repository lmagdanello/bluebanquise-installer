package system

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

const defaultPythonCmd = "/usr/bin/python3"

// OSMapping maps OS IDs to BlueBanquise compatible names.
var OSMapping = map[string]string{
	"rhel":          "rhel",
	"centos":        "rhel",
	"rocky":         "rhel",
	"almalinux":     "rhel",
	"ubuntu":        "ubuntu",
	"debian":        "debian",
	"opensuse-leap": "opensuse-leap",
	"sles":          "opensuse-leap",
}

func DetectOS() (string, string, error) {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		slog.Error("Error detecting OS", "error", err)
		return "", "", err
	}

	var name, version string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			name = strings.TrimPrefix(line, "ID=")
			name = strings.Trim(name, "\"")
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.TrimPrefix(line, "VERSION_ID=")
			version = strings.Trim(version, "\"")
		}
	}

	// Map OS ID to BlueBanquise compatible name
	if mappedName, exists := OSMapping[name]; exists {
		name = mappedName
	}

	// Handle version mapping for RHEL derivatives
	if name == "rhel" {
		switch {
		case strings.Contains(strings.ToLower(version), "stream"):
			// Extract version number from stream
			if strings.Contains(version, "8") {
				version = "8"
			} else if strings.Contains(version, "9") {
				version = "9"
			}
		default:
			// Extract major version (e.g., 9.3 -> 9)
			if strings.Contains(version, ".") {
				version = strings.Split(version, ".")[0]
			}
		}
	}

	return name, version, nil
}

// GetPythonCommand determines the correct Python command based on the operating system.
func GetPythonCommand() (string, error) {
	// Detect OS to determine the correct Python command
	osID, version, err := DetectOS()
	if err != nil {
		slog.Error("Failed to detect OS", "error", err)
		return "", err
	}

	// Determine Python command based on OS
	var pythonCmd string
	switch osID {
	case "rhel":
		switch version {
		case "7":
			pythonCmd = "/opt/rh/rh-python38/root/usr/bin/python3"
		case "8":
			pythonCmd = "/usr/bin/python3.9"
		case "9":
			// Try multiple Python versions for RHEL9
			pythonVersions := []string{
				"/usr/bin/python3.12",
				"/usr/bin/python3.11",
				"/usr/bin/python3.10",
				"/usr/bin/python3.9",
				"/usr/bin/python3",
			}
			for _, version := range pythonVersions {
				if _, err := os.Stat(version); err == nil {
					pythonCmd = version
					break
				}
			}
			if pythonCmd == "" {
				pythonCmd = defaultPythonCmd
			}
		default:
			pythonCmd = defaultPythonCmd
		}
	case "opensuse-leap":
		pythonCmd = "/usr/bin/python3.11"
	default:
		pythonCmd = defaultPythonCmd
	}

	// Verify the Python command exists
	if _, err := os.Stat(pythonCmd); os.IsNotExist(err) {
		slog.Error("Python command not found", "error", err, "python_cmd", pythonCmd)
		return "", err
	}

	slog.Info("Using Python command", "python_cmd", pythonCmd, "os", osID, "version", version)
	return pythonCmd, nil
}

// BuildPython311FromSource builds Python 3.11 from source for Ubuntu 20.04.
func BuildPython311FromSource() error {
	slog.Info("Building Python 3.11 from source for Ubuntu 20.04")
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
		slog.Info("Executing Python build command", "step", i+1, "command", args)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Run(); err != nil {
			slog.Error("Failed to execute Python build command", "error", err, "step", i+1, "command", args)
			return fmt.Errorf("failed to execute command: %v", args)
		}
		slog.Info("Python build step completed", "step", i+1, "command", args)
	}

	slog.Info("Python 3.11 built from source successfully")
	return nil
}

// LinkPython311AsDefault links python3.11 as default in OpenSUSE.
func LinkPython311AsDefault() error {
	slog.Info("Linking python3.11 as default in OpenSUSE")
	fmt.Println("Linking python3.11 as default in opensuse...")

	cmds := [][]string{
		{"update-alternatives", "--install", "/usr/bin/python3", "python3", "/usr/bin/python3.11", "3"},
		{"update-alternatives", "--install", "/usr/bin/python", "python", "/usr/bin/python3.11", "3"},
		{"update-alternatives", "--install", "/usr/bin/pip3", "pip3", "/usr/bin/pip3.11", "3"},
		{"update-alternatives", "--install", "/usr/bin/pip", "pip", "/usr/bin/pip3.11", "3"},
	}

	for i, args := range cmds {
		slog.Info("Executing Python link command", "step", i+1, "command", args)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Run(); err != nil {
			slog.Error("Failed to link python3.11", "error", err, "step", i+1, "command", args)
			return fmt.Errorf("failed to link python3.11: %v", err)
		}
		slog.Info("Python link step completed", "step", i+1, "command", args)
	}

	slog.Info("Python 3.11 linked as default successfully")
	return nil
}

package utils

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SystemCheck verifies if the system has the necessary prerequisites.
func SystemCheck() error {
	LogInfo("Starting system prerequisites check")

	checks := []struct {
		name  string
		check func() error
	}{
		{"root access", checkRootAccess},
		{"python3", checkPython3},
		{"package manager", checkPackageManager},
		{"internet connectivity", checkInternetConnectivity},
	}

	for _, c := range checks {
		LogInfo("Running system check", "check", c.name)
		fmt.Printf("Checking %s... ", c.name)
		if err := c.check(); err != nil {
			LogError(fmt.Sprintf("%s check failed", c.name), err)
			fmt.Printf("FAILED: %v\n", err)
			return fmt.Errorf("%s check failed: %v", c.name, err)
		}
		LogInfo(fmt.Sprintf("%s check passed", c.name))
		fmt.Println("OK")
	}

	LogInfo("All system checks passed")
	return nil
}

func checkRootAccess() error {
	LogInfo("Checking root access")
	if os.Geteuid() != 0 {
		LogError("Root access check failed", nil, "euid", os.Geteuid())
		return fmt.Errorf("root access required")
	}
	LogInfo("Root access confirmed")
	return nil
}

func checkPython3() error {
	LogInfo("Checking Python3 availability")
	if _, err := exec.LookPath("python3"); err != nil {
		LogError("Python3 not found in PATH", err)
		return fmt.Errorf("python3 not found in PATH")
	}
	LogInfo("Python3 found in PATH")
	return nil
}

func checkPackageManager() error {
	LogInfo("Checking package manager availability")
	candidates := []string{"apt-get", "dnf", "yum", "zypper"}
	for _, tool := range candidates {
		if _, err := exec.LookPath(tool); err == nil {
			LogInfo("Package manager found", "manager", tool)
			return nil
		}
	}
	LogError("No supported package manager found", nil, "managers", candidates)
	return fmt.Errorf("no supported package manager found")
}

func checkInternetConnectivity() error {
	LogInfo("Checking internet connectivity")
	// Try to connect to a reliable host
	conn, err := net.Dial("tcp", "8.8.8.8:53")
	if err != nil {
		LogError("No internet connectivity detected", err)
		return fmt.Errorf("no internet connectivity detected")
	}
	defer conn.Close()
	LogInfo("Internet connectivity confirmed")
	return nil
}

// CheckCollectionsPrerequisites valida o diretório de collections offline
func CheckCollectionsPrerequisites(collectionsPath string) error {
	LogInfo("Checking collections prerequisites", "path", collectionsPath)
	if _, err := os.Stat(collectionsPath); os.IsNotExist(err) {
		LogError("Collections path does not exist", err, "path", collectionsPath)
		return fmt.Errorf("collections path does not exist: %s", collectionsPath)
	}
	info, err := os.Stat(collectionsPath)
	if err != nil {
		LogError("Cannot stat collections path", err, "path", collectionsPath)
		return err
	}
	if !info.IsDir() {
		LogError("Collections path is not a directory", nil, "path", collectionsPath)
		return fmt.Errorf("collections path is not a directory: %s", collectionsPath)
	}
	entries, err := os.ReadDir(collectionsPath)
	if err != nil {
		LogError("Cannot read collections directory", err, "path", collectionsPath)
		return err
	}
	if len(entries) == 0 {
		LogError("No collection files found in directory", nil, "path", collectionsPath)
		return fmt.Errorf("no collection files found in directory: %s", collectionsPath)
	}
	LogInfo("Collections directory check passed", "path", collectionsPath)
	return nil
}

// CheckRequirementsPrerequisites verifies prerequisites for requirements offline installation.
func CheckRequirementsPrerequisites(requirementsPath string) error {
	LogInfo("Checking requirements prerequisites", "path", requirementsPath)

	// Check if directory exists
	if _, err := os.Stat(requirementsPath); os.IsNotExist(err) {
		LogError("Requirements path does not exist", err, "path", requirementsPath)
		return fmt.Errorf("requirements path does not exist: %s", requirementsPath)
	}

	// Check if requirements.txt exists
	requirementsFile := filepath.Join(requirementsPath, "requirements.txt")
	if _, err := os.Stat(requirementsFile); os.IsNotExist(err) {
		LogError("requirements.txt not found", err, "file", requirementsFile)
		return fmt.Errorf("requirements.txt not found: %s", requirementsFile)
	}

	// Check if directory contains Python packages
	entries, err := os.ReadDir(requirementsPath)
	if err != nil {
		LogError("Cannot read requirements directory", err, "path", requirementsPath)
		return fmt.Errorf("cannot read requirements directory: %v", err)
	}

	packageFound := false
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(name, ".whl") || strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".tgz") {
				packageFound = true
				break
			}
		}
	}

	if !packageFound {
		LogError("No Python packages found in requirements directory", nil, "path", requirementsPath)
		return fmt.Errorf("no Python packages found in requirements directory: %s", requirementsPath)
	}

	LogInfo("Requirements prerequisites check passed", "path", requirementsPath, "entries", len(entries))
	return nil
}

// ValidatePath validates if a path exists and is accessible.
func ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("invalid tarball path: %s", path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	return nil
}

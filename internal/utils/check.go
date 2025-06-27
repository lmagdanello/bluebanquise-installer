package utils

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SystemCheck verifies if the system has the necessary prerequisites
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

// CheckOfflinePrerequisites verifies prerequisites for offline installation
func CheckOfflinePrerequisites(collectionsPath string) error {
	LogInfo("Checking offline prerequisites", "path", collectionsPath)

	// Check if directory exists
	if _, err := os.Stat(collectionsPath); os.IsNotExist(err) {
		LogError("Collections path does not exist", err, "path", collectionsPath)
		return fmt.Errorf("collections path does not exist: %s", collectionsPath)
	}

	// Check if it contains collection files
	entries, err := os.ReadDir(collectionsPath)
	if err != nil {
		LogError("Cannot read collections directory", err, "path", collectionsPath)
		return fmt.Errorf("cannot read collections directory: %v", err)
	}

	if len(entries) == 0 {
		LogError("Collections directory is empty", nil, "path", collectionsPath)
		return fmt.Errorf("collections directory is empty: %s", collectionsPath)
	}

	LogInfo("Offline prerequisites check passed", "path", collectionsPath, "entries", len(entries))
	return nil
}

// CheckRequirementsPrerequisites verifies prerequisites for requirements offline installation
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

// CheckTarballPrerequisites verifies prerequisites for tarball installation
func CheckTarballPrerequisites(tarballPath string) error {
	LogInfo("Checking tarball prerequisites", "path", tarballPath)

	// Check if path exists
	if _, err := os.Stat(tarballPath); os.IsNotExist(err) {
		LogError("Tarball path does not exist", err, "path", tarballPath)
		return fmt.Errorf("tarball path does not exist: %s", tarballPath)
	}

	info, err := os.Stat(tarballPath)
	if err != nil {
		LogError("Cannot stat tarball path", err, "path", tarballPath)
		return fmt.Errorf("cannot stat tarball path: %v", err)
	}

	if !info.IsDir() {
		// If it's a file, check if it's a tarball
		name := info.Name()
		if !strings.HasSuffix(name, ".tar.gz") && !strings.HasSuffix(name, ".tgz") {
			LogError("File is not a tarball", nil, "path", tarballPath)
			return fmt.Errorf("file is not a tarball: %s", tarballPath)
		}
		LogInfo("Tarball file check passed", "path", tarballPath)
		return nil
	}

	// If it's a directory, check if it contains tarballs
	entries, err := os.ReadDir(tarballPath)
	if err != nil {
		LogError("Cannot read tarball directory", err, "path", tarballPath)
		return fmt.Errorf("cannot read tarball directory: %v", err)
	}

	tarballFound := false
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".tgz") {
				tarballFound = true
				break
			}
		}
	}

	if !tarballFound {
		LogError("No tarball files found in directory", nil, "path", tarballPath)
		return fmt.Errorf("no tarball files found in directory: %s", tarballPath)
	}

	LogInfo("Tarball directory check passed", "path", tarballPath)
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

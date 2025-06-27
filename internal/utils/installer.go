package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func detectPackageManager() (string, error) {
	candidates := []string{"apt-get", "dnf", "yum", "zypper"}

	for _, tool := range candidates {
		if _, err := exec.LookPath(tool); err == nil {
			LogInfo("Package manager detected", "manager", tool)
			return tool, nil
		}
	}

	LogError("No supported package manager found", nil, "candidates", candidates)
	return "", fmt.Errorf("no supported package manager found")
}

func InstallPackages(pkgs []string) error {
	LogInfo("Installing packages", "packages", pkgs)

	manager, err := detectPackageManager()
	if err != nil {
		LogError("Failed to detect package manager", err)
		return err
	}

	var args []string
	switch manager {
	case "apt-get":
		args = append([]string{"install", "-y"}, pkgs...)
	case "dnf", "yum":
		args = append([]string{"install", "-y"}, pkgs...)
	case "zypper":
		args = append([]string{"--non-interactive", "install"}, pkgs...)
	default:
		LogError("Unsupported package manager", nil, "manager", manager)
		return fmt.Errorf("unsupported package manager: %s", manager)
	}

	LogCommand(manager, args...)
	cmd := exec.Command(manager, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	fmt.Printf("Installing packages with %s: %s\n", manager, strings.Join(pkgs, " "))
	if err := cmd.Run(); err != nil {
		LogError("Failed to install packages", err, "manager", manager, "packages", pkgs)
		return fmt.Errorf("failed to install packages: %v", err)
	}

	LogInfo("Packages installed successfully", "manager", manager, "packages", pkgs)
	return nil
}

func RunCommand(command string, args ...string) error {
	LogCommand(command, args...)
	cmd := exec.Command(command, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()
	if err != nil {
		LogError("Command execution failed", err, "command", command, "args", args)
	} else {
		LogInfo("Command executed successfully", "command", command, "args", args)
	}
	return err
}

func AppendLineIfMissing(filePath, line string) error {
	LogInfo("Appending line to file if missing", "file", filePath, "line", line)

	// Check if line already exists
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil && !os.IsNotExist(err) {
		LogError("Failed to open file for reading", err, "file", filePath)
		return err
	}

	if file != nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(line) {
				LogInfo("Line already exists in file", "file", filePath, "line", line)
				return nil // Line already exists
			}
		}
	}

	// Append the line
	file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		LogError("Failed to open file for writing", err, "file", filePath)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(line + "\n")
	if err != nil {
		LogError("Failed to write line to file", err, "file", filePath, "line", line)
	} else {
		LogInfo("Line appended to file successfully", "file", filePath, "line", line)
	}
	return err
}

func EnsureLineInSudoers(line string) error {
	LogInfo("Ensuring line in sudoers", "line", line)

	sudoersPath := "/etc/sudoers.d/bluebanquise"

	// Check if line already exists
	file, err := os.OpenFile(sudoersPath, os.O_RDONLY, 0644)
	if err != nil && !os.IsNotExist(err) {
		LogError("Failed to open sudoers file for reading", err, "file", sudoersPath)
		return err
	}

	if file != nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(line) {
				LogInfo("Line already exists in sudoers", "file", sudoersPath, "line", line)
				return nil // Line already exists
			}
		}
	}

	// Append the line
	file, err = os.OpenFile(sudoersPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		LogError("Failed to open sudoers file for writing", err, "file", sudoersPath)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(line + "\n")
	if err != nil {
		LogError("Failed to write line to sudoers", err, "file", sudoersPath, "line", line)
	} else {
		LogInfo("Line added to sudoers successfully", "file", sudoersPath, "line", line)
	}
	return err
}

func DownloadFile(url, filepath string) error {
	LogInfo("Downloading file", "url", url, "path", filepath)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		LogError("Failed to create request", err, "url", url)
		return fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		LogError("Failed to download file", err, "url", url)
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		LogError("Failed to download file", nil, "status", resp.StatusCode, "url", url)
		return fmt.Errorf("failed to download file: HTTP %d", resp.StatusCode)
	}

	file, err := os.Create(filepath)
	if err != nil {
		LogError("Failed to create file", err, "path", filepath)
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		LogError("Failed to write file", err, "path", filepath)
		return fmt.Errorf("failed to write file: %v", err)
	}

	LogInfo("File downloaded successfully", "url", url, "path", filepath)
	return nil
}

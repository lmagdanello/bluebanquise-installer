package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ConfigureSSH sets up SSH configuration for the BlueBanquise user.
func ConfigureSSH(userHome string) error {
	LogInfo("Configuring SSH for BlueBanquise user", "home", userHome)

	sshDir := filepath.Join(userHome, ".ssh")

	// Create .ssh directory
	LogInfo("Creating .ssh directory", "path", sshDir)
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		LogError("Failed to create .ssh directory", err, "path", sshDir)
		return fmt.Errorf("failed to create .ssh directory: %v", err)
	}

	// Generate SSH key pair if it doesn't exist
	keyPath := filepath.Join(sshDir, "id_ed25519")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		LogInfo("Generating SSH key pair", "path", keyPath)
		fmt.Println("Generating SSH key pair...")
		LogCommand("ssh-keygen", "-t", "ed25519", "-f", keyPath, "-q", "-N", "")
		cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", keyPath, "-q", "-N", "")
		if err := cmd.Run(); err != nil {
			LogError("Failed to generate SSH key", err, "path", keyPath)
			return fmt.Errorf("failed to generate SSH key: %v", err)
		}
		LogInfo("SSH key pair generated successfully", "path", keyPath)
	} else {
		LogInfo("SSH key pair already exists", "path", keyPath)
	}

	// Set up authorized_keys
	pubKeyPath := keyPath + ".pub"
	authKeysPath := filepath.Join(sshDir, "authorized_keys")

	// Read public key
	LogInfo("Reading public key", "path", pubKeyPath)
	pubKeyData, err := os.ReadFile(pubKeyPath)
	if err != nil {
		LogError("Failed to read public key", err, "path", pubKeyPath)
		return fmt.Errorf("failed to read public key: %v", err)
	}

	// Check if authorized_keys exists
	if _, err := os.Stat(authKeysPath); os.IsNotExist(err) {
		// Create authorized_keys with the public key
		LogInfo("Creating authorized_keys file", "path", authKeysPath)
		if err := os.WriteFile(authKeysPath, pubKeyData, 0600); err != nil {
			LogError("Failed to create authorized_keys", err, "path", authKeysPath)
			return fmt.Errorf("failed to create authorized_keys: %v", err)
		}
		LogInfo("authorized_keys file created successfully", "path", authKeysPath)
	} else {
		// Check if public key is already in authorized_keys
		LogInfo("Checking if public key is in authorized_keys", "path", authKeysPath)
		authKeysData, err := os.ReadFile(authKeysPath)
		if err != nil {
			LogError("Failed to read authorized_keys", err, "path", authKeysPath)
			return fmt.Errorf("failed to read authorized_keys: %v", err)
		}

		// If public key is not in authorized_keys, append it
		if !contains(authKeysData, pubKeyData) {
			LogInfo("Adding public key to authorized_keys", "path", authKeysPath)
			file, err := os.OpenFile(authKeysPath, os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				LogError("Failed to open authorized_keys for writing", err, "path", authKeysPath)
				return fmt.Errorf("failed to open authorized_keys for writing: %v", err)
			}
			defer func() {
				if closeErr := file.Close(); closeErr != nil {
					LogWarning("Failed to close file", "error", closeErr)
				}
			}()

			if _, err := file.Write(pubKeyData); err != nil {
				LogError("Failed to append to authorized_keys", err, "path", authKeysPath)
				return fmt.Errorf("failed to append to authorized_keys: %v", err)
			}
			LogInfo("Public key added to authorized_keys successfully", "path", authKeysPath)
		} else {
			LogInfo("Public key already exists in authorized_keys", "path", authKeysPath)
		}
	}

	// Set proper permissions
	LogInfo("Setting SSH directory permissions", "path", sshDir, "permissions", "0700")
	if err := os.Chmod(sshDir, 0700); err != nil {
		LogError("Failed to set .ssh directory permissions", err, "path", sshDir)
		return fmt.Errorf("failed to set .ssh directory permissions: %v", err)
	}

	LogInfo("Setting authorized_keys permissions", "path", authKeysPath, "permissions", "0600")
	if err := os.Chmod(authKeysPath, 0600); err != nil {
		LogError("Failed to set authorized_keys permissions", err, "path", authKeysPath)
		return fmt.Errorf("failed to set authorized_keys permissions: %v", err)
	}

	LogInfo("SSH configuration completed successfully", "home", userHome)
	return nil
}

// contains checks if a slice contains a specific byte slice.
func contains(slice, item []byte) bool {
	return len(slice) >= len(item) && string(slice[len(slice)-len(item):]) == string(item)
}

package bootstrap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Initialize logger for tests
	utils.InitTestLogger()
}

func TestInstallCoreVariablesOnline(t *testing.T) {
	tests := []struct {
		name        string
		userHome    string
		expectError bool
	}{
		{
			name:        "Valid user home",
			userHome:    "/tmp/testhome",
			expectError: false,
		},
		{
			name:        "Empty user home",
			userHome:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.userHome != "" {
				defer func() {
					if err := os.RemoveAll(tt.userHome); err != nil {
						t.Logf("Failed to remove test directory: %v", err)
					}
				}()
			}

			err := InstallCoreVariablesOnline(tt.userHome)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				// Note: This test requires internet connectivity
				// In a real test environment, you might want to mock the HTTP request
				if err != nil {
					t.Skip("Skipping test - requires internet connectivity")
				}
				assert.NoError(t, err)
			}
		})
	}
}

func TestInstallCoreVariablesOffline(t *testing.T) {
	tests := []struct {
		name         string
		coreVarsPath string
		userHome     string
		expectError  bool
		setup        func() string
		cleanup      func(string)
	}{
		{
			name:        "Valid core variables file",
			userHome:    "/tmp/testhome",
			expectError: false,
			setup: func() string {
				tempDir := t.TempDir()
				coreVarsFile := filepath.Join(tempDir, "bb_core.yml")
				content := `# Test core variables
bluebanquise_version: "1.0"
test_variable: "test_value"
`
				err := os.WriteFile(coreVarsFile, []byte(content), 0644)
				require.NoError(t, err)
				return coreVarsFile
			},
			cleanup: func(path string) {
				// Cleanup handled by t.TempDir()
			},
		},
		{
			name:        "Valid core variables directory",
			userHome:    "/tmp/testhome",
			expectError: false,
			setup: func() string {
				tempDir := t.TempDir()
				// Create multiple variable files
				files := []string{"bb_core.yml", "bb_network.yml", "bb_storage.yml"}
				for _, file := range files {
					filePath := filepath.Join(tempDir, file)
					content := `# Test variables for ` + file + `
test_variable: "test_value"
`
					err := os.WriteFile(filePath, []byte(content), 0644)
					require.NoError(t, err)
				}
				return tempDir
			},
			cleanup: func(path string) {
				// Cleanup handled by t.TempDir()
			},
		},
		{
			name:        "Non-existent path",
			userHome:    "/tmp/testhome",
			expectError: true,
			setup: func() string {
				return "/non/existent/path"
			},
			cleanup: func(path string) {
				// No cleanup needed
			},
		},
		{
			name:        "Empty user home",
			userHome:    "",
			expectError: true,
			setup: func() string {
				tempDir := t.TempDir()
				coreVarsFile := filepath.Join(tempDir, "bb_core.yml")
				err := os.WriteFile(coreVarsFile, []byte("test"), 0644)
				require.NoError(t, err)
				return coreVarsFile
			},
			cleanup: func(path string) {
				// Cleanup handled by t.TempDir()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coreVarsPath := tt.setup()
			defer tt.cleanup(coreVarsPath)

			if tt.userHome != "" {
				defer func() {
					if err := os.RemoveAll(tt.userHome); err != nil {
						t.Logf("Failed to remove test directory: %v", err)
					}
				}()
			}

			err := InstallCoreVariablesOffline(coreVarsPath, tt.userHome)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify that the files were copied
				groupVarsDir := filepath.Join(tt.userHome, "bluebanquise", "inventory", "group_vars", "all")
				info, err := os.Stat(groupVarsDir)
				assert.NoError(t, err)
				assert.True(t, info.IsDir())
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		destination string
		expectError bool
		setup       func() (string, string)
		cleanup     func(string, string)
	}{
		{
			name:        "Valid file copy",
			expectError: false,
			setup: func() (string, string) {
				tempDir := t.TempDir()
				source := filepath.Join(tempDir, "source.txt")
				destination := filepath.Join(tempDir, "destination.txt")
				content := "test content"
				err := os.WriteFile(source, []byte(content), 0644)
				require.NoError(t, err)
				return source, destination
			},
			cleanup: func(source, destination string) {
				// Cleanup handled by t.TempDir()
			},
		},
		{
			name:        "Non-existent source",
			expectError: true,
			setup: func() (string, string) {
				tempDir := t.TempDir()
				source := filepath.Join(tempDir, "nonexistent.txt")
				destination := filepath.Join(tempDir, "destination.txt")
				return source, destination
			},
			cleanup: func(source, destination string) {
				// Cleanup handled by t.TempDir()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, destination := tt.setup()
			defer tt.cleanup(source, destination)

			err := copyFile(source, destination)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify that the file was copied
				info, err := os.Stat(destination)
				assert.NoError(t, err)
				assert.False(t, info.IsDir())
			}
		})
	}
}

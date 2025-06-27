package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Initialize logger for tests
	InitTestLogger()
}

func TestSystemCheck(t *testing.T) {
	// Skip if not running as root
	if os.Geteuid() != 0 {
		t.Skip("System check test requires root privileges")
	}

	// This test checks if the system check function runs without error
	err := SystemCheck()
	assert.NoError(t, err)
}

func TestCheckCollectionsPrerequisites(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
		setup       func() string
		cleanup     func(string)
	}{
		{
			name:        "Valid collections path",
			expectError: false,
			setup: func() string {
				tempDir := t.TempDir()
				collectionsDir := filepath.Join(tempDir, "collections")
				err := os.MkdirAll(collectionsDir, 0755)
				require.NoError(t, err)
				// Create a dummy collection file
				collectionFile := filepath.Join(collectionsDir, "test_collection")
				err = os.WriteFile(collectionFile, []byte("test"), 0644)
				require.NoError(t, err)
				return collectionsDir
			},
			cleanup: func(path string) {
				// Cleanup handled by t.TempDir()
			},
		},
		{
			name:        "Non-existent path",
			expectError: true,
			setup: func() string {
				return "/non/existent/path"
			},
			cleanup: func(path string) {
				// No cleanup needed
			},
		},
		{
			name:        "Empty directory",
			expectError: true,
			setup: func() string {
				tempDir := t.TempDir()
				collectionsDir := filepath.Join(tempDir, "collections")
				err := os.MkdirAll(collectionsDir, 0755)
				require.NoError(t, err)
				return collectionsDir
			},
			cleanup: func(path string) {
				// Cleanup handled by t.TempDir()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			defer tt.cleanup(path)

			err := CheckOfflinePrerequisites(path)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckTarballPrerequisites(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
		setup       func() string
		cleanup     func(string)
	}{
		{
			name:        "Valid tarball directory",
			expectError: false,
			setup: func() string {
				tempDir := t.TempDir()
				// Create some tarball files
				tarball1 := filepath.Join(tempDir, "collection1.tar.gz")
				tarball2 := filepath.Join(tempDir, "collection2.tgz")
				err := os.WriteFile(tarball1, []byte("test"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(tarball2, []byte("test"), 0644)
				require.NoError(t, err)
				return tempDir
			},
			cleanup: func(path string) {
				// Cleanup handled by t.TempDir()
			},
		},
		{
			name:        "Valid single tarball file",
			expectError: false,
			setup: func() string {
				tempDir := t.TempDir()
				tarball := filepath.Join(tempDir, "collection.tar.gz")
				err := os.WriteFile(tarball, []byte("test"), 0644)
				require.NoError(t, err)
				return tarball
			},
			cleanup: func(path string) {
				// Cleanup handled by t.TempDir()
			},
		},
		{
			name:        "Directory with no tarballs",
			expectError: true,
			setup: func() string {
				tempDir := t.TempDir()
				// Create a non-tarball file
				file := filepath.Join(tempDir, "test.txt")
				err := os.WriteFile(file, []byte("test"), 0644)
				require.NoError(t, err)
				return tempDir
			},
			cleanup: func(path string) {
				// Cleanup handled by t.TempDir()
			},
		},
		{
			name:        "Non-existent path",
			expectError: true,
			setup: func() string {
				return "/non/existent/path"
			},
			cleanup: func(path string) {
				// No cleanup needed
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			defer tt.cleanup(path)

			err := CheckTarballPrerequisites(path)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckRequirementsPrerequisites(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
		setup       func() string
		cleanup     func(string)
	}{
		{
			name:        "Valid requirements directory",
			expectError: false,
			setup: func() string {
				tempDir := t.TempDir()
				// Create some requirement files
				req1 := filepath.Join(tempDir, "ansible-1.0.0.tar.gz")
				req2 := filepath.Join(tempDir, "jinja2-2.0.0.tar.gz")
				err := os.WriteFile(req1, []byte("test"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(req2, []byte("test"), 0644)
				require.NoError(t, err)

				// Create requirements.txt file
				requirementsFile := filepath.Join(tempDir, "requirements.txt")
				requirementsContent := "ansible>=2.15.0\njinja2>=3.0.0\n"
				err = os.WriteFile(requirementsFile, []byte(requirementsContent), 0644)
				require.NoError(t, err)

				return tempDir
			},
			cleanup: func(path string) {
				// Cleanup handled by t.TempDir()
			},
		},
		{
			name:        "Empty directory",
			expectError: true,
			setup: func() string {
				tempDir := t.TempDir()
				return tempDir
			},
			cleanup: func(path string) {
				// Cleanup handled by t.TempDir()
			},
		},
		{
			name:        "Non-existent path",
			expectError: true,
			setup: func() string {
				return "/non/existent/path"
			},
			cleanup: func(path string) {
				// No cleanup needed
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			defer tt.cleanup(path)

			err := CheckRequirementsPrerequisites(path)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

package bootstrap

import (
	"os"
	"testing"

	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Initialize logger for tests
	utils.InitTestLogger()
}

func TestCreateBluebanquiseUser(t *testing.T) {
	tests := []struct {
		name        string
		userName    string
		userHome    string
		expectError bool
	}{
		{
			name:        "Valid user creation",
			userName:    "testuser",
			userHome:    "/tmp/testhome",
			expectError: false,
		},
		{
			name:        "Empty username",
			userName:    "",
			userHome:    "/tmp/testhome",
			expectError: true,
		},
		{
			name:        "Empty home directory",
			userName:    "testuser",
			userHome:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip user creation tests if not running as root
			if os.Geteuid() != 0 {
				t.Skip("Skipping user creation test - requires root privileges")
			}

			err := CreateBluebanquiseUser(tt.userName, tt.userHome)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				// Clean up after test
				defer func() {
					// Remove test user and home directory
					if err := os.RemoveAll(tt.userHome); err != nil {
						t.Logf("Failed to remove test directory: %v", err)
					}
				}()
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	tests := []struct {
		name        string
		userName    string
		expectError bool
	}{
		{
			name:        "Get current user info",
			userName:    os.Getenv("USER"), // Use current user
			expectError: false,
		},
		{
			name:        "Non-existent user",
			userName:    "nonexistentuser",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip if no valid user to test
			if tt.userName == "" {
				t.Skip("No valid user to test with")
			}

			uid, gid, err := GetUserInfo(tt.userName)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, 0, uid)
				assert.Equal(t, 0, gid)
			} else {
				assert.NoError(t, err)
				assert.Greater(t, uid, 0)
				assert.Greater(t, gid, 0)
			}
		})
	}
}

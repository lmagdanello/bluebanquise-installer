package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
)

func CreateBluebanquiseUser(userName, userHome string) error {
	utils.LogInfo("Creating BlueBanquise user", "user", userName, "home", userHome)
	fmt.Printf("Creating %s user... ", userName)

	// Default UID/GID for bluebanquise user
	uid := "377"
	gid := "377"

	// Check if group exists
	if err := exec.Command("getent", "group", userName).Run(); err != nil {
		utils.LogInfo("Creating group", "group", userName, "gid", gid)
		cmd := exec.Command("groupadd", "--gid", gid, userName)
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Run(); err != nil {
			utils.LogError("Failed to create group", err, "group", userName, "gid", gid)
			return fmt.Errorf("failed to create group: %v", err)
		}
	} else {
		utils.LogInfo("Group already exists", "group", userName)
	}

	// Check if user exists
	if err := exec.Command("getent", "passwd", userName).Run(); err != nil {
		utils.LogInfo("Creating user", "user", userName, "uid", uid, "gid", gid, "home", userHome)
		cmd := exec.Command("useradd",
			"--gid", gid,
			"--uid", uid,
			"--create-home",
			"--home-dir", userHome,
			"--shell", "/bin/bash",
			"--system", userName)
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Run(); err != nil {
			utils.LogError("Failed to create user", err, "user", userName, "uid", uid, "gid", gid)
			return fmt.Errorf("failed to create user: %v", err)
		}
	} else {
		utils.LogInfo("User already exists", "user", userName)
	}

	// Create sudoers entry
	sudoers := fmt.Sprintf("%s ALL=(ALL:ALL) NOPASSWD:ALL\n", userName)
	sudoersPath := fmt.Sprintf("/etc/sudoers.d/%s", userName)
	utils.LogInfo("Creating sudoers entry", "user", userName, "path", sudoersPath)
	if err := os.WriteFile(sudoersPath, []byte(sudoers), 0644); err != nil {
		utils.LogError("Failed to write sudoers file", err, "path", sudoersPath)
		return fmt.Errorf("failed to write sudoers file: %v", err)
	}

	utils.LogInfo("BlueBanquise user created successfully", "user", userName, "home", userHome)
	fmt.Println("OK")
	return nil
}

// GetUserInfo returns UID and GID for a given user
func GetUserInfo(userName string) (int, int, error) {
	utils.LogInfo("Getting user info", "user", userName)

	// Get UID
	uidCmd := exec.Command("id", "-u", userName)
	uidBytes, err := uidCmd.Output()
	if err != nil {
		utils.LogError("Failed to get UID", err, "user", userName)
		return 0, 0, fmt.Errorf("failed to get UID for user %s: %v", userName, err)
	}
	uid, err := strconv.Atoi(string(uidBytes[:len(uidBytes)-1])) // Remove newline
	if err != nil {
		utils.LogError("Failed to parse UID", err, "user", userName, "uid_bytes", string(uidBytes))
		return 0, 0, fmt.Errorf("failed to parse UID: %v", err)
	}

	// Get GID
	gidCmd := exec.Command("id", "-g", userName)
	gidBytes, err := gidCmd.Output()
	if err != nil {
		utils.LogError("Failed to get GID", err, "user", userName)
		return 0, 0, fmt.Errorf("failed to get GID for user %s: %v", userName, err)
	}
	gid, err := strconv.Atoi(string(gidBytes[:len(gidBytes)-1])) // Remove newline
	if err != nil {
		utils.LogError("Failed to parse GID", err, "user", userName, "gid_bytes", string(gidBytes))
		return 0, 0, fmt.Errorf("failed to parse GID: %v", err)
	}

	utils.LogInfo("User info retrieved", "user", userName, "uid", uid, "gid", gid)
	return uid, gid, nil
}

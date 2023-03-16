// This file contains reusable code relating to backups

package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*
Creates a backup of any file and adds a unix timestamp.
Used to backup postgresql.auto.conf
@srcPath - path to the file being backed up (postgresql.auto.conf)
@backupDir - path to directory to backup file in (/usr/local/postgrescrutiniser/backups)
*/
func BackupFile(srcPath string, backupDir string, appUser *User, logger *Logger) error {
	// 1. Get the filename from the source path
	filename := filepath.Base(srcPath)

	// 2. Get the current Unix timestamp
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// 3. Backup file to backupDir
	destPath := fmt.Sprintf("%s/%s_%s", backupDir, filename, timestamp)
	cmd := exec.Command("sudo", "cp", srcPath, destPath)
	if err := cmd.Run(); err != nil {
		logger.LogError(fmt.Sprintf("error creating backup: %s", strings.Join(cmd.Args, " ")))
		return err
	}

	// 4. Change owner to postgrescrutiniser
	cmd = exec.Command("sudo", "chown", appUser.Username+".", destPath) // sudo chown user. /path
	if err := cmd.Run(); err != nil {
		logger.LogError(fmt.Sprintf("error creating backup: %s", strings.Join(cmd.Args, " ")))
		return err
	}

	return nil
}

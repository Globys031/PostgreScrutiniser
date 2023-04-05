// This file contains reusable code relating to backups

package utils

import (
	"database/sql"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
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
		logger.LogError(fmt.Errorf("error creating backup: %s", strings.Join(cmd.Args, " ")))
		return err
	}

	// 4. Change owner to postgrescrutiniser
	cmd = exec.Command("sudo", "chown", appUser.Username+".", destPath) // sudo chown user. /path
	if err := cmd.Run(); err != nil {
		logger.LogError(fmt.Errorf("error creating backup: %s", strings.Join(cmd.Args, " ")))
		return err
	}

	return nil
}

// Function for getting when a backup was created
func GetDateTime(path string, logger *Logger) (*time.Time, error) {
	// 1. Get the unix timestamp from the file path
	regex, _ := regexp.Compile(`(\d{10})$`)
	timestamp := regex.FindStringSubmatch(path)

	// Convert it to a date time format that javascript will be able to parse
	unixStamp, err := strconv.ParseInt(timestamp[0], 10, 64)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to convert unix timestamp to datetime: %v", err))
		return nil, err
	}
	datetime := time.Unix(unixStamp, 0)

	return &datetime, nil
}

// Reloads postgresql.conf and postgresql.auto.conf
func ReloadConfiguration(db *sql.DB, logger *Logger) error {
	// Not using `pg_reload_conf()` as some settings only take effect after restarting postgresql itself
	cmd := exec.Command("sudo", "systemctl", "restart", "postgresql*")
	if err := cmd.Run(); err != nil {
		logger.LogError(fmt.Errorf("failed to restart PostgreSQL service: %v", err))
		return err
	}

	return nil
}

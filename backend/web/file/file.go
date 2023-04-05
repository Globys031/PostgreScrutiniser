package file

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
)

// Gets a list of currently available backups
// @path - path to where backups are located
// @currentFile - current postgresql.auto.conf file location
func ListBackups(path, currentFile string, logger *utils.Logger) (*[]BackupFile, error) {
	// 1. Read list of backups
	files, err := ioutil.ReadDir(path)
	if err != nil {
		err = fmt.Errorf("failed to list backups: %v", err)
		logger.LogError(err)
		return nil, err
	}

	// 2. Return only backup files. Takes into account there being other files as well
	var backups []BackupFile
	for _, file := range files {
		filename := file.Name()

		match, _ := regexp.MatchString("postgresql.auto.conf", filename)
		if match {
			// 3. Get file diff and return final response object
			fullBackupPath := path + "/" + file.Name()
			diff, err := CompareBackup(fullBackupPath, currentFile, logger);
			if err != nil {
				return nil, err
			}

			datetime, _ := utils.GetDateTime(path+filename, logger)
			backupFile := BackupFile{
				Name: filename,
				Time: *datetime,
				Diff: diff,
			}
			backups = append(backups, backupFile)
		}
	}
	return &backups, err
}

/*
compares backup to current postgresql.auto.conf file
@backupFile - full path to backup postgresql.auto.conf file
@currentFile - full path to currently used postgresql.auto.conf file
*/
func CompareBackup(backupFile, currentFile string, logger *utils.Logger) ([]FileDiffLine, error) {
	// 1. Get content of currently used conf in bytes. Need elevated privileges to do this
	cmd := exec.Command("sudo", "cat", currentFile)
	stdout, _ := cmd.StdoutPipe()

	cmd.Start()
	currentBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		logger.LogError(fmt.Errorf("failed reading current postgresql.auto.conf file: %v", err))
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		logger.LogError(fmt.Errorf("failed reading current postgresql.auto.conf file: %v", err))
		panic(err)
	}

	// 2. Get content of the backup file
	backupBytes, err := ioutil.ReadFile(backupFile)
	if err != nil {
		logger.LogError(fmt.Errorf("failed reading backup postgresql.auto.conf file: %v", err))
		return nil, err
	}

	// 3. Get diff between the two files
	dmp := diffmatchpatch.New()
	linesCurrent, linesBackup, lineArray := dmp.DiffLinesToChars(string(currentBytes), string(backupBytes))
	diffByChar := dmp.DiffMain(linesCurrent, linesBackup, false)
	diffsByLine := dmp.DiffCharsToLines(diffByChar, lineArray)

	// 4. prepare result for JSON
	var lineDiffs []FileDiffLine
	for _, diff := range diffsByLine {
		diffType := FileDiffLineType(diff.Type.String())
		diffText := diff.Text // Result will get messed up without this intermediate variable
		lineDiffs = append(lineDiffs, FileDiffLine{Line: diffText, Type: diffType})
	}

	return lineDiffs, nil
}

/*
Replaces current postgresql.auto.conf file with backup file and reloads configuration
@backupFile - full path to backup postgresql.auto.conf file
@currentFile - full path to currently used postgresql.auto.conf file
*/
func RestoreBackup(postgresUsername string, backupFile, currentFile string, db *sql.DB, logger *utils.Logger) error {
	cmd := exec.Command("sudo", "mv", backupFile, currentFile)
	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("error replacing current postgresql.auto.conf with backup: %v", err)
		logger.LogError(err)
		return err
	}
	cmd = exec.Command("sudo", "chown", postgresUsername+".", currentFile)
	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("error replacing current postgresql.auto.conf with backup: %v", err)
		logger.LogError(err)
		return err
	}

	err := utils.ReloadConfiguration(db, logger)
	return err
}

// Removes all backups of postgresql.auto.conf
func RemoveBackups(backupDir string, logger *utils.Logger) error {
	files, err := filepath.Glob(backupDir + "/postgresql.auto.conf*")
	if err != nil {
		logger.LogError(fmt.Errorf("error finding postgresql.auto.conf backups: %v", err))
		return err
	}
	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			logger.LogError(fmt.Errorf("error removing file %s: %v", file, err))
			return err
		}
	}
	return nil
}

// Removes postgresql.auto.conf backup
func RemoveBackup(backupFile string, logger *utils.Logger) error {
	err := os.Remove(backupFile)
	if err != nil {
		logger.LogError(fmt.Errorf("error removing %s: %v", backupFile, err))
	}
	return err
}

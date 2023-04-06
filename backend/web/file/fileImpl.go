package file

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FileImpl struct {
	BackupDir        string // directory in which backups are saved
	CurrentFile      string // full path to currently used postgresql.auto.conf
	PostgresUsername string
	AppUser           *utils.User
	Logger           *utils.Logger
	DbHandler        *sql.DB
	Validate         *validator.Validate
}

type AutoConfBackup struct {
	// name of the backup file
	Name string `json:"name" validate:"required,backup"`
}

// Validate request parameter fits postgresql.auto.conf backup regex
// This extra guard is needed for cases where there might be more files in backups dir
// that we don't want to accidentally expose
func (impl *FileImpl) validateBackup(c *gin.Context, backupName string) error {
	request := AutoConfBackup{Name: backupName}
	if err := impl.Validate.Struct(request); err != nil {
		err := fmt.Errorf(`Parameter must match regex: postgresql.auto.conf_(\d{10})$`)
		errorMsg := &ErrorMessage{
			ErrorMessage: err.Error(),
		}
		c.JSON(http.StatusBadRequest, &errorMsg)
		return err
	}
	return nil
}

// Lists all backups
func (impl *FileImpl) GetBackups(c *gin.Context) {
	// Due to how `oapi-codegen` generates code, we have to manually
	// exit at each endpoint if middleware aborted
	if len(c.Errors) > 0 || c.Writer.Status() >= 400 {
		return
	}

	data, err := ListBackups(impl.BackupDir, impl.CurrentFile, impl.Logger)
	if err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, &errorMsg)
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, &errorMsg)
		return
	}

	c.Data(http.StatusAccepted, "application/json", jsonData)
}

// Replaces current postgresql.auto.conf file with backup file and reloads configuration
func (impl *FileImpl) PutBackup(c *gin.Context, backupName string) {
	if len(c.Errors) > 0 || c.Writer.Status() >= 400 {
		return
	}

	// 1. Validate parameter fits regex
	err := impl.validateBackup(c, backupName)
	if err != nil {
		return
	}

	// 2. Replace backup
	fullPath := impl.BackupDir + "/" + backupName
	if err := RestoreBackup(impl.PostgresUsername, fullPath, impl.CurrentFile, impl.AppUser, impl.DbHandler, impl.Logger); err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, &errorMsg)
		return
	}
}

func (impl *FileImpl) DeleteBackups(c *gin.Context) {
	if len(c.Errors) > 0 || c.Writer.Status() >= 400 {
		return
	}

	if err := RemoveBackups(impl.BackupDir, impl.Logger); err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, &errorMsg)
		return
	}
}

func (impl *FileImpl) DeleteBackup(c *gin.Context, backupName string) {
	if len(c.Errors) > 0 || c.Writer.Status() >= 400 {
		return
	}

	// 1. Validate parameter fits regex
	err := impl.validateBackup(c, backupName)
	if err != nil {
		return
	}

	// 2. Delete backup
	fullPath := impl.BackupDir + "/" + backupName
	if err := RemoveBackup(fullPath, impl.Logger); err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, &errorMsg)
		return
	}
}
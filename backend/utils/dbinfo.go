package utils

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"

	_ "github.com/lib/pq"
)

/*
Function for getting database info from ~/.pgpass
@appUser - appUser = "postgrescrutiniser"
*/
func ParsePgpassFile(appUser *User, logger *Logger) (string, string, string, string, string, error) {
	// Get the home directory of our application user
	usr, err := user.Lookup(appUser.Username)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed to find %s home directory: %v", appUser.Username, err.Error()))
		return "", "", "", "", "", err
	}
	homeDir := usr.HomeDir

	// Open the .pgpass file for reading
	pgpassFile, err := os.Open(homeDir + "/.pgpass")
	if err != nil {
		logger.LogError(fmt.Errorf("Failed to find %s home directory: %v", appUser.Username, err.Error()))
		return "", "", "", "", "", err
	}
	defer pgpassFile.Close()

	// Assume that first line contains our data
	scanner := bufio.NewScanner(pgpassFile)
	scanner.Scan()
	line := scanner.Text()
	fields := strings.Split(line, ":")

	return fields[0], fields[1], fields[2], fields[3], fields[4], nil
}

// Intitiate new database connection and return handler for it.
func InitDbConnection(hostname string, user string, passwd string, port string, logger *Logger) (*sql.DB, error) {
	if hostname == "" || user == "" || passwd == "" || port == "" {
		err := fmt.Errorf("Could not initiate database connection bcause one of the fields was empty")
		logger.LogError(err)
		return nil, err
	}

	connString := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable", hostname, user, passwd, port)

	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		logger.LogError(fmt.Errorf("Could not initiate database connection: %v", err))
		return nil, err
	}

	return dbConn, nil
}

func CloseDbConnection(dbHandler *sql.DB, logger *Logger) error {
	if dbHandler == nil {
		return nil
	}

	if err := dbHandler.Close(); err != nil {
		logger.LogError(fmt.Errorf("Could not initiate database connection: %v", err))
		return err
	}
	return nil
}

// Returns path to postgresql.conf if it exists.
func FindConfigFile(dbHandler *sql.DB, logger *Logger) (string, error) {
	row := dbHandler.QueryRow("SHOW config_file")

	var filePath string
	err := row.Scan(&filePath)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed finding postgresql.conf: %v", err))
		return "", err
	}

	// check to confirm postgresql.conf file exists
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		err = fmt.Errorf("postgresql.conf was found with `SHOW config_file` but could not be opened: %v", err)
		logger.LogError(fmt.Errorf("Failed finding postgresql.conf: %v", err))
		return "", err
	}

	return filePath, nil
}

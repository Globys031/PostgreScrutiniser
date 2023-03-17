// File meant for handling error/warning related logic (might add info logging in the future)
// TO DO: add log rotation

package utils

import (
	"log"
	"os"
)

type Logger struct {
	logDir  string // Path to directory where error.log resides in
	console *log.Logger
	file    *log.Logger
}

// Creates `/var/log/postgrescrutiniser/error.log` if it doesn't exist
// and returns a `Logger` object that will be used throughout the application
func InitLogging() *Logger {
	// 1. Check if logDir exists in filesystem
	logDir := "/var/log/postgrescrutiniser"
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Create error.log for appending
	console := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logFile, err := os.OpenFile(logDir+"/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	file := log.New(logFile, "", log.Ldate|log.Ltime)

	// 3. Return logger that will be used through the entire application
	return &Logger{
		logDir:  logDir,
		console: console,
		file:    file,
	}
}

// Log warnings to both console and error.log
func (logger *Logger) LogWarning(message error) {
	logger.console.Println("WARNING:", message.Error())
	logger.file.Println("WARNING:", message.Error())
}

// Log errors to both console and error.log
func (logger *Logger) LogError(message error) {
	logger.console.Println("ERROR:", message.Error())
	logger.file.Println("ERROR:", message.Error())
}

// Log fatal error to both console and error.log and exit
func (logger *Logger) LogFatal(message error) {
	logger.console.Println("FATAL:", message.Error())
	logger.file.Println("FATAL:", message.Error())
	os.Exit(1)
}

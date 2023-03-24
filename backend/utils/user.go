// This file contains reusable code relating to managing users and their permissions

package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

type User struct {
	Username string
	Password string
	Uid      int
	Gid      int
}

/*
Gets uid and gid of user
@username - user whose uid & gid we're trying to fetch
*/
func GetUserIds(username string, logger *Logger) (int, int, error) {
	u, err := user.Lookup(username)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to get uid and gid for user %s: %s", username, err))
		return 0, 0, err
	}
	uid, err := strconv.Atoi(u.Uid)
	gid, err := strconv.Atoi(u.Gid)
	return uid, gid, nil
}

// Function for confirming whether password passed as a parameter matches the one in /etc/shadow
func PasswordMatches(username string, password string, logger *Logger) bool {
	// 1. Get password from /etc/shadow
	// getent shadow postgrescrutiniser | awk -F ':' '{print $2}'
	shadowHash, err := getPasswordHash(username, logger)
	if err != nil {
		return false
	}

	// 2. hash @password using the salt from /etc/shadow
	salt, err := getHashSalt(shadowHash, logger)
	if err != nil {
		return false
	}
	generatedHash, err := generateSHA512Password(password, salt, logger)
	if err != nil {
		return false
	}

	// 3. If the two are the same, @password is correct
	if shadowHash == generatedHash {
		return true
	}
	return false
}

// Function for getting salt from @hash
func getHashSalt(hash string, logger *Logger) (string, error) {
	parts := strings.Split(hash, "$")
	if len(parts) < 3 {
		err := fmt.Errorf("Invalid hash format")
		logger.LogError(err)
		return "", err
	}
	return parts[2], nil
}

// Function for reading the hash from /etc/shadow of @username
func getPasswordHash(username string, logger *Logger) (string, error) {
	// 1. Run `sudo getent shadow postgrescrutiniser | awk -F ':' '{print $2}'`
	var out bytes.Buffer

	cmdGetent := exec.Command("sudo", "getent", "shadow", username)
	cmdAwk := exec.Command("awk", "-F", ":", "{print $2}")

	cmdAwk.Stdin, _ = cmdGetent.StdoutPipe()
	cmdAwk.Stdout = &out

	if err := cmdAwk.Start(); err != nil {
		logger.LogError(fmt.Errorf("Could not get shadow hash: %v", err))
		return "", err
	}
	if err := cmdGetent.Start(); err != nil {
		logger.LogError(fmt.Errorf("Could not get shadow hash: %v", err))
		return "", err
	}
	if err := cmdAwk.Wait(); err != nil {
		logger.LogError(fmt.Errorf("Could not get shadow hash: %v", err))
		return "", err
	}

	// 2. Return password hash from output
	passwordHash := strings.TrimSpace(out.String())
	return passwordHash, nil
}

// Function used to compare hash present in /etc/shadow with hash from supplied password
func generateSHA512Password(password string, salt string, logger *Logger) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("openssl", "passwd", "-6", "-salt", salt, password)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		logger.LogError(fmt.Errorf("Could not generate SHA512 hash: %v", err))
		return "", err
	}

	hashedPassword := strings.ReplaceAll(out.String(), "\n", "")
	return hashedPassword, nil
}

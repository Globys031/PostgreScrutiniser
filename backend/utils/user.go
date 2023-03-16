// This file contains reusable code relating to managing users and their permissions

package utils

import (
	"fmt"
	"os/user"
	"strconv"
)

type User struct {
	Username string
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
		logger.LogError(fmt.Sprintf("failed to get uid and gid for user %s: %s", username, err.Error()))
		return 0, 0, err
	}
	uid, err := strconv.Atoi(u.Uid)
	gid, err := strconv.Atoi(u.Gid)
	return uid, gid, nil
}

// // Function for switching what user the application is running under
// // assumes that user has correct permissions in /etc/sudoers for this
// func switchUsers(user *User, logger *Logger) error {
// 	// switch users
// 	cmd := exec.Command("sudo", "-u", user.Username, "-i")
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	if err := cmd.Run(); err != nil {
// 		return err
// 	}

// 	// Set the UID and GID of the current process to the target user
// 	if err := syscall.Setuid(user.Uid); err != nil {
// 		logger.LogError(fmt.Sprintf("failed to set uid and gid for user %s: %s", user.Username, err.Error()))
// 		return err
// 	}
// 	if err := syscall.Setgid(int(user.Gid)); err != nil {
// 		logger.LogError(fmt.Sprintf("failed to set uid and gid for user %s: %s", user.Username, err.Error()))
// 		return err
// 	}
// 	return nil
// }

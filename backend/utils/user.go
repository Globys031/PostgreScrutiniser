// This file contains reusable code relating to managing users and their permissions

package utils

import (
	"fmt"
	"os/user"
	"strconv"
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

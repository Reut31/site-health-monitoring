package collectors

import (
	"github.com/shirou/gopsutil/v4/host"
)

type LoggedInUser struct {
	Username string
	Terminal string
	Host     string
}

func GetLoggedInUsers() ([]LoggedInUser, error) {
	var users []LoggedInUser
	var newUser LoggedInUser
	result, err := host.Users()
	if err != nil {
		return nil, err
	}
	for _, user := range result {
		newUser.Username = user.User
		newUser.Terminal = user.Terminal
		newUser.Host = user.Host
		users = append(users, newUser)
	}
	return users, nil
}

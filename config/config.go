package config

import (
	"os/user"
	"path/filepath"
)

func DefaultPathForUser(usr *user.User) string {

	home := usr.HomeDir
	config := filepath.Join(home, ".config")
	root := filepath.Join(config, "secretbox")
	path := filepath.Join(root, "salt")

	return path
}

package config

import (
	"fmt"
	"os/user"
)

const (
	DefaultPath = ".x.yml"
)

// config struct
type Config struct {

	OsUser *user.User

	AbsolutePath string
}

func (c *Config) LoadConfig() string {
	return fmt.Sprintf("%s/%s",c.OsUser.HomeDir, DefaultPath)
}

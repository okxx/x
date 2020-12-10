package instance

import (
	"golang.org/x/crypto/ssh"
)

type Instance struct {
	Name           string           `yaml:"name"`
	Alias          string           `yaml:"alias"`
	Host           string           `yaml:"host"`
	User           string           `yaml:"user"`
	Port           string           `yaml:"port"`
	KeyPath        string           `yaml:"keypath"`
	Passphrase     string           `yaml:"passphrase"`
	Password       string           `yaml:"password"`
}

func InitInstance() *Instance {
	instance  := &Instance{}
	return instance
}

// get user
func (i *Instance) GetUser() string {
	if i.User == "" {
		return "root"
	}
	return i.User
}

// get password
func (i *Instance) GetPassword() ssh.AuthMethod {
	if i.Password == "" {
		return nil
	}
	return ssh.Password(i.Password)
}

// get port
func (i *Instance) GetPort() string {
	if i.Port == "" {
		return "22"
	}
	return i.Port
}
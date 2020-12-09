package engine

import (
	"fmt"
	"github.com/laracro/x/internal/config"
	"github.com/laracro/x/internal/instance"
	"os/user"
)

type Engine struct {
	*config.Config

	*instance.Instance
	Instances []instance.Instance
}

// load engine
func LoadEngine() *Engine {

	// init engine
	engine := &Engine{}

	// init config
	engine.Config = &config.Config{}
	engine.Config.OsUser = getOsUser()

	// load Config
	engine.LoadConfig()

	// load instance
	engine.Instance = instance.InitInstance()

	return engine
}

// get os user
func getOsUser () *user.User {
	osUser,err := user.Current()
	if err != nil {
		panic(err)
	}
	return osUser
}

// fetch all instance
func (x *Engine) FetchAll() {
	for _,item := range x.Instances {
		fmt.Printf("｜%s｜%s｜%s\n",item.User,item.Host,item.Port)
	}
}


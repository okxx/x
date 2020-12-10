package engine

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/laracro/x/internal/client"
	"github.com/laracro/x/internal/config"
	"github.com/laracro/x/internal/instance"
	"log"
	"os/user"
	"reflect"
)

type Engine struct {
	*config.Config

	client.Client

	*instance.Instance
	Instances map[string]instance.Instance

	SimpleTable *simpletable.Table
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

	engine.SimpleTable = simpletable.New()

	// load instance
	engine.Instance = instance.InitInstance()

	engine.Instances = map[string]instance.Instance{}

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

func (x *Engine) AddInstance(instance instance.Instance) bool {
	x.Instances[instance.Name] = instance
	return true
}

func (x *Engine) GetInstance(name string) {
	if v,ok := x.Instances[name]; !ok {
		log.Fatal("unknown instance.")
	} else {
		x.Instance = &v
	}
}

// fetch all instance
func (x *Engine) FetchAll() {

	fields := reflect.TypeOf(instance.Instance{})
	for i:=0; i < fields.NumField(); i++ {
		cell := &simpletable.Cell{}
		cell.Text = fields.Field(i).Name
		cell.Align = simpletable.AlignCenter
		x.SimpleTable.Header.Cells = append(x.SimpleTable.Header.Cells,cell)
	}

	if len(x.Instances) > 0 {
		for _, row := range x.Instances {
			value := reflect.ValueOf(row)
			var cells []*simpletable.Cell
			for i:=0; i < value.NumField(); i++ {
				cell := &simpletable.Cell{}
				cell.Align = simpletable.AlignCenter
				cell.Text = value.Field(i).String()
				cells = append(cells,cell)
			}
			x.SimpleTable.Body.Cells = append(x.SimpleTable.Body.Cells, cells)
		}
	}
	fmt.Println(x.SimpleTable.String())
}


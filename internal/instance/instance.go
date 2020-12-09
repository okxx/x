package instance

type Instance struct {
	Host 		string
	Port 		string
	User 		string
	Password 	string
	Alias 		[]string
}

func InitInstance() *Instance {
	instance  := &Instance{}
	return instance
}

func (i *Instance) Open() {

}

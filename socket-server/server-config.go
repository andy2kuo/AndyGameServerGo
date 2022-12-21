package socketserver

type AppSetting struct {
	Server    ServerSetting
	Operation OperationSetting
}

func (AppSetting) Name() string {
	return "AppSetting"
}

type ServerSetting struct {
	Name        string `default:"Socket Server"`
	Environment string `default:"dev"`
	Port        int    `default:"8309"`
	TimeOut     int    `default:"30"`
}

type OperationSetting struct {
	RunMaxTime int `default:"5"`
}

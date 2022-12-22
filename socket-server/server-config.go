package socketserver

type AppSetting struct {
	Server    ServerSetting
	Operation OperationSetting
}

func (AppSetting) Name() string {
	return "AppSetting"
}

type ServerSetting struct {
	Name         string `default:"Socket Server"`
	Environment  string `default:"dev"`
	Port         int    `default:"8309"`
	TimeOut      int    `default:"30"`
	ReadBuffer   int    `default:"1024"`
	ReadTimeOut  int    `default:"5"`
	WriteBuffer  int    `default:"1024"`
	WriteTimeOut int    `default:"5"`
}

type OperationSetting struct {
	RunMaxTime int `default:"5"`
}
